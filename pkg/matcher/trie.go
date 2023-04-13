package matcher

type trie struct {
	tn    *node
	code  []string
	pos   []byte
	count uint32
}

// trie 树
type node struct {
	ch  map[rune]*node
	idx uint32
}

func NewTrie() *trie {
	t := new(trie)
	t.tn = new(node)
	t.tn.ch = make(map[rune]*node, 1000)
	t.code = make([]string, 0, 10000)
	t.pos = make([]byte, 0, 10000)

	t.code = append(t.code, "")
	t.pos = append(t.pos, 0)
	return t
}

func (t *trie) Insert(word, code string, pos int) {
	tn := t.tn
	for _, v := range word {
		if tn.ch == nil {
			tn.ch = make(map[rune]*node)
			tn.ch[v] = new(node)
		} else if _, ok := tn.ch[v]; !ok {
			tn.ch[v] = new(node)
		}
		tn = tn.ch[v]
	}
	// 同一个词取码长较短的
	if tn.idx == 0 {
		t.code = append(t.code, code)
		t.pos = append(t.pos, byte(pos))
		t.count++
		tn.idx = t.count
	} else if len(t.code[tn.idx]) > len(code) {
		t.code[tn.idx] = code
		t.pos[tn.idx] = byte(pos)
	}
}

func (t *trie) Build() {
}

// 前缀树最长匹配
func (t *trie) Match(text []rune) (int, string, int) {
	var wordLen int
	var code string
	var pos byte

	tn := t.tn
	for p := 0; p < len(text); {
		tn = tn.ch[text[p]]
		p++
		if tn == nil {
			break
		}
		if tn.idx != 0 {
			wordLen = p
			code = t.code[tn.idx]
			pos = t.pos[tn.idx]
		}
	}
	return wordLen, code, int(pos)
}
