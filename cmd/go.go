package cmd

import (
	"fmt"
	"time"

	"github.com/imetool/gosmq/pkg/smq"
)

var conf = &struct {
	Text []string // 文本
	Dict []string // 码表

	Single       bool   // 单字模式
	Algo         string // 匹配算法
	Stable       bool   // 按码表顺序(覆盖algo)
	PressSpaceBy string // 空格按键方式 left|right|both
	Clean        bool   // 只统计词库中的词条

	Verbose bool // 输出全部数据
	Split   bool // 输出分词数据
	Stat    bool // 输出词条数据
	Json    bool // 输出json数据
	HTML    bool // 保存 html 结果

	Hidden bool // 隐藏 cli 结果展示
	Merge  bool // 合并多文本的结果
}{}

func init() {
	rootCmd.Flags().StringArrayVarP(&conf.Text, "text", "t", nil, "文本文件或文件夹，可以为多个")
	rootCmd.Flags().StringArrayVarP(&conf.Dict, "dict", "i", nil, "码表文件或文件夹，可以为多个")

	rootCmd.Flags().BoolVarP(&conf.Single, "single", "s", false, "启用单字模式")
	rootCmd.Flags().BoolVarP(&conf.Stable, "stable", "", false, "按码表顺序")
	rootCmd.Flags().StringVarP(&conf.PressSpaceBy, "space", "k", "both", "空格按键方式 left|right|both")
	rootCmd.Flags().BoolVarP(&conf.Clean, "clean", "c", false, "只统计词库中的词条")

	rootCmd.Flags().BoolVarP(&conf.Verbose, "verbose", "v", false, "输出全部数据")
	rootCmd.Flags().BoolVarP(&conf.Split, "split", "", false, "输出分词数据")
	rootCmd.Flags().BoolVarP(&conf.Stat, "stat", "", false, "输出词条数据")
	rootCmd.Flags().BoolVarP(&conf.Json, "json", "", false, "输出 json 数据")
	rootCmd.Flags().BoolVarP(&conf.HTML, "html", "", false, "保存 html 结果")

	rootCmd.Flags().BoolVarP(&conf.Hidden, "hidden", "", false, "隐藏 cli 结果展示")
	rootCmd.Flags().BoolVarP(&conf.Merge, "merge", "m", false, "合并多文本的结果")
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
		conf.HTML = true
	}
	// 开始计时
	start := time.Now()
	texts := make([]string, 0, len(conf.Text))
	for _, v := range conf.Text {
		texts = append(texts, getFiles(v)...)
	}
	fmt.Println("载入文本：")
	for _, v := range texts {
		fmt.Println("-> ", v)
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
			Clean:        conf.Clean,
			Split:        conf.Split,
			Stat:         conf.Stat,
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
	textLenTotal := 0
	var printEnd = func() {
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
		fmt.Printf("共载入 %d 个码表，%d 个文本，总字数 %d，总耗时：%v\n", len(dicts), len(texts), textLenTotal, time.Since(start))
	}

	if conf.Merge {
		resArr := smq.ParallelMerge(texts, dicts)
		for _, res := range resArr {
			if !conf.Hidden {
				printSep()
				Output([]*smq.Result{res})
			}
			OutputHTML([]*smq.Result{res}, conf.HTML)
			OutPutJson(res, conf.Json)
			textLenTotal = res.TextLen
		}
		printEnd()
		return
	}

	smq.Parallel(texts, dicts, func(v []*smq.Result) {
		if len(v) == 0 {
			return
		}
		textLenTotal += v[0].TextLen
		if !conf.Hidden {
			printSep()
			Output(v)
		}
		OutputHTML(v, conf.HTML)
		for _, res := range v {
			OutPutJson(res, conf.Json)
		}
	})
	printEnd()
}
