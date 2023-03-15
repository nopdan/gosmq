package matcher

import (
	"fmt"

	"github.com/imetool/dtool/pkg/table"
)

// trie 树
type trie struct {
	ch   map[rune]*trie
	code string
	pos  int
}

func NewTrie() *trie {
	return new(trie)
}

func (t *trie) Insert(e table.Entry) {
	for _, v := range e.Word {
		if t.ch == nil {
			t.ch = make(map[rune]*trie)
			t.ch[v] = new(trie)
		} else if t.ch[v] == nil {
			t.ch[v] = new(trie)
		}
		t = t.ch[v]
	}
	// 若已存在，保留原来的
	if t.code == "" {
		t.code = e.Code
		t.pos = e.Pos
	}
}

func (t *trie) Build(tb table.Table) {
	fmt.Println("匹配算法：trie(hashMap impl)")
	for i := range tb {
		t.Insert(tb[i])
	}
}

// 前缀树最长匹配
func (t *trie) Match(text []rune, p int) (int, string, int) {
	j := 0     // 已匹配的字数
	i := 0     // 有编码的匹配
	code := "" // 编码
	pos := 0
	for p+j < len(text) {
		t = t.ch[text[p+j]]
		j++
		if t == nil {
			break
		}
		if t.code != "" {
			i = j
			code = t.code
			pos = t.pos
		}
	}
	return i, code, pos
}
