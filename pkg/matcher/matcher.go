package matcher

type Matcher interface {
	// 插入一个词条 word code order
	Insert(string, string, int)
	// 匹配下一个词
	Match([]rune, int) (int, string, int)
}

// 匹配算法
func New(alg string) Matcher {
	var m Matcher
	switch alg {
	case "strie", "s":
		m = NewSTrie()
	case "longest", "l":
		m = NewLongest()
	case "order", "o":
		m = NewOrder()
	case "trie", "t":
		m = NewTrie()
	default: // "trie"
		m = NewTrie()
	}
	return m
}
