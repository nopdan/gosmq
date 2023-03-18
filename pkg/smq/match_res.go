package smq

type CodePosCount struct {
	Code  string
	Pos   int
	Count int
}

// 匹配一段文字得到的信息
type matchRes struct {
	dictIdx int // 码表索引

	wordSlice []string
	codeSlice []string
	statData  map[string]*CodePosCount

	keysDist  [128]int
	notHanMap map[rune]struct{}
	lackMap   map[rune]struct{}

	Commits int

	NotHanCount int // 非汉字计数
	LackCount   int

	WordsCommitsCount int
	WordsCharsCount   int
	WordsFirstCount   int

	CollisionCommitsCount int
	CollisionCharsCount   int

	CodeLenDist   []int
	WordsDist     []int
	CollisionDist []int

	toTalEq10  int
	CombsCount int

	SameFingers int
	Hands       struct {
		LL int
		LR int
		RL int
		RR int
	}
	Combs struct {
		DoubleHit            int
		TribleHit            int
		SingleSpan           int
		MultiSpan            int
		LongFingersDisturb   int
		LittleFingersDisturb int
	}
}

// 将每次匹配得到的信息追加到总结果
func (res *Result) append(mRes *matchRes, dict *Dict) {
	if dict.Verbose {
		for k, v := range mRes.statData {
			if _, ok := res.statData[k]; !ok {
				res.statData[k] = v
			} else {
				res.statData[k].Count += v.Count
			}
		}
	}
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

	for i := 33; i < 128; i++ {
		res.keysDist[i] += mRes.keysDist[i]
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
