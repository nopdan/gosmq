package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/gosmq/pkg/smq"
	"github.com/spf13/cobra"
)

var multiCmd = &cobra.Command{
	Use:   "multi",
	Short: "多文本测评",
	Run: func(cmd *cobra.Command, args []string) {
		multiCli()
	},
}
var multi = &struct {
	Texts        []string // 文本
	Folder       string   // 从文件夹读取文本
	Dict         string   // 码表
	Single       bool     // 单字模式
	Algo         string   // 匹配算法
	PressSpaceBy string   // 空格按键方式 left|right|both
	Verbose      bool     // 输出详细数据
}{}

func init() {
	multiCmd.PersistentFlags().StringArrayVarP(&multi.Texts, "text", "t", nil, "文本路径")
	multiCmd.PersistentFlags().StringVarP(&multi.Folder, "folder", "f", "", "从文件夹读取文本")
	multiCmd.PersistentFlags().StringVarP(&multi.Dict, "dict", "i", "", "码表路径")
	multiCmd.PersistentFlags().BoolVarP(&multi.Single, "single", "s", false, "启用单字模式")
	multiCmd.PersistentFlags().StringVarP(&multi.Algo, "algo", "a", "strie", "匹配算法(strie|trie)")
	multiCmd.PersistentFlags().StringVarP(&multi.PressSpaceBy, "space", "p", "both", "空格按键方式 left|right|both")
	multiCmd.PersistentFlags().BoolVarP(&multi.Verbose, "verbose", "v", false, "输出详细数据")
}

func multiCli() {

	start := time.Now()
	if multi.Dict == "" {
		fmt.Println("没有输入码表")
		return
	}
	if len(multi.Texts) == 0 && multi.Folder == "" {
		fmt.Println("没有输入文本")
		return
	}

	if multi.Folder != "" {
		multi.Texts = make([]string, 0)

		files, err := os.ReadDir(multi.Folder)
		if err != nil {
			panic(err)
		}
		fmt.Printf("载入 %s 下的文件: \n", multi.Folder)
		if !strings.HasSuffix(multi.Folder, "\\") {
			multi.Folder += "\\"
		}

		for _, file := range files {
			if !file.IsDir() {
				multi.Texts = append(multi.Texts, multi.Folder+file.Name())
				fmt.Printf("-> %s\n", file.Name())
			}
		}
	}

	dict := &dict.Dict{
		Single:       multi.Single,
		Algorithm:    multi.Algo,
		PressSpaceBy: multi.PressSpaceBy,
		Verbose:      multi.Verbose,
	}
	dict.Load(multi.Dict)

	printSep()
	textTotalLen := 0
	for _, text := range multi.Texts {
		mid := time.Now()
		// 初始化赛码器
		s := &smq.Smq{}
		err := s.Load(text)
		if err != nil {
			fmt.Println("Error! 读取文件失败：", err)
			continue
		}
		res := s.Eval(dict)
		textTotalLen += res.Basic.TextLen
		fmt.Printf("该文本耗时：%v\n", time.Since(mid))
		printSep()
		Output([]*smq.Result{res}, s.Name)
	}

	fmt.Printf("共载入 %d 个文本，总字数 %d，总耗时：%v\n", len(multi.Texts), textTotalLen, time.Since(start))
}
