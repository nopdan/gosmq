package smq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/goutil/util"
)

type Smq struct {
	Name string // 文本名
	Text []byte // 文本
}

// 从文件添加文本
func (s *Smq) Load(path string) error {
	rd, err := util.Read(path)
	if err != nil {
		return err
	}
	s.Name = util.GetFileName(path)
	s.Text, _ = io.ReadAll(rd)
	fmt.Println("从文件初始化赛码器...", path)
	return nil
}

func (s *Smq) LoadString(name, text string) {
	if text != "" {
		fmt.Println("从字符串初始化赛码器...", name)
	}
	s.Name = name
	s.Text = []byte(text)
}

// 计算一个码表
func (smq *Smq) Eval(dict *dict.Dict) *Result {
	res := newResult()
	mRes := newMatchRes(10)
	brd := bufio.NewReader(bytes.NewReader(smq.Text))
	for {
		line, err := brd.ReadString('\n')
		text := []rune(line)
		mr := newMatchRes(len(text) / 3)
		mr.match(text, dict.Matcher, dict.Verbose, res)
		res.feel(mr.codes, dict)

		mRes.append(mr)
		if err != nil {
			break
		}
	}
	res.stat(mRes, dict)
	res.statFeel(dict)
	if dict.Verbose {
		OutputDetail(smq.Name, res, mRes)
	}
	return res
}

// 计算多个码表
func (smq *Smq) EvalDicts(dicts []*dict.Dict) []*Result {
	ret := make([]*Result, len(dicts))

	var wg sync.WaitGroup
	for i := range dicts {
		wg.Add(1)
		go func(j int) {
			ret[j] = smq.Eval(dicts[j])
			wg.Done()
		}(i)
	}
	wg.Wait()
	return ret
}
