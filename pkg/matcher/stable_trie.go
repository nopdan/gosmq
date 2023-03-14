package matcher

import (
	"sort"

	"github.com/imetool/dtool/pkg/table"
)

// 稳定 trie 树
type sTrie struct {
	children map[rune]*sTrie
	code     string
	pos      int
	line     uint32
}

func NewSTrie() *sTrie {
	return new(sTrie)
}

var orderLine uint32 = 0

func (t *sTrie) Insert(e table.Entry) {
	for _, v := range e.Word {
		if t.children == nil {
			t.children = make(map[rune]*sTrie)
			t.children[v] = new(sTrie)
		} else if t.children[v] == nil {
			t.children[v] = new(sTrie)
		}
		t = t.children[v]
	}
	if t.code == "" {
		t.code = e.Code
		t.pos = e.Pos
		orderLine++
		t.line = orderLine
	}
}

func (st *sTrie) InsertAll(t table.Table) {
	for i := range t {
		st.Insert(t[i])
	}
}

// 前缀树按码表序匹配
func (t *sTrie) Match(text []rune, p int) (int, string, int) {
	j := 0 // 已匹配的字数
	i := 0 // 有编码的匹配
	dict := t
	type res_tmp struct {
		i    int
		code string
		pos  int
		line uint32
	}
	res := make([]res_tmp, 0, 10)
	for p+j < len(text) {
		dict = dict.children[text[p+j]]
		j++
		if dict == nil {
			break
		}
		if dict.code != "" {
			i = j
			res = append(res, res_tmp{i, dict.code, dict.pos, dict.line})
		}
	}
	if len(res) == 0 {
		return 0, "", 1
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].line < res[j].line
	})
	return res[0].i, res[0].code, res[0].pos
}
