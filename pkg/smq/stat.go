package smq

import (
	"sort"
)

func (res *Result) stat() {
	for i, v := range res.Words.Dist {
		res.Basic.TextLen += i * v
	}
	res.TextLen = res.Basic.TextLen

	res.Name = res.DictName         // 旧的
	res.Basic.DictLen = res.DictLen // 旧的
	// Basic

	// 非汉字
	res.Basic.NotHans = len(res.notHanMap)
	notHanList := make([]rune, 0, len(res.notHanMap))
	for k := range res.notHanMap {
		notHanList = append(notHanList, k)
	}
	sort.Slice(notHanList, func(i, j int) bool {
		return notHanList[i] < notHanList[j]
	})
	res.Basic.NotHan = string(notHanList)

	// 缺字
	res.Basic.Lacks = len(res.lackMap)
	lackList := make([]rune, 0, len(res.lackMap))
	for k := range res.lackMap {
		lackList = append(lackList, k)
	}
	sort.Slice(lackList, func(i, j int) bool {
		return lackList[i] < lackList[j]
	})
	res.Basic.Lack = string(lackList)

	// 打词数
	res.Words.Commits.Rate = div(res.Words.Commits.Count, res.Basic.Commits)
	res.Words.Chars.Rate = div(res.Words.Chars.Count, res.Basic.TextLen)
	// 选重数
	res.Collision.Commits.Rate = div(res.Collision.Commits.Count, res.Basic.Commits)
	res.Collision.Chars.Rate = div(res.Collision.Chars.Count, res.Basic.TextLen)
	// 码长
	for i, v := range res.CodeLen.Dist {
		res.CodeLen.Total += i * v
	}
	res.CodeLen.PerChar = div(res.CodeLen.Total, res.Basic.TextLen)

	res.statFeel()
}

func (res *Result) statFeel() {
	// keys
	for i := 33; i < 128; i++ {
		key := string(byte(i))
		if _, ok := res.Keys[key]; !ok {
			res.Keys[key] = new(CountRate)
		}
		res.Keys[key].Count = res.keysDist[i]
	}
	for _, v := range res.Keys {
		v.Rate = div(v.Count, res.CodeLen.Total)
	}
	// combs
	res.Combs.Equivalent = div(res.toTalEq10/10, res.Combs.Count)
	res.Combs.DoubleHit.Rate = div(res.Combs.DoubleHit.Count, res.Combs.Count)
	res.Combs.TribleHit.Rate = div(res.Combs.TribleHit.Count, res.Combs.Count)
	res.Combs.SingleSpan.Rate = div(res.Combs.SingleSpan.Count, res.Combs.Count)
	res.Combs.MultiSpan.Rate = div(res.Combs.MultiSpan.Count, res.Combs.Count)
	res.Combs.LongFingersDisturb.Rate = div(res.Combs.LongFingersDisturb.Count, res.Combs.Count)
	res.Combs.LittleFingersDisturb.Rate = div(res.Combs.LittleFingersDisturb.Count, res.Combs.Count)
	// hands
	res.Hands.Left.Count = res.Hands.LL.Count + res.Hands.LR.Count/2 + res.Hands.RL.Count/2
	res.Hands.Right.Count = res.Combs.Count - res.Hands.Left.Count
	res.Hands.Same.Count = res.Hands.LL.Count + res.Hands.RR.Count
	res.Hands.Diff.Count = res.Combs.Count - res.Hands.Same.Count
	res.Hands.Left.Rate = div(res.Hands.Left.Count, res.Combs.Count)
	res.Hands.Right.Rate = div(res.Hands.Right.Count, res.Combs.Count)
	res.Hands.Same.Rate = div(res.Hands.Same.Count, res.Combs.Count)
	res.Hands.Diff.Rate = div(res.Hands.Diff.Count, res.Combs.Count)
	// fingers
	for i := 33; i < 128; i++ {
		if keyPos := KeyPosArr[i]; keyPos.Fin == 0 {
			res.Fingers.Dist[10].Count += res.keysDist[i]
		} else if keyPos.Fin == 10 {
			res.Fingers.Dist[0].Count += res.keysDist[i]
		} else {
			res.Fingers.Dist[keyPos.Fin].Count += res.keysDist[i]
		}
	}
	for i := range res.Fingers.Dist {
		v := &res.Fingers.Dist[i]
		v.Rate = div(v.Count, res.CodeLen.Total)
	}
	res.Fingers.Same.Rate = div(res.Fingers.Same.Count, res.Combs.Count)
	res.Fingers.Diff.Count = res.Hands.Same.Count - res.Fingers.Same.Count
	res.Fingers.Diff.Rate = div(res.Fingers.Diff.Count, res.Combs.Count)
	res.Hands.LL.Rate = div(res.Hands.LL.Count, res.Combs.Count)
	res.Hands.LR.Rate = div(res.Hands.LR.Count, res.Combs.Count)
	res.Hands.RL.Rate = div(res.Hands.RL.Count, res.Combs.Count)
	res.Hands.RR.Rate = div(res.Hands.RR.Count, res.Combs.Count)
}
