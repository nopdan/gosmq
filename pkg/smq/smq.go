package smq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/cxcn/gosmq/pkg/dict"
	"github.com/cxcn/gosmq/pkg/util"
)

type Smq struct {
	Name   string       // 文本名
	Text   []byte       // 文本
	Inputs []*dict.Dict // 码表选项
}

// 初始化一个赛码器
func New(name, path string) Smq {
	rd, err := util.Read(path)
	if err != nil {
		panic("Error! 读取文件失败：" + path)
	}
	if name == "" {
		name = util.GetFileName(path)
	}
	text, _ := io.ReadAll(rd)
	fmt.Println("从文件初始化赛码器...", name)
	return Smq{name, text, []*dict.Dict{}}
}

func NewFromString(name, text string) Smq {
	if text != "" {
		fmt.Println("从字符串初始化赛码器...", name)
	}
	return Smq{name, []byte(text), []*dict.Dict{}}
}

// 添加一个码表
func (smq *Smq) Add(dict *dict.Dict) {
	// 合法输入
	if dict.Legal {
		dict.Init()
		smq.Inputs = append(smq.Inputs, dict)
		fmt.Println("添加了一个码表：", dict.Name)
	}
}

// 开始计算
func (smq *Smq) Run() []*Result {
	smqLen := len(smq.Inputs)
	ret := make([]*Result, 0, smqLen)
	for i := 0; i < len(smq.Inputs); i++ {
		ret = append(ret, newResult())
		// fmt.Println(smq.Inputs[i])
	}

	var wg sync.WaitGroup
	for i := range smq.Inputs {
		wg.Add(1)
		go func(j int) {
			res, dict := ret[j], smq.Inputs[j]
			brd := bufio.NewReader(bytes.NewReader(smq.Text))
			for {
				line, err := brd.ReadString('\n')
				codes := res.match([]rune(line), dict)
				res.feel(codes, dict)
				if err != nil {
					break
				}
			}
			res.stat(dict)
			if dict.OutputDetail {
				OutputDetail(smq.Name, res)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return ret
}
