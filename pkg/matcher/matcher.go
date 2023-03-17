package matcher

import "github.com/imetool/dtool/pkg/table"

type Matcher interface {
	// 插入一个词条 word code pos
	// Insert(string, string, int)
	Build(table.Table)
	// 匹配下一个词
	Match([]rune, int) (int, string, int)
}

// 匹配算法
func New(alg string) Matcher {
	var m Matcher
	switch alg {
	case "single":
		// fmt.Println("匹配算法：单字专用 hashMap(with rune key)")
		m = NewSingle()
	case "longest", "l":
		// fmt.Println("匹配算法：最长匹配")
		m = NewLongest()
	case "strie", "s":
		// fmt.Println("匹配算法：稳定的 trie(hashMap impl)")
		m = NewSTrie()
	case "trie", "t":
		// fmt.Println("匹配算法：trie(hashMap impl)")
		m = NewTrie()
	default: // 默认 trie 算法
		m = NewTrie()
	}
	return m
}
