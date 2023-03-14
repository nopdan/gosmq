package smq

// count and rate
type CountRate struct {
	Count int
	Rate  float64
}

type CodePosCount struct {
	Code  string
	Pos   int
	Count int
}

type Result struct {
	Name      string
	Basic     basic
	Words     words     // 打词
	Collision collision // 选重
	CodeLen   codeLen   // 码长

	Keys    keys  // 按键统计
	Combs   combs // 按键组合
	Fingers fingers
	Hands   hands

	toTalEq10 int // 总当量*10
	mapKeys   map[byte]int
	notHanMap map[rune]struct{}
	lackMap   map[rune]struct{}
	// codes     string
}

// 基础
type basic struct {
	DictLen     int    // 词条数
	TextLen     int    // 文本字数
	NotHan      string // 非汉字
	NotHans     int    // 非汉字数（去重）
	NotHanCount int    // 非汉字计数
	Lack        string // 缺字
	Lacks       int    // 缺字数（去重）
	LackCount   int    // 缺字计数
	Commits     int    // 上屏数
}

// 打词
type words struct {
	Commits CountRate // 打词数
	Chars   CountRate // 打词字数
	Dist    []int     // 词长分布统计

	FirstCount int // 首选词
}

// 选重
type collision struct {
	Commits CountRate // 选重数
	Chars   CountRate // 选重字数
	Dist    []int     // 选重分布统计
}

// 码长
type codeLen struct {
	Total   int     // 全部码长
	PerChar float64 // 字均码长
	Dist    []int   // 码长分布统计
}

// 按键 左空格_，右空格+
type keys map[string]*CountRate

// 按键组合
type combs struct {
	Count      int     // 按键组合数
	Equivalent float64 // 当量

	DoubleHit  CountRate // 同键双击
	TribleHit  CountRate // 同键三连击
	SingleSpan CountRate // 小跨排
	MultiSpan  CountRate // 大跨排

	LongFingersDisturb   CountRate // 错手
	LittleFingersDisturb CountRate // 小指干扰
}

type fingers struct {
	Dist [11]*CountRate // 手指分布，按键盘上的列，第11个是41键以外的
	Same CountRate      // 同指
	Diff CountRate      // 异指（同手）
}

type hands struct {
	Left  CountRate // 左手
	Right CountRate // 右手
	Same  CountRate // 同手
	Diff  CountRate // 异手

	LL CountRate `json:"LeftToLeft"`   // 左左
	LR CountRate `json:"LeftToRight"`  // 左右
	RL CountRate `json:"RightToLeft"`  // 右左
	RR CountRate `json:"RightToRight"` // 右右
}

func newResult() *Result {
	res := new(Result)
	res.mapKeys = make(map[byte]int, 128)
	res.notHanMap = make(map[rune]struct{}, 100)
	res.lackMap = make(map[rune]struct{}, 10)
	res.Words.Dist = make([]int, 1)
	res.Collision.Dist = make([]int, 1)
	res.CodeLen.Dist = make([]int, 1)
	res.Keys = make(keys)
	const ALL_KEYS = "1234567890qwertyuiopasdfghjkl;'zxcvbnm,./_+"
	for i := 0; i < len(ALL_KEYS); i++ {
		res.Keys[string(ALL_KEYS[i])] = new(CountRate)
		res.mapKeys[ALL_KEYS[i]] = 0
	}
	for i := 0; i < 11; i++ {
		res.Fingers.Dist[i] = new(CountRate)
	}
	return res
}
