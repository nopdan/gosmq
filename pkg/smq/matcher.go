package smq

import "github.com/cxcn/gosmq/pkg/matcher"

type Matcher interface {
	// 插入一个词条 word code order
	Insert(string, string, int)
	// 读取完码表后的操作
	Handle()
	// 匹配下一个词
	Match([]rune, int) (int, string, int)
}

// 匹配算法
func (dict *Dict) match() {
	if dict.Matcher == nil {
		switch dict.Algorithm {
		case "strie","s":
			dict.Matcher = matcher.NewSTrie()
		case "longest","l":
			dict.Matcher = matcher.NewLongest()
		case "order","o":
			dict.Matcher = matcher.NewOrder()
		case "trie","t":
			dict.Matcher = matcher.NewTrie()
		default: // "trie"
			dict.Matcher = matcher.NewTrie()
		}
	}
}
