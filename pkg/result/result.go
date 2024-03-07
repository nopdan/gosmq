package result

// count and rate
type CountRate struct {
	Count int
	Rate  float64
}

type Result struct {
	segments []struct {
		PartIdx int
		Segment []WordCode
	}
	statData map[string]*CodePosCount

	Info   Info   // 文章和码表信息
	Commit commit // 上屏
	Pair   pair   // 按键组合
	Keys   keys   // 按键统计
	Han    han    // 非汉字以及缺字
	// 各种分布
	Dist struct {
		CodeLen   []int   // 码长
		WordLen   []int   // 词长
		Collision []int   // 选重
		Finger    [11]int // 手指
	}
	// 码长
	CodeLen struct {
		Total   int
		PerChar float64
	}
	LeftHand  int // 左手按键数
	RightHand int // 右手按键数
}

type Info struct {
	TextName string
	TextLen  int // 文本字数
	DictName string
	DictLen  int  // 词条数
	Single   bool // 是否为单字码表
}

type han struct {
	NotHan      string // 非汉字
	NotHans     int    // 非汉字数（去重）
	NotHanCount int    // 非汉字计数
	Lack        string // 缺字
	Lacks       int    // 缺字数（去重）
	LackCount   int    // 缺字计数
}

// 码长
type codeLen struct {
	Total   int     // 全部码长
	PerChar float64 // 字均码长
	Dist    []int   // 码长分布统计
}

// 按键 左空格_，右空格+
type keys map[string]*CountRate

type hands struct {
	Left  CountRate // 左手
	Right CountRate // 右手
}
