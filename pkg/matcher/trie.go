package matcher

// trie 树
type trie struct {
	children map[rune]*trie
	code     string
	pos      int
}

func NewTrie() *trie {
	return new(trie)
}

func (t *trie) Insert(word, code string, pos int) {
	for _, v := range word {
		if t.children == nil {
			t.children = make(map[rune]*trie)
			t.children[v] = new(trie)
		} else if t.children[v] == nil {
			t.children[v] = new(trie)
		}
		t = t.children[v]
	}
	if t.code == "" || len(code) < len(t.code) {
		t.code = code
		t.pos = pos
	}
}

// 前缀树最长匹配
func (t *trie) Match(text []rune, p int) (int, string, int) {
	j := 0 // 已匹配的字数
	i := 0 // 有编码的匹配
	dict := t
	code := "" // 编码
	pos := 0
	for p+j < len(text) {
		dict = dict.children[text[p+j]]
		j++
		if dict == nil {
			break
		}
		if dict.code != "" {
			i = j
			code = dict.code
			pos = dict.pos
		}
	}
	return i, code, pos
}
