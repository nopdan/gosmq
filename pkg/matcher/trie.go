package matcher

// trie 树
type trie struct {
	ch   map[rune]*trie
	code string
	pos  int
}

func NewTrie() *trie {
	t := new(trie)
	t.ch = make(map[rune]*trie, 1000)
	return t
}

func (t *trie) Insert(word, code string, pos int) {
	for _, v := range word {
		if t.ch == nil {
			t.ch = make(map[rune]*trie)
			t.ch[v] = new(trie)
		} else if t.ch[v] == nil {
			t.ch[v] = new(trie)
		}
		t = t.ch[v]
	}
	// 同一个词取码长较短的
	if t.code == "" || len(t.code) > len(code) {
		t.code = code
		t.pos = pos
	}
}

func (t *trie) Build() {
}

// 前缀树最长匹配
func (t *trie) Match(text []rune) (int, string, int) {
	var wordLen int
	var code string
	var pos int

	for p := 0; p < len(text); {
		t = t.ch[text[p]]
		p++
		if t == nil {
			break
		}
		if t.code != "" {
			wordLen = p
			code = t.code
			pos = t.pos
		}
	}
	return wordLen, code, pos
}
