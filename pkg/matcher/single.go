package matcher

import (
	"bytes"
	"unicode/utf8"
)

type Single struct {
	dict map[rune]*struct {
		code string
		pos  int
	}
}

func NewSingle() *Single {
	s := new(Single)
	s.dict = make(map[rune]*struct {
		code string
		pos  int
	}, 1024)
	return s
}

func (s *Single) Insert(word, code string, pos int) {
	char, _ := utf8.DecodeRuneInString(word)
	cp, ok := s.dict[char]
	if ok {
		if len(cp.code) < len(code) {
			// 同一个字取码长较短的
			s.dict[char].pos = pos
		}
		return
	}
	s.dict[char] = &struct {
		code string
		pos  int
	}{
		code: code,
		pos:  pos,
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
	if v, ok := s.dict[ch]; ok {
		res.Code = v.code
		res.Pos = v.pos
	}
}
