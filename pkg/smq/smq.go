package smq

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/goutil/util"
)

type Smq struct {
	Name   string    // 文本名
	reader io.Reader // 文本
	bufLen int
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
func (smq *Smq) Eval(di *dict.Dict) *Result {
	resArr := smq.EvalDicts([]*dict.Dict{di})
	return resArr[0]
}

type wcIdx struct {
	idx     int
	wordSli []string
	codeSli []string
}

// 计算多个码表
func (smq *Smq) EvalDicts(dicts []*dict.Dict) []*Result {
	resArr := make([]*Result, len(dicts))
	for i := range dicts {
		resArr[i] = newResult()
	}
	brd := bufio.NewReader(smq.reader)

	var wg sync.WaitGroup
	var lock sync.Mutex
	ch := make(chan struct{}, 16)
	for idx := 0; ; idx++ {
		text, err := SplitStep(brd, smq.bufLen)
		for i := range dicts {
			wg.Add(1)
			ch <- struct{}{}
			go func(text []rune, i, idx int) {
				defer wg.Done()
				mRes := match(text, dicts[i])
				mRes.dictIdx = i
				// 加锁操作
				lock.Lock()
				resArr[i].append(mRes, dicts[i])
				if dicts[i].Split {
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
		if dict.Split {
			sort.Slice(resArr[i].wcIdxs, func(j, k int) bool {
				return resArr[i].wcIdxs[j].idx < resArr[i].wcIdxs[k].idx
			})
		}
	}

	for i, dict := range dicts {
		resArr[i].stat(dict)
		resArr[i].statFeel(dict)
		OutputDetail(dict, smq.Name, resArr[i])
	}

	return resArr
}

// 将每次匹配得到的信息追加到总结果
func (res *Result) append(mRes *matchRes, dict *dict.Dict) {
	if dict.Stat {
		for k, v := range mRes.statData {
			if _, ok := res.statData[k]; !ok {
				res.statData[k] = v
			} else {
				res.statData[k].Count += v.Count
			}
		}
	}
	res.Basic.TextLen += mRes.TextLen
	res.Basic.Commits += mRes.Commits
	res.Basic.NotHanCount += mRes.NotHanCount
	res.Basic.LackCount += mRes.LackCount

	res.Words.Commits.Count += mRes.WordsCommitsCount
	res.Words.Chars.Count += mRes.WordsCharsCount
	res.Words.FirstCount += mRes.WordsFirstCount
	res.Collision.Commits.Count += mRes.CollisionCommitsCount
	res.Collision.Chars.Count += mRes.CollisionCharsCount

	res.toTalEq10 += mRes.toTalEq10
	res.Combs.Count += mRes.CombsCount
	res.Fingers.Same.Count += mRes.SameFingers

	res.Hands.LL.Count += mRes.Hands.LL
	res.Hands.LR.Count += mRes.Hands.LR
	res.Hands.RL.Count += mRes.Hands.RL
	res.Hands.RR.Count += mRes.Hands.RR

	res.Combs.DoubleHit.Count += mRes.Combs.DoubleHit
	res.Combs.TribleHit.Count += mRes.Combs.TribleHit
	res.Combs.SingleSpan.Count += mRes.Combs.SingleSpan
	res.Combs.MultiSpan.Count += mRes.Combs.MultiSpan
	res.Combs.LongFingersDisturb.Count += mRes.Combs.LongFingersDisturb
	res.Combs.LittleFingersDisturb.Count += mRes.Combs.LittleFingersDisturb

	for k, v := range mRes.mapKeys {
		res.mapKeys[k] += v
	}
	for k := range mRes.notHanMap {
		res.notHanMap[k] = struct{}{}
	}
	for k := range mRes.lackMap {
		res.lackMap[k] = struct{}{}
	}

	for i, v := range mRes.CodeLenDist {
		AddToVal(&res.CodeLen.Dist, i, v)
	}

	for i, v := range mRes.WordsDist {
		AddToVal(&res.Words.Dist, i, v)
	}
	for i, v := range mRes.CollisionDist {
		AddToVal(&res.Collision.Dist, i, v)
	}
}
