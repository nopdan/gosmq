package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
		err := filepath.Walk(multi.Folder, func(path string, info os.FileInfo, err error) error {
			multi.Texts = append(multi.Texts, path)
			return nil
		})
		multi.Texts = multi.Texts[1:]
		if len(multi.Texts) == 0 {
			fmt.Println(multi.Folder, "目录下没有文件！")
			return
		}
		fmt.Printf("载入 %s 下的文件: \n", multi.Folder)
		for _, v := range multi.Texts {
			fmt.Printf("-> %s\n", v)
		}
		if err != nil {
			panic(err)
		}
	}

	dict := &dict.Dict{
		Single:       multi.Single,
		Algorithm:    multi.Algo,
		PressSpaceBy: multi.PressSpaceBy,
		Verbose:      multi.Verbose,
	}
	dict.Load(multi.Dict)

	for _, text := range multi.Texts {
		// 初始化赛码器
		s := &smq.Smq{}
		err := s.Load(text)
		if err != nil {
			fmt.Println("Error! 读取文件失败：", err)
			continue
		}
		res := s.Eval(dict)
		fmt.Println("----------------------")
		Output([]*smq.Result{res}, s.Name)
	}
}
