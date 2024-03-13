package smq

import (
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/matcher"
	"github.com/nopdan/gosmq/pkg/result"
	"github.com/nopdan/gosmq/pkg/util"
)

func (c *Config) match(buffer []byte, partIdx int, dict *data.Dict) *result.MatchRes {
	mRes := result.NewMatchRes()
	mRes.Segment.PartIdx = partIdx

	feel := NewFeeling(mRes, dict.SpacePref)
	brd := bytes.NewReader(buffer)
	res := new(matcher.Result)

	hanHandler := func(ch rune) {
		isHan := unicode.Is(unicode.Han, ch)
		if isHan {
			mRes.Dist.LackHan[ch]++
		} else {
			mRes.Dist.NotHan[ch]++
		}
	}

	process := func(res *matcher.Result) {
		mRes.Commit.Count++
		mRes.Char.Count += res.Length
		util.Increase(&mRes.Dist.WordLen, res.Length)
		util.Increase(&mRes.Dist.Collision, res.Pos)
		util.Increase(&mRes.Dist.CodeLen, len(res.Code))
		if res.Code == "######" {
			feel.Invalid()
		} else {
			for i := range len(res.Code) {
				feel.Process(res.Code[i])
			}
		}
		if res.Pos >= 2 {
			mRes.Commit.Collision++
			mRes.Char.Collision += res.Length
		}
		// 匹配到词组
		if res.Length >= 2 {
			mRes.Commit.Word++
			mRes.Char.Word += res.Length
			if res.Pos == 1 {
				mRes.Commit.WordFirst++ // 首选词
				mRes.Char.WordFirst += res.Length
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
			if mRes.Segment.Builder == nil {
				mRes.Segment.Builder = new(strings.Builder)
			}
			sb := mRes.Segment.Builder
			sb.WriteString(word)
			sb.WriteByte('\t')
			sb.WriteString(res.Code)
			sb.WriteByte('\n')
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

	// 处理两字宽标点符号 破折号 —— 省略号 ……
	_2Width := func(res *matcher.Result, brd *bytes.Reader) bool {
		if res.Char != '—' && res.Char != '…' {
			return false
		}
		ch2, _, err := brd.ReadRune()
		if err == nil || res.Char != ch2 {
			return false
		}
		// 不计打词 Length 保持 1
		if res.Char == '—' {
			// 中文破折号 —— 占用 6 字节，不计打词
			res.SetChar(0).SetSize(6).SetPos(1)
			res.Code = "=-"
			process(res)
		} else if res.Char == '…' {
			// 中文省略号 …… 占用 6 字节，不计打词
			res.SetChar(0).SetSize(6).SetPos(1)
			res.Code = "=6"
			process(res)
		} else {
			_ = brd.UnreadRune()
			return false
		}
		return true
	}

	for brd.Len() > 0 {
		// 开始匹配
		dict.Matcher.Match(brd, res)
		mRes.TextLen += res.Length

		// 匹配成功
		if res.Pos > 0 {
			process(res)
			continue
		}

		// 匹配失败了
		if c.Clean {
			feel.Invalid()
			continue
		}

		// 跳过空白符
		if unicode.IsSpace(res.Char) {
			continue
		}
		hanHandler(res.Char)

		// 单字符符号
		punct := convertPunct(res.Char)
		if punct != "" {
			res.Code = punct
			res.Pos = 1
			process(res)
			continue
		}

		// 两个字符的符号
		if ok := _2Width(res, brd); ok {
			continue
		}

		res.Code = "######"
		res.Pos = 0
		process(res)
	}
	return mRes
}
