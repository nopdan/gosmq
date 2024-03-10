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
	res.Commit = m.Commit
	res.Char = m.Char
	res.Pair = m.Pair
	res.Dist.CodeLen = m.Dist.CodeLen
	res.Dist.WordLen = m.Dist.WordLen
	res.Dist.Collision = m.Dist.Collision
	res.Dist.Key = make(map[string]int)

	res.Info.TextLen = m.TextLen
	// 总码长 == 总按键数
	for i, v := range res.Dist.CodeLen {
		res.Keys.Count += i * v
	}
	// 字均码长 = 总按键数 / 总字数
	res.Keys.CodeLen = div(res.Keys.Count, res.Char.Count)
	// 按键分布
	for i := byte(33); i < 128; i++ {
		count := m.Dist.Key[i]
		if count == 0 {
			continue
		}
		res.Dist.Key[string(i)] = count
		// 左右手
		isLeft, finger := feeling.KeyPos(i)
		if isLeft {
			res.Keys.LeftHand += count
		} else {
			res.Keys.RightHand += count
		}
		// 指法
		res.Dist.Finger[finger] += count
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
