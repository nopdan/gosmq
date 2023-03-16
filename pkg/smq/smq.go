package smq

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/goutil/util"
)

type Smq struct {
	Name   string    // 文本名
	reader io.Reader // 文本
	bufLen int64
}

// 从文件添加文本
func (s *Smq) Load(path string) error {
	s.Name = util.GetFileName(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fi, _ := f.Stat()
	if fi.Size() < 4<<20 {
		// 4MB 以下 64KB
		s.bufLen = 64 << 10
	} else {
		// 其他 256KB
		s.bufLen = 256 << 10
	}
	// fmt.Println("buffer size", s.bufLen)
	s.reader = util.NewReader(f)
	fmt.Println("从文件初始化赛码器...", path)
	return nil
}

func (s *Smq) LoadString(name, text string) {
	if text != "" {
		fmt.Println("从字符串初始化赛码器...", name)
	}
	s.Name = name
	s.reader = strings.NewReader(text)
	// s.Text = []byte(text)
}

// 计算一个码表
func (smq *Smq) Eval(dict *dict.Dict) *Result {
	res := newResult()
	mRes := newMatchRes(10)
	brd := bufio.NewReader(smq.reader)

	// for count := 0; ; count++ {
	for {
		var text []rune

		buffer := make([]byte, smq.bufLen)
		n, err := io.ReadFull(brd, buffer)
		buffer = buffer[:n]

		// 分割文本
		for {
			b, err := brd.ReadByte()
			// 控制字符 直接切分
			if b < 33 {
				text = []rune(string(buffer))
				break
			}
			// utf-8 前缀
			if b >= 0b11000000 {
				brd.UnreadByte()
			} else {
				buffer = append(buffer, b)
			}
			// EOF
			if err != nil {
				text = []rune(string(buffer))
				break
			}
			// 读到合法字符，开始读 rune
			if b < 128 || b >= 0b11000000 {
				text = []rune(string(buffer))
			OUT:
				// 超过限制读不到分割符直接 break
				for lim := int64(0); lim < smq.bufLen; lim++ {
					rn, _, err := brd.ReadRune()
					if rn < 33 {
						break
					}
					text = append(text, rn)
					if err != nil {
						break
					}
					switch rn {
					case '。', '？', '！', '》':
						break OUT
					}
				}
				break
			}
		}

		// f, _ := os.OpenFile("test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		// fmt.Fprintf(f, "\n-------------------------------- %d --------------------------------\n", count)
		// f.WriteString(string(text))

		mr := newMatchRes(len(text) / 3)
		mr.match(text, dict, res)
		mRes.append(mr)

		if err != nil {
			break
		}
	}
	res.stat(mRes, dict)
	res.statFeel(dict)

	OutputDetail(dict, smq.Name, res, mRes)
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
