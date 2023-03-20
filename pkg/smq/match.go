package smq

import (
	"unicode"
	"unsafe"
)

func match(buffer []byte, dict *Dict) *matchRes {

	// 初始化
	mRes := new(matchRes)
	mRes.wordSlice = make([]string, 0, len(buffer)/3)
	mRes.codeSlice = make([]string, 0, len(buffer)/3)
	mRes.statData = make(map[string]*CodePosCount)

	mRes.notHanMap = make(map[rune]struct{})
	mRes.lackMap = make(map[rune]struct{})
	mRes.CodeLenDist = make([]int, 0)
	mRes.WordsDist = make([]int, 0)
	mRes.CollisionDist = make([]int, 0)

	// 前面的键
	var last2Key, lastKey byte
	var last KeyPos

	Handler := func(word, code string, wordLen, pos int) {
		AddTo(&mRes.WordsDist, wordLen)
		AddTo(&mRes.CollisionDist, pos)
		AddTo(&mRes.CodeLenDist, len(code))

		for i := 0; i < len(code); i++ {
			tmpKey, tmp := mRes.feel(last2Key, lastKey, code[i], last, dict)
			last2Key = lastKey
			lastKey, last = tmpKey, tmp
		}

		// 启用分词
		if dict.Split {
			mRes.wordSlice = append(mRes.wordSlice, word)
			mRes.codeSlice = append(mRes.codeSlice, code)
		}
		// 启用统计
		if dict.Stat {
			if _, ok := mRes.statData[word]; !ok {
				mRes.statData[word] = &CodePosCount{code, pos, 1}
			} else {
				mRes.statData[word].Count++
			}
		}
	}
	// 是否判断缺字
	HanHandler := func(char rune, lack bool) {
		isHan := unicode.Is(unicode.Han, char)
		// 非汉字
		if !isHan {
			mRes.notHanMap[char] = struct{}{}
			mRes.NotHanCount++
		}
		// 缺汉字
		if lack && isHan {
			mRes.lackMap[char] = struct{}{}
			mRes.LackCount++
		}
	}

	text := []rune(*(*string)(unsafe.Pointer(&buffer)))
	for p := 0; p < len(text); {
		// 跳过空白字符
		if text[p] < 33 {
			p++
			continue
		}
		switch text[p] {
		case 65533, '　':
			p++
			continue
		}
		mRes.Commits++

		wordLen, code, pos := dict.matcher.Match(text[p:])
		// 匹配到了
		if wordLen != 0 {
			sWord := string(text[p : p+wordLen])
			// 打词
			if wordLen >= 2 {
				mRes.WordsCommitsCount++
				mRes.WordsCharsCount += wordLen
				if pos == 1 {
					mRes.WordsFirstCount++
				}
			}
			// 选重
			if pos >= 2 {
				mRes.CollisionCommitsCount++
				mRes.CollisionCharsCount += wordLen
			}
			// 对每个字都进行判断
			for i := 0; i < wordLen; i++ {
				HanHandler(text[p+i], false)
			}
			Handler(sWord, code, wordLen, pos)
			p += wordLen
			continue
		}

		// 匹配不到
		if dict.Clean {
			mRes.Commits--
			p++
			continue
		}

		HanHandler(text[p], true)
		sWord := string(text[p])
		// 是否为符号
		code = PunctToCode(text[p])
		if code != "" {
			Handler(sWord, code, 1, 1)
			p++
			continue
		}
		// 单独处理这两个符号，不作为打词
		if p+1 < len(text) {
			flag := false
			switch string(text[p : p+2]) {
			case "——":
				Handler("——", "=-", 1, 1)
				flag = true
			case "……":
				Handler("……", "=6", 1, 1)
				flag = true
			}
			if flag {
				HanHandler(text[p+1], false)
				p += 2
				continue
			}
		}
		// 找不到的符号，设为 "####"
		Handler(sWord, "####", 1, 1)
		p++
	}
	return mRes
}
