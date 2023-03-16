package smq

import (
	"unicode"

	"github.com/imetool/gosmq/internal/dict"
)

type CodePosCount struct {
	Code  string
	Pos   int
	Count int
}

type matchRes struct {
	wordSlice []string
	codeSlice []string
	statData  map[string]*CodePosCount
	// pos       []int // 选重
}

func newMatchRes(cap int) *matchRes {
	res := new(matchRes)
	res.wordSlice = make([]string, 0, cap)
	res.codeSlice = make([]string, 0, cap)
	// res.pos = make([]int, 0, cap)
	res.statData = make(map[string]*CodePosCount)
	return res
}

// 只对 pos 和两个 slice
func (mr *matchRes) append(res *matchRes) {
	// mr.pos = append(mr.pos, res.pos...)
	mr.wordSlice = append(mr.wordSlice, res.wordSlice...)
	mr.codeSlice = append(mr.codeSlice, res.codeSlice...)
	for k, v := range res.statData {
		if _, ok := mr.statData[k]; !ok {
			mr.statData[k] = v
		} else {
			mr.statData[k].Count++
		}
	}
}

func (mr *matchRes) match(text []rune, dict *dict.Dict, res *Result) {

	res.Basic.TextLen += len(text)
	// 前面的键
	var last2Key, lastKey byte
	var last KeyPos
	codeHandler := func(code string) {
		for i := 0; i < len(code); i++ {
			tmpKey, tmp := res.newFeel(last2Key, lastKey, code[i], last, dict)
			last2Key = lastKey
			lastKey, last = tmpKey, tmp
		}
		AddTo(&res.CodeLen.Dist, len(code))
	}

	for p := 0; p < len(text); {
		// 跳过空白字符
		if text[p] < 33 {
			res.Basic.TextLen--
			p++
			continue
		}
		switch text[p] {
		case 65533, '　':
			res.Basic.TextLen--
			p++
			continue
		}
		res.Basic.Commits++

		i, code, pos := dict.Matcher.Match(text, p)
		// 匹配到了
		if i != 0 {
			// 对每个字都进行判断
			for j := 0; j < i; j++ {
				// 非汉字
				isHan := unicode.Is(unicode.Han, text[p+j])
				if !isHan {
					res.notHanMap[text[p+j]] = struct{}{}
					res.Basic.NotHanCount++
				}
			}

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
			AddTo(&res.Collision.Dist, pos)
			codeHandler(code)

			if dict.Split {
				word := string(text[p : p+i])
				mr.wordSlice = append(mr.wordSlice, word)
				mr.codeSlice = append(mr.codeSlice, code)
				// mr.pos = append(mr.pos, pos)
			}
			if dict.Stat {
				word := string(text[p : p+i])
				if _, ok := mr.statData[word]; !ok {
					mr.statData[word] = &CodePosCount{code, pos, 1}
				} else {
					mr.statData[word].Count++
				}
			}
			p += i
			continue
		}

		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			res.notHanMap[text[p]] = struct{}{}
			res.Basic.NotHanCount++
		}

		// 匹配不到

		fh := func(w, c string) {
			AddTo(&res.Words.Dist, 1) // 符号不作为打词
			AddTo(&res.Collision.Dist, 1)
			codeHandler(c)

			if dict.Split {
				mr.wordSlice = append(mr.wordSlice, w)
				mr.codeSlice = append(mr.codeSlice, c)
				// mr.pos = append(mr.pos, 1)
			}
			if dict.Stat {
				if _, ok := mr.statData[w]; !ok {
					mr.statData[w] = &CodePosCount{c, 1, 1}
				} else {
					mr.statData[w].Count++
				}
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
		// 找不到的符号，设为 "####"
		AddTo(&res.Words.Dist, 1)
		AddTo(&res.Collision.Dist, 1)

		code = "####"
		codeHandler(code)

		if dict.Split {
			mr.wordSlice = append(mr.wordSlice, string(text[p]))
			mr.codeSlice = append(mr.codeSlice, code)
			// mr.pos = append(mr.pos, 1)
		}
		if dict.Stat {
			word := string(text[p])
			if _, ok := mr.statData[word]; !ok {
				mr.statData[word] = &CodePosCount{code, 1, 1}
			} else {
				mr.statData[word].Count++
			}
		}
		p++
		continue
	}
}
