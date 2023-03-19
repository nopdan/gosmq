package smq

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/imetool/goutil/util"
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
	resArr := t.Race([]*Dict{di})
	return resArr[0]
}

type wcIdx struct {
	idx     int
	wordSli []string
	codeSli []string
}

// 一篇文章计算多个码表
func (t *Text) Race(dicts []*Dict) []*Result {
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
				resArr[i].append(mRes, dicts[i])
				if dicts[i].Verbose {
					resArr[i].wcIdxs = append(resArr[i].wcIdxs, wcIdx{idx, mRes.wordSlice, mRes.codeSlice})
				}
				lock.Unlock()
				<-ch
			}(text, i, idx)
		}
		if err != nil {
			break
		}
	}
	wg.Wait()

	for i, dict := range dicts {
		if dict.Verbose {
			sort.Slice(resArr[i].wcIdxs, func(j, k int) bool {
				return resArr[i].wcIdxs[j].idx < resArr[i].wcIdxs[k].idx
			})
		}
	}

	for i := range dicts {
		resArr[i].stat()
	}
	return resArr
}

// 多篇文章计算多个码表
func Parallel(texts []string, dicts []*Dict) [][]*Result {
	resArr := make([][]*Result, len(texts))

	var wg sync.WaitGroup
	var lock sync.Mutex
	ch := make(chan struct{}, 16)
	for i, text := range texts {
		ch <- struct{}{}
		wg.Add(1)
		go func(text string, i int) {
			t := &Text{}
			t.Load(text)
			res := t.Race(dicts)
			lock.Lock()
			resArr[i] = res
			lock.Unlock()
			<-ch
			wg.Done()
		}(text, i)
	}
	wg.Wait()
	return resArr
}
