package smq

import (
	"strings"
	"unicode"

	"github.com/imetool/gosmq/internal/dict"
)

var PUNCTS = dict.GetPuncts()

func (res *Result) match(text []rune, dict *dict.Dict) string {
	var sb strings.Builder
	sb.Grow(len(text))
	res.Basic.TextLen += len(text)
	for p := 0; p < len(text); {
		// 删掉空白字符
		switch text[p] {
		case 65533, '\n', '\r', '\t', ' ', '　':
			p++
			res.Basic.TextLen--
			continue
		}
		// 非汉字
		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			res.Basic.NotHanCount++
			res.mapNotHan[text[p]] = struct{}{}
		}

		i, code, pos := dict.Matcher.Match(text, p)
		// 缺字
		if i == 0 {
			if isHan {
				res.Basic.LackCount++
				res.mapLack[text[p]] = struct{}{}
			}
			sb.WriteByte(' ')
			p++
			continue
		}

		sb.WriteString(code)
		if i == 2 && PUNCTS[string(text[p:p+2])] != "" {
			AddTo(&res.wordsDist, 1)
		} else {
			AddTo(&res.wordsDist, i) // 词长分布
		}

		if pos != 1 {
			res.Collision.Chars.Count += i // 选重字数
		} else if i != 1 {
			res.Words.FirstCount++ // 首选词
		}
		AddTo(&res.collDist, pos)       // 选重分布
		AddTo(&res.codeDist, len(code)) // 码长分布

		if dict.Verbose {
			word := string(text[p : p+i])
			res.Data.WordSlice = append(res.Data.WordSlice, word)
			res.Data.CodeSlice = append(res.Data.CodeSlice, code)
			if _, ok := res.Data.Details[word]; !ok {
				res.Data.Details[word] = &CodePosCount{Code: code, Pos: i}
			}
			res.Data.Details[word].Count++
		}
		p += i
	}
	return sb.String()
}
