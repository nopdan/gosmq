package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/imetool/gosmq/internal/dict"
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
	goCmd.PersistentFlags().StringArrayVarP(&conf.Dict, "dict", "i", nil, "码表路径，可以为多个")

	goCmd.PersistentFlags().BoolVarP(&conf.Single, "single", "s", false, "启用单字模式")
	goCmd.PersistentFlags().StringVarP(&conf.Algo, "algo", "", "trie", "匹配算法(trie|strie)")
	goCmd.PersistentFlags().BoolVarP(&conf.Stable, "stable", "", false, "按码表顺序(覆盖algo)")
	goCmd.PersistentFlags().StringVarP(&conf.PressSpaceBy, "space", "k", "both", "空格按键方式 left|right|both")
	goCmd.PersistentFlags().BoolVarP(&conf.Split, "split", "", false, "输出分词数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Stat, "stat", "", false, "输出词条数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Json, "json", "", false, "输出json数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出全部数据")
	goCmd.PersistentFlags().BoolVarP(&conf.Hidden, "hidden", "", false, "隐藏 cli 结果展示")
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

	// 判断是否为文件夹
	if len(conf.Text) == 1 {
		path := conf.Text[0]
		fi, err := os.Stat(path)
		if err != nil {
			fmt.Println("找不到文件或文件夹", path)
			panic(err)
		}
		if fi.IsDir() {
			texts := make([]string, 0)
			files, err := os.ReadDir(path)
			if err != nil {
				panic(err)
			}
			conf.isFolder = true
			fmt.Printf("载入 %s 下的文本: \n", path)
			if !strings.HasSuffix(path, "\\") {
				path += "\\"
			}
			for _, file := range files {
				if !file.IsDir() {
					texts = append(texts, path+file.Name())
					fmt.Printf("-> %s\n", file.Name())
				}
			}
			conf.Text = texts
		}
	}

	newDict := func() *dict.Dict {
		return &dict.Dict{
			Single:       conf.Single,
			Algorithm:    conf.Algo,
			PressSpaceBy: conf.PressSpaceBy,
			Split:        conf.Split,
			Stat:         conf.Stat,
			Json:         conf.Json,
		}
	}

	// 对多文本启用多协程
	if len(conf.Text) > 1 {
		start := time.Now()
		dicts := make([]*dict.Dict, len(conf.Dict))
		for i := range dicts {
			dicts[i] = newDict()
			dicts[i].Load(conf.Dict[i])
			fmt.Println("载入码表：", dicts[i].Name)
		}
		if len(conf.Dict) != 1 {
			fmt.Printf("比赛开始，共 %d 个码表\n", len(dicts))
		}
		textTotalLen := int64(0)
		var wg sync.WaitGroup
		ch := make(chan struct{}, 8)
		for _, text := range conf.Text {
			ch <- struct{}{}
			wg.Add(1)
			go func(text string) {
				mid := time.Now()
				// 初始化赛码器
				s := &smq.Smq{}
				err := s.Load(text)
				if err != nil {
					fmt.Println("Error! 读取文件失败：", err)
					<-ch
					wg.Done()
					return
				}
				res := s.EvalDicts(dicts)
				atomic.AddInt64(&textTotalLen, int64(res[0].Basic.TextLen))
				if !conf.isFolder {
					fmt.Println("载入文本：", s.Name)
				}
				if !conf.Hidden {
					fmt.Printf("此文本耗时：%v\n", time.Since(mid))
					printSep()
					Output(res, s.Name)
				}
				<-ch
				wg.Done()
			}(text)
		}
		wg.Wait()
		fmt.Printf("共载入 %d 个文本，总字数 %d，总耗时：%v\n", len(conf.Text), textTotalLen, time.Since(start))
		return
	}

	// 下面针对单文本
	start := time.Now()
	// 初始化赛码器
	s := &smq.Smq{}
	s.Load(conf.Text[0])

	dicts := make([]*dict.Dict, 0)
	// 添加码表
	for _, v := range conf.Dict {
		dict := newDict()
		dict.Load(v)
		fmt.Println("载入码表：", dict.Name)
		dicts = append(dicts, dict)
	}

	// 开始赛码
	if len(conf.Dict) != 1 {
		fmt.Printf("比赛开始，共 %d 个码表\n", len(dicts))
	}
	res := s.EvalDicts(dicts)
	fmt.Printf("总耗时：%v\n", time.Since(start))
	if !conf.Hidden {
		printSep()
		Output(res, s.Name)
	}
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
	s := &smq.Smq{}
	if info.Text == "" {
		fmt.Println("没有输入文本")
		return
	} else {
		err := s.Load(info.Text)
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

	d := &dict.Dict{
		Single:       info.Single,
		Algorithm:    algo,
		PressSpaceBy: info.PressSpaceBy,
	}
	d.Load(info.Dict)
	// 开始赛码
	res := s.Eval(d)
	fmt.Printf("耗时：%v\n", time.Since(start))
	printSep()
	Output([]*smq.Result{res}, s.Name)
}
