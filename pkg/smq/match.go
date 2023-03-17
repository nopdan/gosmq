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

// 匹配一段文字得到的信息
type matchRes struct {
	dictIdx int // 码表索引
	textIdx int // 文本段索引

	wordSlice []string
	codeSlice []string
	statData  map[string]*CodePosCount

	mapKeys   map[byte]int
	notHanMap map[rune]struct{}
	lackMap   map[rune]struct{}

	TextLen int
	Commits int

	NotHanCount int // 非汉字计数
	LackCount   int

	WordsCommitsCount int
	WordsCharsCount   int
	WordsFirstCount   int

	CollisionCommitsCount int
	CollisionCharsCount   int

	CodeLenDist   []int
	WordsDist     []int
	CollisionDist []int

	toTalEq10  int
	CombsCount int

	SameFingers int
	Hands       struct {
		LL int
		LR int
		RL int
		RR int
	}
	Combs struct {
		DoubleHit            int
		TribleHit            int
		SingleSpan           int
		MultiSpan            int
		LongFingersDisturb   int
		LittleFingersDisturb int
	}
}

func match(text []rune, dict *dict.Dict) *matchRes {

	// 初始化
	mRes := new(matchRes)
	mRes.wordSlice = make([]string, 0, len(text)/3)
	mRes.codeSlice = make([]string, 0, len(text)/3)
	mRes.statData = make(map[string]*CodePosCount)

	mRes.mapKeys = make(map[byte]int)
	mRes.notHanMap = make(map[rune]struct{})
	mRes.lackMap = make(map[rune]struct{})
	mRes.CodeLenDist = make([]int, 0)
	mRes.WordsDist = make([]int, 0)
	mRes.CollisionDist = make([]int, 0)

	mRes.TextLen = len(text)
	// 前面的键
	var last2Key, lastKey byte
	var last KeyPos
	codeHandler := func(code string) {
		for i := 0; i < len(code); i++ {
			tmpKey, tmp := mRes.newFeel(last2Key, lastKey, code[i], last, dict)
			last2Key = lastKey
			lastKey, last = tmpKey, tmp
		}
		AddTo(&mRes.CodeLenDist, len(code))
	}

	for p := 0; p < len(text); {
		// 跳过空白字符
		if text[p] < 33 {
			mRes.TextLen--
			p++
			continue
		}
		switch text[p] {
		case 65533, '　':
			mRes.TextLen--
			p++
			continue
		}
		mRes.Commits++

		i, code, pos := dict.Matcher.Match(text, p)
		// 匹配到了
		if i != 0 {
			// 对每个字都进行判断
			for j := 0; j < i; j++ {
				// 非汉字
				isHan := unicode.Is(unicode.Han, text[p+j])
				if !isHan {
					mRes.notHanMap[text[p+j]] = struct{}{}
					mRes.NotHanCount++
				}
			}

			// 打词
			if i >= 2 {
				mRes.WordsCommitsCount++
				mRes.WordsCharsCount += i
				if pos == 1 {
					mRes.WordsFirstCount++
				}
			}
			// 选重
			if pos >= 2 {
				mRes.CollisionCommitsCount++
				mRes.CollisionCharsCount += i
			}
			AddTo(&mRes.WordsDist, i)
			AddTo(&mRes.CollisionDist, pos)
			codeHandler(code)

			if dict.Split {
				word := string(text[p : p+i])
				mRes.wordSlice = append(mRes.wordSlice, word)
				mRes.codeSlice = append(mRes.codeSlice, code)
			}
			if dict.Stat {
				word := string(text[p : p+i])
				if _, ok := mRes.statData[word]; !ok {
					mRes.statData[word] = &CodePosCount{code, pos, 1}
				} else {
					mRes.statData[word].Count++
				}
			}
			p += i
			continue
		}

		isHan := unicode.Is(unicode.Han, text[p])
		if !isHan {
			mRes.notHanMap[text[p]] = struct{}{}
			mRes.NotHanCount++
		}

		// 匹配不到

		fh := func(w, c string) {
			AddTo(&mRes.WordsDist, 1) // 符号不作为打词
			AddTo(&mRes.CollisionDist, 1)
			codeHandler(c)

			if dict.Split {
				mRes.wordSlice = append(mRes.wordSlice, w)
				mRes.codeSlice = append(mRes.codeSlice, c)
			}
			if dict.Stat {
				if _, ok := mRes.statData[w]; !ok {
					mRes.statData[w] = &CodePosCount{c, 1, 1}
				} else {
					mRes.statData[w].Count++
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
			mRes.lackMap[text[p]] = struct{}{}
			mRes.LackCount++
		}
		// 找不到的符号，设为 "####"
		AddTo(&mRes.WordsDist, 1)
		AddTo(&mRes.CollisionDist, 1)

		code = "####"
		codeHandler(code)

		if dict.Split {
			mRes.wordSlice = append(mRes.wordSlice, string(text[p]))
			mRes.codeSlice = append(mRes.codeSlice, code)
			// mRes.pos = append(mRes.pos, 1)
		}
		if dict.Stat {
			word := string(text[p])
			if _, ok := mRes.statData[word]; !ok {
				mRes.statData[word] = &CodePosCount{code, 1, 1}
			} else {
				mRes.statData[word].Count++
			}
		}
		p++
		continue
	}
	return mRes
}
