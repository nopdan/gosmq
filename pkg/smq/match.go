package smq

import (
	"bytes"
	"io"
	"unicode"

	"github.com/nopdan/gosmq/pkg/dict"
	"github.com/nopdan/gosmq/pkg/matcher"
	"github.com/nopdan/gosmq/pkg/result"
	"github.com/nopdan/gosmq/pkg/util"
)

func (c *Config) match(buffer []byte, dict *dict.Dict) *result.MatchRes {
	mRes := result.NewMatchRes()
	feel := NewFeeling(mRes, dict.SpacePref)
	brd := bytes.NewReader(buffer)
	res := new(matcher.Result)

	process := func(res *matcher.Result) {
		mRes.Commit.Count++
		util.Increase(&mRes.Dist.WordLen, res.Length)
		util.Increase(&mRes.Dist.Collision, res.Pos)
		util.Increase(&mRes.Dist.CodeLen, len(res.Code))
		for i := range len(res.Code) {
			feel.Process(res.Code[i])
		}
		if res.Pos >= 2 {
			mRes.Commit.Collision++
			mRes.Commit.CollisionChars += res.Length
		}
		// 匹配到词组
		if res.Length >= 2 {
			mRes.Commit.Word++
			mRes.Commit.WordChars += res.Length
			if res.Pos == 1 {
				mRes.Commit.WordFirst++ // 首选词
			}
		}
		if !c.Split && !c.Stat {
			return
		}

		var word string
		if res.Char > 32 {
			word = string([]rune{res.Char})
		} else {
			brd.Seek(-1*int64(res.Size), io.SeekCurrent)
			data := make([]byte, res.Size)
			brd.Read(data)
			word = util.UnsafeToString(data)
		}
		// 启用分词
		if c.Split {
			mRes.Segment = append(mRes.Segment, result.WordCode{
				Word: word,
				Code: res.Code,
			})
		}
		// 启用统计
		if c.Stat {
			if _, ok := mRes.StatData[word]; !ok {
				mRes.StatData[word] = &result.CodePosCount{
					Code:  res.Code,
					Pos:   res.Pos,
					Count: 1}
			} else {
				mRes.StatData[word].Count++
			}
		}
	}

	for brd.Len() > 0 {
		// 跳过空白字符
		ch, _, _ := brd.ReadRune()
		if ch < 33 || ch == 65533 || ch == '　' {
			continue
		}
		_ = brd.UnreadRune()

		// 开始匹配
		dict.Matcher.Match(brd, res)

		// 匹配成功
		if res.Pos > 0 {
			process(res)
			continue
		}

		// 匹配失败了
		if c.Clean {
			continue
		}
		res.Pos = 1

		// 两个字符的符号
		if res.Char == '—' || res.Char == '…' {
			ch2, _, err := brd.ReadRune()
			if err != nil {
				if res.Char == '—' && ch2 == '—' {
					// 中文破折号 —— 占用 6 字节，不计打词
					res.SetChar(0).SetCode("=-").SetSize(6)
					process(res)
					continue
				} else if res.Char == '…' && ch2 == '…' {
					// 中文省略号 …… 占用 6 字节，不计打词
					res.SetChar(0).SetCode("=6").SetSize(6)
					process(res)
					continue
				}
			}
			_ = brd.UnreadRune()
		}
		// 单字符符号
		punct := convertPunct(res.Char)
		if punct != "" {
			res.Code = punct
			process(res)
			continue
		}
		isHan := unicode.Is(unicode.Han, res.Char)
		if isHan {
			mRes.Dist.LackHan[ch]++
		} else {
			mRes.Dist.NotHan[ch]++
		}
		res.Code = "######"
		process(res)
	}
	return mRes
}
