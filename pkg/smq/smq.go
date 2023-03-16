package smq

import (
	"bufio"
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
	return nil
}

func (s *Smq) LoadString(name, text string) {
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
		text, err := SplitStep(brd, smq.bufLen)
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
	resArr := make([]*Result, len(dicts))
	mResArr := make([]*matchRes, len(dicts))
	for i := range dicts {
		resArr[i] = newResult()
		mResArr[i] = newMatchRes(10)
	}

	brd := bufio.NewReader(smq.reader)
	var wg sync.WaitGroup
	// var lock sync.Mutex
	for {
		text, err := SplitStep(brd, smq.bufLen)
		for i := range dicts {
			wg.Add(1)
			go func(j int) {
				mr := newMatchRes(len(text) / 3)
				// lock.Lock()
				mr.match(text, dicts[j], resArr[j])
				mResArr[j].append(mr)
				// lock.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()
		if err != nil {
			break
		}
	}
	for i, dict := range dicts {
		resArr[i].stat(mResArr[i], dict)
		resArr[i].statFeel(dict)
		OutputDetail(dict, smq.Name, resArr[i], mResArr[i])
	}

	return resArr
}
