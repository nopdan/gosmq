package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/imetool/gosmq/pkg/smq"
	"github.com/spf13/cobra"
)

type Basic struct {
}

var conf = &struct {
	Text []string // 文本
	Dict []string // 码表

	Single       bool   // 单字模式
	Algo         string // 匹配算法
	Stable       bool   // 按码表顺序(覆盖algo)
	PressSpaceBy string // 空格按键方式 left|right|both
	Split        bool   // 输出分词数据
	Stat         bool   // 输出词条数据
	Json         bool   // 输出json数据
	Verbose      bool   // 输出全部数据
	Hidden       bool   // 隐藏 cli 结果展示
	Clean        bool   // 只统计词库中的词条
	Merge        bool   // 合并一码表多文本的结果

	isFolder bool
}{}

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "命令行模式赛码",
	Run: func(cmd *cobra.Command, args []string) {
		goCli()
	},
}

func init() {
	goCmd.PersistentFlags().StringArrayVarP(&conf.Text, "text", "t", nil, "文本路径，可以为多个文件，或一个文件夹")
	goCmd.PersistentFlags().StringArrayVarP(&conf.Dict, "dict", "i", nil, "码表路径，可以为多个文件，或一个文件夹")

	goCmd.PersistentFlags().BoolVarP(&conf.Single, "single", "s", false, "启用单字模式")
	goCmd.PersistentFlags().StringVarP(&conf.Algo, "algo", "", "trie", "匹配算法(trie|strie)")
	goCmd.PersistentFlags().BoolVarP(&conf.Stable, "stable", "", false, "按码表顺序(覆盖algo)")
	goCmd.PersistentFlags().StringVarP(&conf.PressSpaceBy, "space", "k", "both", "空格按键方式 left|right|both")
	goCmd.PersistentFlags().BoolVarP(&conf.Split, "split", "", false, "输出分词数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Stat, "stat", "", false, "输出词条数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Json, "json", "", false, "输出json数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出全部数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Hidden, "hidden", "", false, "隐藏 cli 结果展示")
	goCmd.PersistentFlags().BoolVarP(&conf.Clean, "clean", "c", false, "只统计词库中的词条")
	goCmd.PersistentFlags().BoolVarP(&conf.Merge, "merge", "m", false, "合并一码表多文本的结果")
}

func goCli() {
	if len(conf.Dict) == 0 || len(conf.Text) == 0 {
		fmt.Println("输入有误")
		return
	}
	if conf.Stable {
		conf.Algo = "strie"
	}
	if conf.Verbose {
		conf.Split = true
		conf.Stat = true
		conf.Json = true
	}

	flag := 0
	if conf.Split {
		flag |= smq.S_SPLIT
	}
	if conf.Stat {
		flag |= smq.S_STAT
	}
	if conf.Json {
		flag |= smq.S_JSON
	}

	// 开始计时
	start := time.Now()
	texts := make([]string, 0, len(conf.Text))
	for _, v := range conf.Text {
		texts = append(texts, getFiles(v)...)
	}
	fmt.Println("载入文本：")
	for _, v := range texts {
		if !conf.Hidden {
			fmt.Println("-> ", v)
		}
	}
	fmt.Println()

	dictNames := make([]string, 0, len(conf.Dict))
	for _, v := range conf.Dict {
		dictNames = append(dictNames, getFiles(v)...)
	}
	newDict := func() *smq.Dict {
		return &smq.Dict{
			Single:       conf.Single,
			Algorithm:    conf.Algo,
			PressSpaceBy: conf.PressSpaceBy,
			Verbose:      flag != 0,
			Clean:        conf.Clean,
		}
	}
	dicts := make([]*smq.Dict, 0, len(dictNames))
	fmt.Println("载入码表：")
	dictStartTime := time.Now()
	mid := time.Now()
	for _, v := range dictNames {
		d := newDict()
		d.Load(v)
		dicts = append(dicts, d)
		if !conf.Hidden {
			if len(dictNames) == 1 {
				fmt.Println("=> ", v)
			} else {
				fmt.Println("=> ", v, "\t耗时：", time.Since(mid))
				mid = time.Now()
			}
		}
	}
	fmt.Printf("载入码表耗时：%v\n\n", time.Since(dictStartTime))

	// race
	fmt.Println("比赛开始……")
	textLenTotal := 0
	resArr := smq.Parallel(texts, dicts)

	if conf.Merge {
		resArr2 := transpose(resArr)
		for _, res := range resArr2 {
			res2 := smq.MergeResults(res, conf.Stat)
			printSep()
			Output([]*smq.Result{res2})
			res2.Output(flag)
			textLenTotal = res2.TextLen
		}
		fmt.Printf("共载入 %d 个码表，%d 个文本，总字数 %d，总耗时：%v\n", len(dicts), len(texts), textLenTotal, time.Since(start))
		return
	}

	for _, v := range resArr {
		if len(v) == 0 {
			break
		}
		textLenTotal += v[0].TextLen
		if !conf.Hidden {
			printSep()
			Output(v)
		}
		for _, res := range v {
			res.Output(flag)
		}
	}

	fmt.Printf("共载入 %d 个码表，%d 个文本，总字数 %d，总耗时：%v\n", len(dicts), len(texts), textLenTotal, time.Since(start))
}

func goWithSurvey() {
	handle := func(err error) {
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
	}

	var info = &struct {
		Text         string
		Dict         string
		PressSpaceBy string
		Single       bool
		Stable       bool
	}{}

	err := survey.AskOne(&survey.Input{
		Message: "文本:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &info.Text, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Input{
		Message: "码表:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &info.Dict, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Select{
		Message: "空格按键方式:",
		Options: []string{"both", "left", "right"},
	}, &info.PressSpaceBy)
	handle(err)

	err = survey.AskOne(&survey.Confirm{
		Message: "单字模式:",
		Default: false,
	}, &info.Single)
	handle(err)

	err = survey.AskOne(&survey.Confirm{
		Message: "按码表顺序:",
		Default: false,
	}, &info.Stable)
	handle(err)

	fmt.Printf("\n\n")

	start := time.Now()
	// 初始化赛码器
	t := &smq.Text{}
	if info.Text == "" {
		fmt.Println("没有输入文本")
		return
	} else {
		err := t.Load(info.Text)
		if err != nil {
			log.Panic("Error! 读取文件失败：", err)
		}
	}

	var algo string
	if info.Stable {
		algo = "strie"
	} else {
		algo = "trie"
	}

	d := &smq.Dict{
		Single:       info.Single,
		Algorithm:    algo,
		PressSpaceBy: info.PressSpaceBy,
	}
	d.Load(info.Dict)
	// 开始赛码
	res := t.RaceOne(d)
	fmt.Printf("耗时：%v\n", time.Since(start))
	printSep()
	Output([]*smq.Result{res})
}
