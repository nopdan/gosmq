package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/gosmq/pkg/smq"
	"github.com/spf13/cobra"
)

var conf = &struct {
	Text         string   // 文本
	Dict         []string // 码表
	Single       bool     // 单字模式
	Algo         string   // 匹配算法
	PressSpaceBy string   // 空格按键方式 left|right|both
	Verbose      bool     // 输出详细数据
}{}

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "命令行模式赛码",
	Run: func(cmd *cobra.Command, args []string) {
		goCli()
	},
}

func init() {
	goCmd.PersistentFlags().StringVarP(&conf.Text, "text", "t", "", "文本路径")
	goCmd.PersistentFlags().StringArrayVarP(&conf.Dict, "dict", "i", nil, "码表路径")
	goCmd.PersistentFlags().BoolVarP(&conf.Single, "single", "s", false, "启用单字模式")
	goCmd.PersistentFlags().StringVarP(&conf.Algo, "algo", "a", "strie", "匹配算法")
	goCmd.PersistentFlags().StringVarP(&conf.PressSpaceBy, "space", "p", "both", "空格按键方式 left|right|both")
	goCmd.PersistentFlags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出详细数据")
}

func goCli() {
	if len(conf.Dict) == 0 {
		fmt.Println("输入有误")
		return
	}
	// 初始化赛码器
	s := &smq.Smq{}
	start := time.Now()
	if conf.Text == "" {
		fmt.Println("没有输入文本")
		return
	} else {
		s.Load(conf.Text)
	}

	// 添加码表
	for _, v := range conf.Dict {
		dict := &dict.Dict{
			Single:       conf.Single,
			Algorithm:    conf.Algo,
			PressSpaceBy: conf.PressSpaceBy,
			Verbose:      conf.Verbose,
		}
		dict.Load(v)
		s.Add(dict)
	}
	fmt.Printf("构建码表耗时：%v\n", time.Since(start))

	// 开始赛码
	fmt.Printf("比赛开始，一共 %d 个码表\n", len(s.Inputs))
	mid := time.Now()
	res := s.Run()
	fmt.Printf("比赛结束，耗时：%v\n", time.Since(mid))
	fmt.Printf("总耗时：%v\n", time.Since(start))
	if len(res) == 0 {
		return
	}
	fmt.Println("----------------------")
	Output(res, s.Name)
}

func goWithSurvey() {
	handle := func(err error) {
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
	}

	err := survey.AskOne(&survey.Input{
		Message: "文本:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &conf.Text, survey.WithValidator(survey.Required))
	handle(err)

	var tmp string
	err = survey.AskOne(&survey.Input{
		Message: "码表:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &tmp, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Select{
		Message: "空格按键方式:",
		Options: []string{"both", "left", "right"},
	}, &conf.PressSpaceBy)
	handle(err)

	err = survey.AskOne(&survey.Confirm{
		Message: "单字模式:",
		Default: false,
	}, &conf.Single)
	handle(err)

	var greedy bool
	err = survey.AskOne(&survey.Confirm{
		Message: "贪心匹配:",
		Default: false,
	}, &greedy)
	handle(err)

	err = survey.AskOne(&survey.Confirm{
		Message: "输出详细数据:",
		Default: false,
	}, &conf.Verbose)
	handle(err)

	fmt.Println("\n", conf)

	conf.Dict = []string{tmp}

	// 初始化赛码器
	s := &smq.Smq{}
	start := time.Now()
	if conf.Text == "" {
		fmt.Println("没有输入文本")
		return
	} else {
		s.Load(conf.Text)
	}
	if greedy {
		conf.Algo = "trie"
	} else {
		conf.Algo = "strie"
	}
	d := &dict.Dict{
		Single:       conf.Single,
		Algorithm:    conf.Algo,
		PressSpaceBy: conf.PressSpaceBy,
		Verbose:      conf.Verbose,
	}
	d.Load(conf.Dict[0])
	s.Add(d)
	fmt.Printf("构建码表耗时：%v\n", time.Since(start))
	// 开始赛码
	mid := time.Now()
	res := s.Run()
	fmt.Printf("比赛结束，耗时：%v\n", time.Since(mid))
	fmt.Printf("总耗时：%v\n", time.Since(start))
	if len(res) == 0 {
		return
	}
	fmt.Println("----------------------")
	Output(res, s.Name)
}
