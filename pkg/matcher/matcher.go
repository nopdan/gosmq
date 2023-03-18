package matcher

type Matcher interface {
	// 插入一个词条 word code pos
	Insert(string, string, int)
	// 构建
	Build()
	// 匹配下一个词，返回匹配到的词长，编码和候选位置
	Match([]rune) (int, string, int)
}

// 匹配算法
func New(alg string) Matcher {
	var m Matcher
	switch alg {
	case "single":
		m = NewSingle()
		// fmt.Println("匹配算法：单字专用 hashMap(with rune key)")
	case "strie", "s":
		m = NewStableTrie()
		// fmt.Println("匹配算法：稳定的 trie(hashMap impl)")
	default: // 默认 trie 算法
		m = NewTrie()
		// fmt.Println("匹配算法：trie(hashMap impl)")
	}
	return m
}
