package matcher

type single map[rune]codePos

func NewSingle() *single {
	t := make(single, 2e4)
	return &t
}

func (s *single) Insert(word, code string, pos int) {
	char := []rune(word)[0]
	// 不替换原有的
	if co, ok := (*s)[char]; !ok || co.code == "" || len(code) < len(co.code) {
		(*s)[char] = codePos{code, pos}
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
