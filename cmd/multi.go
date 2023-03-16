package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
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
	Texts  []string // 文本
	Folder string   // 从文件夹读取文本
	Dict   string   // 码表
	Basic
}{}

func init() {
	multiCmd.PersistentFlags().StringArrayVarP(&multi.Texts, "text", "t", nil, "文本路径")
	multiCmd.PersistentFlags().StringVarP(&multi.Folder, "folder", "f", "", "从文件夹读取文本")
	multiCmd.PersistentFlags().StringVarP(&multi.Dict, "dict", "i", "", "码表路径")
	multiCmd.PersistentFlags().BoolVarP(&multi.Single, "single", "s", false, "启用单字模式")
	multiCmd.PersistentFlags().StringVarP(&multi.Algo, "algo", "a", "strie", "匹配算法(strie|trie)")
	multiCmd.PersistentFlags().StringVarP(&multi.PressSpaceBy, "space", "p", "both", "空格按键方式 left|right|both")
	multiCmd.PersistentFlags().BoolVarP(&multi.Verbose, "verbose", "v", false, "输出详细数据")
	multiCmd.PersistentFlags().BoolVarP(&multi.Split, "split", "", false, "输出分词数据")
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
		fmt.Printf("载入 %s 下的文本: \n", multi.Folder)
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
		Stat:         multi.Verbose,
		Json:         multi.Verbose,
		Split:        multi.Split,
	}
	dict.Load(multi.Dict)
	fmt.Println("载入码表：", dict.Name)

	printSep()
	textTotalLen := int64(0)
	var wg sync.WaitGroup
	ch := make(chan struct{}, 8)
	for _, text := range multi.Texts {
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
			res := s.Eval(dict)
			atomic.AddInt64(&textTotalLen, int64(res.Basic.TextLen))
			if multi.Folder == "" {
				fmt.Println("载入文本：", s.Name)
			}
			fmt.Printf("此文本耗时：%v\n", time.Since(mid))
			printSep()
			Output([]*smq.Result{res}, s.Name)
			<-ch
			wg.Done()
		}(text)
	}
	wg.Wait()

	fmt.Printf("共载入 %d 个文本，总字数 %d，总耗时：%v\n", len(multi.Texts), textTotalLen, time.Since(start))
}
