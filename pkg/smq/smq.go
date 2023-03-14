package smq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/goutil/util"
)

type Smq struct {
	Name   string       // 文本名
	Text   []byte       // 文本
	Inputs []*dict.Dict // 码表选项
}

// 从文件添加文本
func (s *Smq) Load(path string) {
	rd, err := util.Read(path)
	if err != nil {
		log.Panic("Error! 读取文件失败：", path)
	}
	s.Name = util.GetFileName(path)
	s.Text, _ = io.ReadAll(rd)
	fmt.Println("从文件初始化赛码器...", path)
}

func (s *Smq) LoadString(name, text string) {
	if text != "" {
		fmt.Println("从字符串初始化赛码器...", name)
	}
	s.Name = name
	s.Text = []byte(text)
}

// 添加一个码表
func (smq *Smq) Add(dict *dict.Dict) {
	smq.Inputs = append(smq.Inputs, dict)
	fmt.Println("添加了一个码表：", dict.Name)
}

// 开始计算
func (smq *Smq) Run() []*Result {
	smqLen := len(smq.Inputs)
	ret := make([]*Result, smqLen)
	for i := 0; i < smqLen; i++ {
		ret[i] = newResult()
		// fmt.Println(smq.Inputs[i])
	}

	var wg sync.WaitGroup
	matchResults := make([]*matchRes, smqLen)
	for j := 0; j < smqLen; j++ {
		mr := newMatchRes(10)
		matchResults[j] = mr
	}

	for i := range smq.Inputs {
		wg.Add(1)
		go func(j int) {
			res, dict := ret[j], smq.Inputs[j]
			brd := bufio.NewReader(bytes.NewReader(smq.Text))
			for {
				line, err := brd.ReadString('\n')
				text := []rune(line)
				mr := newMatchRes(len(text) / 3)
				mr.match(text, dict.Matcher, dict.Verbose, res)
				res.feel(mr.codes, dict)

				matchResults[j].append(mr)
				if err != nil {
					break
				}
			}
			res.stat(matchResults[j], dict)
			res.statFeel(dict)
			if dict.Verbose {
				OutputDetail(smq.Name, res, matchResults[j])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return ret
}
