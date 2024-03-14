package matcher

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

type codePos struct {
	code string
	pos  int
}

type Single struct {
	dict map[rune]*codePos
}

func NewSingle() *Single {
	s := new(Single)
	s.dict = make(map[rune]*codePos, 1024)
	return s
}

func (s *Single) Insert(word, code string, pos int) {
	char, _ := utf8.DecodeRuneInString(word)
	if _, ok := s.dict[char]; !ok {
		s.dict[char] = &codePos{
			code: code,
			pos:  pos,
		}
	}
}

func (s *Single) Build() {
}

func (s *Single) Match(brd *bytes.Reader, res *Result) {
	res.Reset()
	ch, size, _ := brd.ReadRune()
	res.Char = ch
	res.Size = size
	res.Length = 1
	if unicode.IsSpace(ch) {
		return
	}
	if v, ok := s.dict[ch]; ok {
		res.Code = v.code
		res.Pos = v.pos
	}
}
