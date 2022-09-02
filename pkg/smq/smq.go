package smq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
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
	Name   string  // 文本名
	Text   []byte  // 文本
	Inputs []*Dict // 码表选项
}

// 初始化一个赛码器
func New(name string, rd io.Reader) Smq {
	fmt.Println("从字节流初始化赛码器...")
	nrd := Tranformer(rd)
	text, _ := io.ReadAll(nrd)
	return Smq{name, text, []*Dict{}}
}

func NewFromString(name, text string) Smq {
	if text != "" {
		fmt.Println("从字符串初始化赛码器...")
	}
	return Smq{name, []byte(text), []*Dict{}}
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
	text, _ := io.ReadAll(rd)
	return Smq{name, text, []*Dict{}}
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

		if dict.OutputDetail {
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

func OutputDetail(textName string, res *Result) {
	// 创建文件夹
	os.Mkdir("result", 0666)
	// 输出分词结果
	var buf strings.Builder
	for i := 0; i < len(res.Data.CodeSlice); i++ {
		buf.WriteString(fmt.Sprintf("%s\t%s\n", res.Data.WordSlice[i], string(res.Data.CodeSlice[i])))
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_分词结果.txt", textName, res.Name), []byte(buf.String()), 0666)
	// 输出词条数据
	buf.Reset()
	buf.WriteString("词条\t编码\t顺序\t次数\n")
	type details struct {
		CoC
		word string
	}
	tmp := make([]details, 0, len(res.Data.Details))
	for k, v := range res.Data.Details {
		tmp = append(tmp, details{*v, k})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Count > tmp[j].Count
	})
	for _, v := range tmp {
		buf.WriteString(v.word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Order))
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Count))
		buf.WriteByte('\n')
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_词条数据.txt", textName, res.Name), []byte(buf.String()), 0666)
	// 输出 json 数据
	res.Data.CodeSlice = []string{}
	res.Data.WordSlice = []string{}
	res.Data.Details = make(map[string]*CoC)
	tmp2, _ := json.MarshalIndent(res, "", "  ")
	os.WriteFile(fmt.Sprintf("result/%s_%s.json", textName, res.Name), tmp2, 0666)
}
