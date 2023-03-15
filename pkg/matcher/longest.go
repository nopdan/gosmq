package matcher

import (
	"fmt"

	"github.com/imetool/dtool/pkg/table"
)

type codePos struct {
	code string
	pos  int
}

// 最长匹配
type longest []map[string]codePos

func NewLongest() *longest {
	l := make(longest, 0, 20)
	return &l
}

func (l *longest) Insert(e table.Entry) {
	i := len([]rune(e.Word)) // 词长
	for i+1 > len(*l) {      // 扩容
		*l = append(*l, make(map[string]codePos))
	}
	// 不替换原有的
	if co, ok := (*l)[i][e.Word]; !ok || co.code == "" || len(e.Code) < len(co.code) {
		(*l)[i][e.Word] = codePos{e.Code, e.Pos}
	}
}

func (l *longest) Build(t table.Table) {
	fmt.Println("匹配算法：最长匹配")
	for i := range t {
		l.Insert(t[i])
	}
}

// 最长匹配
func (l longest) Match(text []rune, p int) (int, string, int) {
	max := len(l) - 1
	if len(text)-p < max {
		max = len(text) - p
	}
	for i := max; i > 0; i-- {
		if v, ok := l[i][string(text[p:p+i])]; ok {
			return i, v.code, v.pos
		}
	}
	return 0, "", 1
}
