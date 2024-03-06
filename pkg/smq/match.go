package smq

import (
	"unicode"
	"unsafe"

	"github.com/nopdan/gosmq/pkg/dict"
	"github.com/nopdan/gosmq/pkg/feeling"
	"github.com/nopdan/gosmq/pkg/result"
	"github.com/nopdan/gosmq/pkg/util"
)

func (c *Config) match(buffer []byte, dict *dict.Dict) *result.MatchRes {

	// 初始化
	mRes := result.NewMatchRes()
	feel := feeling.New(mRes, dict.SpacePref)

	Handler := func(word, code string, wordLen, pos int) {
		util.Increase(&mRes.WordLenDist, wordLen)
		util.Increase(&mRes.CollisionDist, pos)
		util.Increase(&mRes.CodeLenDist, len(code))

		for i := 0; i < len(code); i++ {
			feel.Process(code[i])
		}

		// 启用分词
		if c.Split {
			mRes.Segment = append(mRes.Segment, result.WordCode{
				Word: word,
				Code: code,
			})
		}
		// 启用统计
		if c.Stat {
			if _, ok := mRes.StatData[word]; !ok {
				mRes.StatData[word] = &result.CodePosCount{
					Code:  code,
					Pos:   pos,
					Count: 1}
			} else {
				mRes.StatData[word].Count++
			}
		}
	}
	// 是否判断缺字
	HanHandler := func(char rune, lack bool) {
		isHan := unicode.Is(unicode.Han, char)
		// 非汉字
		if !isHan {
			mRes.NotHanMap[char] = struct{}{}
		}
		// 缺汉字
		if lack && isHan {
			mRes.LackMap[char] = struct{}{}
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
		mRes.Commit.Count++

		wordLen, code, pos := dict.Matcher.Match(text[p:])
		// 匹配到了
		if wordLen != 0 {
			sWord := string(text[p : p+wordLen])
			// 打词
			if wordLen >= 2 {
				mRes.Commit.Word++
				mRes.Commit.WordChars += wordLen
				if pos == 1 {
					mRes.Commit.WordFirst++
				}
			}
			// 选重
			if pos >= 2 {
				mRes.Commit.Collision++
				mRes.Commit.CollisionChars += wordLen
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
		if c.Clean {
			mRes.Commit.Count--
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
		// 找不到的符号，设为 "######"
		Handler(sWord, "######", 1, 1)
		p++
	}
	return mRes
}
