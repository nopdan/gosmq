package result

import (
	"strings"
	"sync"

	"github.com/nopdan/gosmq/pkg/util"
)

type CodePosCount struct {
	Code  string
	Pos   int
	Count int
}

type Segment struct {
	// 分段索引
	PartIdx int
	// 每段的分词结果
	Builder *strings.Builder
}

// 匹配一段文字得到的信息
type MatchRes struct {
	PartIdx  int // 分段索引
	Segment  Segment
	segments []Segment

	// 每个词条对应的编码，以及出现的次数
	StatData map[string]*CodePosCount

	TextLen int // 文章总字数，不管有没有匹配到
	Dist    dist
	Commit  commit
	Char    char
	Pair    pair

	lock sync.Mutex
}

type dist struct {
	Key     [128]int     // 按键分布
	NotHan  map[rune]int // 非汉字
	LackHan map[rune]int // 缺字

	CodeLen   []int // 码长
	WordLen   []int // 词长
	Collision []int // 选重
}

type commit struct {
	Count     int // 上屏数
	Word      int // 打词数
	WordFirst int // 首选词
	Collision int // 选重
}

type char struct {
	Count     int // 总匹配字数
	Word      int // 打词字数
	WordFirst int // 首选词字数
	Collision int // 选重字数
}

type pair struct {
	Count      int     // 按键组合数
	Equivalent float64 // 总当量

	SameFinger int // 同手同指
	DoubleHit  int // 同键双击
	TribleHit  int // 同键三连击
	SingleSpan int // 小跨排
	MultiSpan  int // 大跨排
	Staggered  int // 错手
	Disturb    int // 小指干扰

	LeftToLeft   int // 左左
	LeftToRight  int // 左右
	RightToLeft  int // 右左
	RightToRight int // 右右

	DiffFinger int // 同手不同指
	SameHand   int // 同手
	DiffHand   int // 不同手
}

func NewMatchRes() *MatchRes {
	mRes := new(MatchRes)
	mRes.segments = make([]Segment, 0)
	mRes.StatData = make(map[string]*CodePosCount)

	mRes.Dist.NotHan = make(map[rune]int)
	mRes.Dist.LackHan = make(map[rune]int)
	mRes.Dist.CodeLen = make([]int, 0, 10)
	mRes.Dist.WordLen = make([]int, 0, 10)
	mRes.Dist.Collision = make([]int, 0, 10)
	return mRes
}

// 将每次匹配得到的信息追加到总结果
func (m *MatchRes) Combine(mRes *MatchRes) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if mRes.Segment.Builder != nil {
		m.segments = append(m.segments, mRes.Segment)
	}
	for k, v := range mRes.StatData {
		if _, ok := m.StatData[k]; !ok {
			m.StatData[k] = v
		} else {
			m.StatData[k].Count += v.Count
		}
	}
	for i := 33; i < 128; i++ {
		m.Dist.Key[i] += mRes.Dist.Key[i]
	}
	for k, v := range mRes.Dist.NotHan {
		m.Dist.NotHan[k] += v
	}
	for k, v := range mRes.Dist.LackHan {
		m.Dist.LackHan[k] += v
	}
	for i := range mRes.Dist.CodeLen {
		util.AddTo(mRes.Dist.CodeLen[i], &m.Dist.CodeLen, i)
	}
	for i := range mRes.Dist.WordLen {
		util.AddTo(mRes.Dist.WordLen[i], &m.Dist.WordLen, i)
	}
	for i := range mRes.Dist.Collision {
		util.AddTo(mRes.Dist.Collision[i], &m.Dist.Collision, i)
	}

	m.TextLen += mRes.TextLen
	m.Commit.Count += mRes.Commit.Count
	m.Commit.Word += mRes.Commit.Word
	m.Commit.WordFirst += mRes.Commit.WordFirst
	m.Commit.Collision += mRes.Commit.Collision

	m.Char.Count += mRes.Char.Count
	m.Char.Word += mRes.Char.Word
	m.Char.WordFirst += mRes.Char.WordFirst
	m.Char.Collision += mRes.Char.Collision

	m.Pair.Count += mRes.Pair.Count
	m.Pair.Equivalent += mRes.Pair.Equivalent
	m.Pair.SameFinger += mRes.Pair.SameFinger
	m.Pair.DoubleHit += mRes.Pair.DoubleHit
	m.Pair.TribleHit += mRes.Pair.TribleHit
	m.Pair.SingleSpan += mRes.Pair.SingleSpan
	m.Pair.MultiSpan += mRes.Pair.MultiSpan
	m.Pair.Staggered += mRes.Pair.Staggered
	m.Pair.Disturb += mRes.Pair.Disturb

	m.Pair.LeftToLeft += mRes.Pair.LeftToLeft
	m.Pair.LeftToRight += mRes.Pair.LeftToRight
	m.Pair.RightToLeft += mRes.Pair.RightToLeft
	m.Pair.RightToRight += mRes.Pair.RightToRight
}
