package smq

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/result"
)

type Config struct {
	textList []*data.Text
	dictList []*data.Dict

	Merge bool // 合并同一码表的多个文章结果
	Clean bool // 只统计词库中的词条
	Split bool // 统计分词结果
	Stat  bool // 统计每个词条出现的次数
}

// 逻辑 CPU 数量（线程数）
var NUM_CPU = runtime.NumCPU()

func (c *Config) Reset() {
	c.textList = c.textList[:0]
	c.dictList = c.dictList[:0]
}

func (c *Config) AddText(textList ...*data.Text) {
	if c.textList == nil {
		c.textList = make([]*data.Text, 0)
	}
	for _, text := range textList {
		err := text.Init()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
	}
	c.textList = append(c.textList, textList...)
}

func (c *Config) AddDict(dictList ...*data.Dict) {
	if c.dictList == nil {
		c.dictList = make([]*data.Dict, 0)
	}
	for _, dict := range dictList {
		err := dict.Init()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
	}
	c.dictList = append(c.dictList, dictList...)
}

func (c *Config) Race() [][]*result.Result {
	if len(c.textList) == 0 || len(c.dictList) == 0 {
		fmt.Printf("输入有误\n")
		return nil
	}
	// 限制并发数量
	ch := make(chan *result.MatchRes, NUM_CPU)
	var wg sync.WaitGroup
	for i, text := range c.textList {
		// 分段计算当前文章，pIdx 为每一段的索引
		pIdx := -1
		for {
			text, err := text.Iter()
			pIdx++
			for j, dict := range c.dictList {
				wg.Add(1)
				go func(i, j, pIdx int) {
					defer wg.Done()
					mRes := c.match(text, dict)
					mRes.TextIdx = i
					mRes.DictIdx = j
					mRes.PartIdx = pIdx
					ch <- mRes
				}(i, j, pIdx)
			}
			if err != nil {
				break
			}
		}
	}

	// 文章数量和码表数量
	var tNum, dNum = len(c.textList), len(c.dictList)
	mRes := make([][]*result.MatchRes, tNum)
	for i := range tNum {
		mRes[i] = make([]*result.MatchRes, dNum)
		for j := range dNum {
			mRes[i][j] = result.NewMatchRes()
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// 循环从 ch 通道中接受值
	for part := range ch {
		mRes[part.TextIdx][part.DictIdx].Combine(part)
	}

	if c.Merge {
		for j := range dNum {
			for i := range tNum {
				if i != 0 {
					mRes[0][j].Combine(mRes[i][j])
				}
			}
		}
		res := make([]*result.Result, dNum)
		for j := range dNum {
			info := &result.Info{
				TextName: "合并结果",
				DictName: c.dictList[j].Text.Name,
				DictLen:  c.dictList[j].Length,
				Single:   c.dictList[j].Single,
			}
			res[j] = mRes[0][j].Stat(info)
			if c.Split {
				res[j].OutputSplit()
			}
			if c.Stat {
				res[j].OutputStat()
			}
		}
		return [][]*result.Result{res}
	}

	res := make([][]*result.Result, tNum)
	for i := range tNum {
		res[i] = make([]*result.Result, dNum)
		for j := range dNum {
			info := &result.Info{
				TextName: c.textList[i].Name,
				DictName: c.dictList[j].Text.Name,
				DictLen:  c.dictList[j].Length,
				Single:   c.dictList[j].Single,
			}
			res[i][j] = mRes[i][j].Stat(info)
			if c.Split {
				res[i][j].OutputSplit()
			}
			if c.Stat {
				res[i][j].OutputStat()
			}
		}
	}
	return res
}
