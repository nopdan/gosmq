package matcher

import (
	"fmt"

	"github.com/imetool/dtool/pkg/table"
)

type single map[rune]codePos

func NewSingle() *single {
	t := make(single, 2e4)
	return &t
}

func (s *single) Insert(e table.Entry) {
	char := []rune(e.Word)[0]
	// 不替换原有的
	if co, ok := (*s)[char]; !ok || co.code == "" || len(e.Code) < len(co.code) {
		(*s)[char] = codePos{e.Code, e.Pos}
	}
}

func (s *single) Build(t table.Table) {
	fmt.Println("匹配算法：rune hashMap")
	for i := range t {
		s.Insert(t[i])
	}
}

// 最长匹配
func (s *single) Match(text []rune, p int) (int, string, int) {
	v, ok := (*s)[text[p]]
	if !ok {
		return 0, "", 1
	}
	return 1, v.code, v.pos
}
