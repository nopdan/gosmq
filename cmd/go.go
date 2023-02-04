package cmd

import (
	"fmt"
	"time"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/gosmq/pkg/smq"
	"github.com/spf13/cobra"
)

var conf = &struct {
	Text         string   // 文本
	Dict         []string // 码表
	Single       bool     // 单字模式
	Greedy       bool     // 贪心匹配
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
	goCmd.PersistentFlags().BoolVarP(&conf.Greedy, "greedy", "g", false, "贪心匹配")
	goCmd.PersistentFlags().StringVarP(&conf.PressSpaceBy, "space", "p", "both", "空格按键方式 left|right|both")
	goCmd.PersistentFlags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出详细数据")

}

func goCli() {
	if len(conf.Dict) == 0 {
		fmt.Println("这点参数让我怎么算")
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
		var algo string
		if conf.Greedy {
			algo = "trie"
		} else {
			algo = "strie"
		}
		dict := &dict.Dict{
			Single:       conf.Single,
			Algorithm:    algo,
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
