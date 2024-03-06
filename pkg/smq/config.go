package smq

import (
	"runtime"
	"sync"

	"github.com/nopdan/gosmq/pkg/dict"
	"github.com/nopdan/gosmq/pkg/result"
	"github.com/nopdan/gosmq/pkg/text"
)

type Config struct {
	textList []*text.Text
	dictList []*dict.Dict

	Clean bool // 只统计词库中的词条
	Split bool // 统计分词结果
	Stat  bool // 统计每个词条出现的次数
}

// 逻辑 CPU 数量（线程数）
var NUM_CPU = runtime.NumCPU()

type SmqOption func(*Config)

func New(opts ...SmqOption) *Config {
	c := new(Config)
	c.textList = make([]*text.Text, 0)
	c.dictList = make([]*dict.Dict, 0)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithClean() SmqOption {
	return func(c *Config) {
		c.Clean = true
	}
}

func WithSplit() SmqOption {
	return func(c *Config) {
		c.Split = true
	}
}

func WithStat() SmqOption {
	return func(c *Config) {
		c.Stat = true
	}
}

func (c *Config) AddText(textList ...*text.Text) {
	c.textList = append(c.textList, textList...)
}

func (c *Config) AddDict(dictList ...*dict.Dict) {
	c.dictList = append(c.dictList, dictList...)
}

func (c *Config) Race() [][]*result.Result {
	// 限制并发数量
	ch := make(chan *result.MatchRes, NUM_CPU)
	var wg sync.WaitGroup
	for i, text := range c.textList {
		// 分段计算当前文章，pIdx 为每一段的索引
		pIdx := 0
		for {
			text, err := text.Iter()
			// fmt.Println(util.UnsafeToString(text))
			if err != nil {
				break
			}
			for j, dict := range c.dictList {
				wg.Add(1)
				// go 1.22 修复了 range 循环问题
				go func() {
					defer wg.Done()
					mRes := c.match(text, dict)
					mRes.TextIdx = i
					mRes.DictIdx = j
					mRes.PartIdx = pIdx
					ch <- mRes
				}()
			}
			pIdx++
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

	res := make([][]*result.Result, tNum)
	for i := range tNum {
		res[i] = make([]*result.Result, dNum)
		for j := range dNum {
			// TODO
			mRes[i][j].Print(false)
			res[i][j] = mRes[i][j].ToResult()
		}
	}
	return res
}
