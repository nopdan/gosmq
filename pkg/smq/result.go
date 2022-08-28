package smq

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

	Data export

	toTalEq10 int // 总当量*10
	mapKeys   map[byte]int
	mapNotHan map[rune]struct{}
	mapLack   map[rune]struct{}
	wordsDist []int
	collDist  []int
	codeDist  []int
	// codes     string
}

// count and rate
type CaR struct {
	Count int
	Rate  float64
}

type CoC struct {
	Code  string
	Order int
	Count int
}

// 可能要导出的数据
type export struct {
	WordSlice []string // 分词
	CodeSlice []string // 编码
	Details   map[string]*CoC
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
	Commits CaR   // 打词数
	Chars   CaR   // 打词字数
	Dist    []int // 词长分布统计

	FirstCount int // 首选词
}

// 选重
type collision struct {
	Commits CaR   // 选重数
	Chars   CaR   // 选重字数
	Dist    []int // 选重分布统计
}

// 码长
type codeLen struct {
	Total   int     // 全部码长
	PerChar float64 // 字均码长
	Dist    []int   // 码长分布统计
}

// 按键 左空格_，右空格+
type keys map[string]*CaR

// 按键组合
type combs struct {
	Count      int     // 按键组合数
	Equivalent float64 // 当量

	DoubleHit  CaR // 同键双击
	TribleHit  CaR // 同键三连击
	SingleSpan CaR // 小跨排
	MultiSpan  CaR // 大跨排

	LongFingersDisturb   CaR // 错手
	LittleFingersDisturb CaR // 小指干扰
}

type fingers struct {
	Dist [11]*CaR // 手指分布，按键盘上的列，第11个是41键以外的
	Same CaR      // 同指
	Diff CaR      // 异指（同手）
}

type hands struct {
	Left  CaR // 左手
	Right CaR // 右手
	Same  CaR // 同手
	Diff  CaR // 异手

	LL CaR `json:"LeftToLeft"`   // 左左
	LR CaR `json:"LeftToRight"`  // 左右
	RL CaR `json:"RightToLeft"`  // 右左
	RR CaR `json:"RightToRight"` // 右右
}

func newResult() *Result {
	res := new(Result)
	res.mapKeys = make(map[byte]int)
	res.mapLack = make(map[rune]struct{})
	res.mapNotHan = make(map[rune]struct{})
	res.wordsDist = make([]int, 1)
	res.collDist = make([]int, 1)
	res.codeDist = make([]int, 1)
	res.Data.WordSlice = make([]string, 0, 1024)
	res.Data.CodeSlice = make([]string, 0, 1024)
	res.Data.Details = make(map[string]*CoC)
	res.Keys = make(keys)
	const ALL_KEYS = "1234567890qwertyuiopasdfghjkl;'zxcvbnm,./"
	for i := 0; i < len(ALL_KEYS); i++ {
		res.Keys[string(ALL_KEYS[i])] = new(CaR)
		res.mapKeys[ALL_KEYS[i]] = 0
	}
	for i := 0; i < 11; i++ {
		res.Fingers.Dist[i] = new(CaR)
	}
	return res
}
