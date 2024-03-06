package result

import (
	"fmt"
	"sort"

	"github.com/nopdan/gosmq/pkg/util"
)

type CodePosCount struct {
	Code  string
	Pos   int
	Count int
}

type WordCode struct {
	Word string
	Code string
}

// 匹配一段文字得到的信息
type MatchRes struct {
	TextIdx int // 文章索引
	DictIdx int // 码表索引

	PartIdx  int // 文章分段索引
	Segment  []WordCode
	segments []struct {
		// 分段索引
		PartIdx int
		// 每段的分词结果
		Segment []WordCode
	}

	// 每个词条对应的编码，以及出现的次数
	StatData map[string]*CodePosCount

	KeysDist  [128]int          // 按键分布
	NotHanMap map[rune]struct{} // 非汉字
	LackMap   map[rune]struct{} // 缺失的汉字

	CodeLenDist   []int   // 码长
	WordLenDist   []int   // 词长
	CollisionDist []int   // 选重
	Equivalent    float64 // 总当量

	Commit struct {
		Count     int // 上屏数
		Word      int // 打词数
		WordChars int // 打词字数
		WordFirst int // 首选词

		Collision      int // 选重
		CollisionChars int // 选重字数
	}

	Combs struct {
		Count       int // 按键组合数
		SameFingers int // 同指
		DoubleHit   int // 同键双击
		TribleHit   int // 同键三连击
		SingleSpan  int // 小跨排
		MultiSpan   int // 大跨排
		Staggered   int // 错手
		Disturb     int // 小指干扰
	}

	Hands struct {
		LL int
		LR int
		RL int
		RR int
	}
}

func NewMatchRes() *MatchRes {
	mRes := new(MatchRes)
	mRes.Segment = make([]WordCode, 0)
	mRes.segments = make([]struct {
		PartIdx int
		Segment []WordCode
	}, 0)
	mRes.StatData = make(map[string]*CodePosCount)

	mRes.NotHanMap = make(map[rune]struct{})
	mRes.LackMap = make(map[rune]struct{})
	mRes.CodeLenDist = make([]int, 0, 10)
	mRes.WordLenDist = make([]int, 0, 10)
	mRes.CollisionDist = make([]int, 0, 10)
	return mRes
}

// 将每次匹配得到的信息追加到总结果
func (m *MatchRes) Combine(mRes *MatchRes) {
	// 第一个 MatchRes 为总结果
	if len(m.segments) == 0 {
		m.segments = append(m.segments, struct {
			PartIdx int
			Segment []WordCode
		}{m.PartIdx, m.Segment})
	}
	if len(mRes.Segment) != 0 {
		m.segments = append(m.segments, struct {
			PartIdx int
			Segment []WordCode
		}{mRes.PartIdx, mRes.Segment})
	}
	for k, v := range mRes.StatData {
		if _, ok := m.StatData[k]; !ok {
			m.StatData[k] = v
		} else {
			m.StatData[k].Count += v.Count
		}
	}
	for i := 33; i < 128; i++ {
		m.KeysDist[i] += mRes.KeysDist[i]
	}
	for k := range mRes.NotHanMap {
		m.NotHanMap[k] = struct{}{}
	}
	for k := range mRes.LackMap {
		m.LackMap[k] = struct{}{}
	}
	for i := range mRes.CodeLenDist {
		util.AddTo(mRes.CodeLenDist[i], &m.CodeLenDist, i)
	}
	for i := range mRes.WordLenDist {
		util.AddTo(mRes.WordLenDist[i], &m.WordLenDist, i)
	}
	for i := range mRes.CollisionDist {
		util.AddTo(mRes.CollisionDist[i], &m.CollisionDist, i)
	}
	m.Equivalent += mRes.Equivalent

	m.Commit.Count += mRes.Commit.Count
	m.Commit.Word += mRes.Commit.Word
	m.Commit.WordChars += mRes.Commit.WordChars
	m.Commit.WordFirst += mRes.Commit.WordFirst
	m.Commit.Collision += mRes.Commit.Collision
	m.Commit.CollisionChars += mRes.Commit.CollisionChars

	m.Combs.Count += mRes.Combs.Count
	m.Combs.SameFingers += mRes.Combs.SameFingers
	m.Combs.DoubleHit += mRes.Combs.DoubleHit
	m.Combs.TribleHit += mRes.Combs.TribleHit
	m.Combs.SingleSpan += mRes.Combs.SingleSpan
	m.Combs.MultiSpan += mRes.Combs.MultiSpan
	m.Combs.Staggered += mRes.Combs.Staggered
	m.Combs.Disturb += mRes.Combs.Disturb

	m.Hands.LL += mRes.Hands.LL
	m.Hands.LR += mRes.Hands.LR
	m.Hands.RL += mRes.Hands.RL
	m.Hands.RR += mRes.Hands.RR
}

func (m *MatchRes) Print(detailed bool) {
	if m.segments != nil && detailed {
		sort.Slice(m.segments, func(i, j int) bool {
			return m.segments[i].PartIdx < m.segments[j].PartIdx
		})
		for i := range m.segments {
			for j := range m.segments[i].Segment {
				fmt.Printf("part: %d word: %s code: %s\n",
					m.segments[i].PartIdx, m.segments[i].Segment[j].Word, m.segments[i].Segment[j].Code)
			}
		}
	}

	for b, count := range m.KeysDist {
		if count != 0 {
			fmt.Printf("key: %s count: %d\n", string(rune(b)), count)
		}
	}

	notHan := ""
	for k := range m.NotHanMap {
		notHan += string(k)
	}
	fmt.Printf("not han: %s\n", notHan)
	lackHan := ""
	for k := range m.LackMap {
		lackHan += string(k)
	}
	fmt.Printf("lack han: %s\n", lackHan)

	fmt.Printf("code len dist: %v\n", m.CodeLenDist)
	fmt.Printf("word len dist: %v\n", m.WordLenDist)
	fmt.Printf("collision dist: %v\n", m.CollisionDist)
	fmt.Printf("equivalent: %f\n", m.Equivalent)

	fmt.Printf("commit: %+v\n", m.Commit)
	fmt.Printf("combs: %+v\n", m.Combs)
	fmt.Printf("hands: %+v\n", m.Hands)
}

func (m *MatchRes) ToResult() *Result {
	// TODO
	return nil
}
