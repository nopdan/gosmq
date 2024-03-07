package result

import (
	"sort"

	"github.com/nopdan/gosmq/pkg/feeling"
)

func div(x, y int) float64 {
	return float64(x) / float64(y)
}

func (m *MatchRes) Stat(info *Info) *Result {
	res := new(Result)
	res.segments = m.segments
	res.statData = m.StatData
	res.Info = *info
	res.Keys = make(map[string]*CountRate)
	res.Commit = m.Commit
	res.Pair = m.Pair
	res.Dist.CodeLen = m.Dist.CodeLen
	res.Dist.WordLen = m.Dist.WordLen
	res.Dist.Collision = m.Dist.Collision

	// 文章字数
	for i, v := range res.Dist.WordLen {
		res.Info.TextLen += i * v
	}
	// 总码长
	for i, v := range res.Dist.CodeLen {
		res.CodeLen.Total += i * v
	}
	// 字均码长
	res.CodeLen.PerChar = div(res.CodeLen.Total, res.Info.TextLen)
	// 按键分布
	for i := byte(33); i < 128; i++ {
		cr := new(CountRate)
		cr.Count = m.Dist.Key[i]
		if cr.Count == 0 {
			continue
		}
		cr.Rate = div(cr.Count, res.CodeLen.Total)
		res.Keys[string(i)] = cr
		// 左右手
		isLeft, finger := feeling.KeyPos(i)
		if isLeft {
			res.LeftHand += cr.Count
		} else {
			res.RightHand += cr.Count
		}
		// 指法
		res.Dist.Finger[finger] += cr.Count
	}
	res.Pair.SameHand = res.Pair.LeftToLeft + res.Pair.RightToRight
	res.Pair.DiffHand = res.Pair.Count - res.Pair.SameHand
	res.Pair.DiffFinger = res.Pair.SameHand - res.Pair.SameFinger

	// 非汉字
	res.Han.NotHans = len(m.Dist.NotHan)
	notHanList := make([]rune, 0, len(m.Dist.NotHan))
	for k, v := range m.Dist.NotHan {
		res.Han.NotHanCount += v
		notHanList = append(notHanList, k)
	}
	sort.Slice(notHanList, func(i, j int) bool {
		return notHanList[i] < notHanList[j]
	})
	res.Han.NotHan = string(notHanList)

	// 缺字
	res.Han.Lacks = len(m.Dist.LackHan)
	lackList := make([]rune, 0, len(m.Dist.LackHan))
	for k, v := range m.Dist.LackHan {
		res.Han.LackCount += v
		lackList = append(lackList, k)
	}
	sort.Slice(lackList, func(i, j int) bool {
		return lackList[i] < lackList[j]
	})
	res.Han.Lack = string(lackList)
	return res
}
