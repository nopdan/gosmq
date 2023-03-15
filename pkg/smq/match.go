package smq

import (
	"strings"
	"unicode"

	"github.com/imetool/gosmq/pkg/matcher"
)

// var PUNCTS = dict.GetPuncts()

type matchRes struct {
	codes string

	wordSlice []string
	codeSlice []string
	pos       []int // 选重
}

func newMatchRes(cap int) *matchRes {
	res := new(matchRes)
	res.wordSlice = make([]string, 0, cap)
	res.codeSlice = make([]string, 0, cap)
	res.pos = make([]int, 0, cap)
	return res
}

// 只对 pos 和两个 slice
func (mr *matchRes) append(res *matchRes) {
	mr.pos = append(mr.pos, res.pos...)
	mr.wordSlice = append(mr.wordSlice, res.wordSlice...)
	mr.codeSlice = append(mr.codeSlice, res.codeSlice...)
}

func (mr *matchRes) match(text []rune, m matcher.Matcher, verbose bool, res *Result) {

	var sb strings.Builder
	sb.Grow(len(text))
	res.Basic.TextLen += len(text)
	for p := 0; p < len(text); {
		// 删掉空白字符
		switch text[p] {
		case 65533, '\n', '\r', '\t', ' ', '　':
			res.Basic.TextLen--
			p++
			continue
		}

		res.Basic.Commits++
		// 非汉字
		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			res.notHanMap[text[p]] = struct{}{}
			res.Basic.NotHanCount++
		}

		i, code, pos := m.Match(text, p)

		// 匹配到了
		if i != 0 {
			// 打词
			if i >= 2 {
				res.Words.Commits.Count++
				res.Words.Chars.Count += i
				if pos == 1 {
					res.Words.FirstCount++
				}
			}
			// 选重
			if pos >= 2 {
				res.Collision.Commits.Count++
				res.Collision.Chars.Count += i
			}
			AddTo(&res.Words.Dist, i)
			AddTo(&res.CodeLen.Dist, len(code))
			AddTo(&res.Collision.Dist, pos)

			sb.WriteString(code)
			if verbose {
				word := string(text[p : p+i])
				mr.wordSlice = append(mr.wordSlice, word)
				mr.codeSlice = append(mr.codeSlice, code)
				mr.pos = append(mr.pos, pos)
			}
			p += i
			continue
		}

		// 匹配不到

		fh := func(w, c string) {
			AddTo(&res.Words.Dist, 1) // 符号不作为打词
			AddTo(&res.CodeLen.Dist, 2)
			AddTo(&res.Collision.Dist, 1)
			sb.WriteString(c)
			if verbose {
				mr.wordSlice = append(mr.wordSlice, w)
				mr.codeSlice = append(mr.codeSlice, c)
				mr.pos = append(mr.pos, 1)
			}
			p += 2
		}
		// 单独处理这两个符号
		if p+2 < len(text) {
			switch string(text[p : p+2]) {
			case "——":
				fh("——", "=-")
				continue
			case "……":
				fh("……", "=6")
				continue
			}
		}

		// 缺汉字
		if isHan {
			res.lackMap[text[p]] = struct{}{}
			res.Basic.LackCount++
		}
		// 找不到的符号
		AddTo(&res.Words.Dist, 1)
		AddTo(&res.CodeLen.Dist, 2)
		AddTo(&res.Collision.Dist, 1)

		sb.WriteString("##")
		if verbose {
			mr.wordSlice = append(mr.wordSlice, string(text[p]))
			mr.codeSlice = append(mr.codeSlice, "##")
			mr.pos = append(mr.pos, 1)
		}
		p++
		continue
	}

	mr.codes = sb.String()
}
