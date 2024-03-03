package smq

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/nopdan/gosmq/pkg/util"
)

type Text struct {
	Name string // 文本名

	reader io.Reader // 文本
	bufLen int
}

// 从文件添加文本
func (t *Text) Load(path string) error {
	t.Name = util.GetFileName(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fi, _ := f.Stat()
	if fi.Size() < 4<<20 {
		// 4MB 以下 64KB
		t.bufLen = 64 << 10
	} else {
		// 其他 256KB
		t.bufLen = 256 << 10
	}
	// fmt.Println("buffer size", s.bufLen)
	t.reader = util.NewReader(f)
	return nil
}

func (t *Text) LoadString(name, text string) {
	t.Name = name
	t.reader = strings.NewReader(text)
}

// 计算一个码表
func (t *Text) RaceOne(di *Dict) *Result {
	resArr := t.Race([]*Dict{di}, true)
	return resArr[0]
}

type wcIdx struct {
	idx     int
	wordSli []string
	codeSli []string
}

// 一篇文章计算多个码表，是否输出
func (t *Text) Race(dicts []*Dict, output bool) []*Result {
	resArr := make([]*Result, len(dicts))
	for i := range dicts {
		resArr[i] = newResult()
		resArr[i].TextName = t.Name
		resArr[i].DictName = dicts[i].Name
		resArr[i].DictLen = dicts[i].length
		resArr[i].Single = dicts[i].Single
	}
	brd := bufio.NewReader(t.reader)

	var wg sync.WaitGroup
	var lock sync.Mutex
	ch := make(chan struct{}, 16)
	for idx := 0; ; idx++ {
		text, err := SplitStep(brd, t.bufLen)
		for i := range dicts {
			if dicts[i].length < 100 {
				continue
			}
			wg.Add(1)
			ch <- struct{}{}
			go func(text []byte, i, idx int) {
				defer wg.Done()
				mRes := match(text, dicts[i])
				mRes.dictIdx = i
				// 加锁操作
				lock.Lock()
				resArr[i].append(mRes, dicts[i], i)
				lock.Unlock()
				<-ch
			}(text, i, idx)
		}
		if err != nil {
			break
		}
	}
	wg.Wait()
	for i := range dicts {
		resArr[i].stat()
		if output {
			resArr[i].OutputSplit(dicts[i])
			resArr[i].OutputStat(dicts[i])
		}
	}
	return resArr
}

// 多篇文章计算多个码表，回调函数针对每篇文章生成的结果列表
func Parallel(texts []string, dicts []*Dict, callback func([]*Result)) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 16)
	for i, text := range texts {
		ch <- struct{}{}
		wg.Add(1)
		go func(text string, i int) {
			t := &Text{}
			t.Load(text)
			res := t.Race(dicts, true)
			callback(res)
			<-ch
			wg.Done()
		}(text, i)
	}
	wg.Wait()
}

// 多篇文章计算多个码表，合并同一个码表的多个结果
func ParallelMerge(texts []string, dicts []*Dict) []*Result {
	resArr := make([]*Result, len(dicts))
	// 合并结果不会输出分词
	for i, dict := range dicts {
		resArr[i] = newResult()
		resArr[i].TextName = "总计"
		resArr[i].DictName = dict.Name
		resArr[i].DictLen = dict.length
		resArr[i].Single = dict.Single
		dict.Split = false
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	ch := make(chan struct{}, 16)
	for i, text := range texts {
		ch <- struct{}{}
		wg.Add(1)
		go func(text string, i int) {
			t := &Text{}
			t.Load(text)
			res := t.Race(dicts, false)
			lock.Lock()
			mergeRes(resArr, res, dicts)
			lock.Unlock()
			<-ch
			wg.Done()
		}(text, i)
	}
	wg.Wait()
	for i, res := range resArr {
		res.stat()
		res.OutputStat(dicts[i])
	}
	return resArr
}
