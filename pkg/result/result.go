package result

type Result struct {
	segments []Segment
	statData map[string]*CodePosCount

	Info   Info   // 文章和码表信息
	Commit commit // 上屏
	Char   char   // 上屏字数
	Han    han    // 非汉字以及缺字
	Pair   pair   // 按键组合
	Keys   keys   // 按键统计
	// 各种分布
	Dist struct {
		CodeLen   []int   // 码长
		WordLen   []int   // 词长
		Collision []int   // 选重
		Finger    [11]int // 手指

		// 按键 左空格_，右空格+
		Key map[string]int
	}
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

type keys struct {
	Count     int     // 按键数
	CodeLen   float64 // 字均码长
	LeftHand  int     // 左手按键数
	RightHand int     // 右手按键数
}
