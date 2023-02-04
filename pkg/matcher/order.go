package matcher

import "github.com/imetool/dtool/pkg/table"

// 顺序匹配
type order table.Table

func NewOrder() *order {
	o := make(order, 0, 9999)
	return &o
}

func (o *order) Insert(word, code string, pos int) {
	*o = append(*o, table.Entry{word, code, pos})
}

// 顺序匹配
func (o order) Match(text []rune, p int) (int, string, int) {
	for _, v := range o {
		word := []rune(v.Word)
		if p+len(word) > len(text) {
			continue
		}
		if v.Word == string(text[p:p+len(word)]) {
			return len(word), v.Code, v.Pos
		}
	}
	return 0, "", 1
}
