package smq

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"unicode"

	"github.com/cxcn/gosmq/pkg/data"
)

var (
	KEY_POS = data.GetKeyPos()
	PUNCTS  = data.GetPuncts()
	COMB    = data.GetComb()
)

type Smq struct {
	Name   string    // 文本名
	Text   io.Reader // 文本
	Inputs []*Dict   // 码表选项
}

// 初始化一个赛码器
func New(name string, rd io.Reader) Smq {
	fmt.Println("从字节流初始化赛码器...")
	nrd := Tranformer(rd)
	fmt.Println("文本名：", name)
	return Smq{name, nrd, []*Dict{}}
}

func NewFromString(name, text string) Smq {
	fmt.Println("从字符串初始化赛码器...")
	rd := readFromString(text)
	fmt.Println("文本名：", name)
	return Smq{name, rd, []*Dict{}}
}

func NewFromPath(name, path string) Smq {
	fmt.Println("从文件初始化赛码器...")
	rd, err := readFromPath(path)
	if err != nil {
		log.Panicln("Error! 从文件初始化赛码器，路径：", path)
	}
	if name == "" {
		name = GetFileName(path)
	}
	fmt.Println("文本名：", name)
	return Smq{name, rd, []*Dict{}}
}

// 添加一个码表
func (smq *Smq) Add(dict *Dict) {
	// 合法输入
	if dict.legal {
		dict.init()
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
	brd := bufio.NewReader(smq.Text)
	if smqLen == 1 {
		for {
			line, err := brd.ReadString('\n')
			codes := ret[0].match([]rune(line), smq.Inputs[0])
			ret[0].feel(codes, smq.Inputs[0])
			if err != nil {
				break
			}
		}
	} else {
		var wg sync.WaitGroup
		// 逐行读取文本文件
		for {
			line, err := brd.ReadString('\n')
			for i, v := range smq.Inputs {
				wg.Add(1)
				tmp := ret[i]
				go func(arg *Dict) {
					codes := tmp.match([]rune(line), arg)
					tmp.feel(codes, arg)
					wg.Done()
				}(v)
			}
			wg.Wait()

			if err != nil {
				break
			}
		}
	}

	for i, v := range ret {
		v.stat(smq.Inputs[i])
	}
	return ret
}

func (res *Result) match(text []rune, dict *Dict) string {
	var sb strings.Builder
	sb.Grow(len(text))
	res.Basic.TextLen += len(text)
	for p := 0; p < len(text); {
		// 删掉空白字符
		switch text[p] {
		case 65533, '\n', '\r', '\t', ' ', '　':
			p++
			res.Basic.TextLen--
			continue
		}
		// 非汉字
		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			res.Basic.NotHanCount++
			res.mapNotHan[text[p]] = struct{}{}
		}

		i, code, order := dict.Matcher.Match(text, p)
		// 缺字
		if i == 0 {
			if isHan {
				res.Basic.LackCount++
				res.mapLack[text[p]] = struct{}{}
			}
			sb.WriteByte(' ')
			p++
			continue
		}

		sb.WriteString(code)
		AddTo(&res.wordsDist, i) // 词长分布
		if order != 1 {
			res.Collision.Chars.Count += i // 选重字数
		} else if i != 1 {
			res.Words.FirstCount++ // 首选词
		}
		AddTo(&res.collDist, order)     // 选重分布
		AddTo(&res.codeDist, len(code)) // 码长分布

		if dict.Details {
			word := string(text[p : p+i])
			res.Data.WordSlice = append(res.Data.WordSlice, word)
			res.Data.CodeSlice = append(res.Data.CodeSlice, code)
			if _, ok := res.Data.Details[word]; !ok {
				res.Data.Details[word] = &CoC{Code: code, Order: i}
			}
			res.Data.Details[word].Count++
		}
		p += i
	}
	return sb.String()
}

// 输出 json
func (smq *Smq) ToJson() ([]byte, error) {
	res := smq.Run()
	return ResToJson(res)
}

// 输出 json
func ResToJson(res []*Result) ([]byte, error) {
	return json.Marshal(res)
}
