package matcher

type codePos struct {
	code string
	pos  int
}

type single map[rune]*codePos

func NewSingle() *single {
	t := make(single, 1024)
	return &t
}

func (s *single) Insert(word, code string, pos int) {
	char := []rune(word)[0]
	// 同一个字取码长较短的
	if cp, ok := (*s)[char]; !ok || len(cp.code) > len(code) {
		(*s)[char] = &codePos{code, pos}
	}
}

func (s *single) Build() {
}

func (s *single) Match(text []rune) (int, string, int) {
	if v, ok := (*s)[text[0]]; ok {
		return 1, v.code, v.pos
	}
	return 0, "", 1
}
