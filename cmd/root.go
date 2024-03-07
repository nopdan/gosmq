package cmd

import (
	"fmt"
	"time"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/smq"
)

var conf = &struct {
	Texts []string // 文本
	Dicts []string // 码表

	Single    bool   // 单字模式
	Stable    bool   // 按码表顺序
	SpacePref string // 空格按键方式 left|right|both
	Clean     bool   // 只统计词库中的词条
	Split     bool   // 统计分词结果
	Stat      bool   // 统计每个词条出现的次数

	Verbose bool // 输出全部数据
	Json    bool // 输出json数据
	HTML    bool // 保存 html 结果

	Hidden bool // 隐藏 cli 结果展示
	Merge  bool // 合并多文本的结果
}{}

func init() {
	rootCmd.Flags().StringArrayVarP(&conf.Texts, "text", "t", nil, "文本文件或文件夹，可以为多个")
	rootCmd.Flags().StringArrayVarP(&conf.Dicts, "dict", "i", nil, "码表文件或文件夹，可以为多个")

	rootCmd.Flags().BoolVarP(&conf.Single, "single", "s", false, "启用单字模式")
	rootCmd.Flags().BoolVarP(&conf.Stable, "stable", "", false, "按码表顺序")
	rootCmd.Flags().StringVarP(&conf.SpacePref, "space", "k", "both", "空格按键方式 left|right|both")
	rootCmd.Flags().BoolVarP(&conf.Clean, "clean", "c", false, "只统计词库中的词条")

	rootCmd.Flags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出全部数据")
	rootCmd.Flags().BoolVarP(&conf.Split, "split", "", false, "输出分词数据")
	rootCmd.Flags().BoolVarP(&conf.Stat, "stat", "", false, "输出词条数据")
	rootCmd.Flags().BoolVarP(&conf.Json, "json", "", false, "输出 json 数据")
	rootCmd.Flags().BoolVarP(&conf.HTML, "html", "", false, "保存 html 结果")

	rootCmd.Flags().BoolVarP(&conf.Hidden, "hidden", "", false, "隐藏 cli 结果展示")
	rootCmd.Flags().BoolVarP(&conf.Merge, "merge", "m", false, "合并多文本的结果")
}

func _root() {
	if len(conf.Dicts) == 0 || len(conf.Texts) == 0 {
		fmt.Println("输入有误")
		return
	}
	if conf.Verbose {
		conf.Split = true
		conf.Stat = true
		conf.Json = true
		conf.HTML = true
	}

	smq := &smq.Config{
		Merge: conf.Merge,
		Clean: conf.Clean,
		Split: conf.Split,
		Stat:  conf.Stat,
	}

	// 开始计时
	start := time.Now()
	texts := make([]string, 0, len(conf.Texts))
	for _, v := range conf.Texts {
		for _, v := range getFiles(v) {
			texts = append(texts, v)
			smq.AddText(&data.Text{Path: v})
		}
	}

	if !conf.Hidden {
		fmt.Println("载入文本：")
		for _, v := range texts {
			fmt.Println("-> ", v)
		}
		fmt.Println()
	} else {
		fmt.Printf("载入 %d 个文本\n", len(texts))
	}

	dictNames := make([]string, 0, len(conf.Dicts))
	for _, v := range conf.Dicts {
		dictNames = append(dictNames, getFiles(v)...)
	}
	fmt.Println("载入码表：")
	dictStartTime := time.Now()
	mid := time.Now()
	for _, v := range dictNames {
		dict := &data.Dict{
			Text:      &data.Text{Path: v},
			Single:    conf.Single,
			SpacePref: conf.SpacePref,
		}
		if conf.Stable {
			dict.Algorithm = "ordered"
		}
		smq.AddDict(dict)

		if len(dictNames) == 1 {
			fmt.Println("=> ", v)
		} else {
			fmt.Println("=> ", v, "\t耗时：", time.Since(mid))
			mid = time.Now()
		}
	}
	fmt.Printf("载入码表耗时：%v\n\n", time.Since(dictStartTime))

	// race
	fmt.Println("比赛开始...")
	var printEnd = func(textLen int) {
		if conf.Split {
			if conf.Merge {
				fmt.Println("--merge 不会输出分词结果。")
			} else {
				fmt.Println("已输出分词结果")
			}
		}
		if conf.Stat {
			fmt.Println("已输出词条统计数据")
		}
		if conf.HTML {
			fmt.Println("已保存 html 结果")
		}
		if conf.Json {
			fmt.Println("已输出 json 数据")
		}
		fmt.Printf("共载入 %d 个码表，%d 个文本，总字数 %d，总耗时：%v\n", len(dictNames), len(texts), textLen, time.Since(start))
	}

	textLen := 0
	res := smq.Race()
	for _, v := range res {
		Output(v)
		for _, vv := range v {
			textLen += vv.Info.TextLen
		}
	}

	printEnd(textLen)
}
