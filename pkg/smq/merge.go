package smq

// 合并一段文本多个码表生成的结果
func mergeRes(resA, resB []*Result, dicts []*Dict) {
	if len(resA) != len(resB) || len(resA) != len(dicts) {
		panic("程序运行错误：merge，两个结果长度不相等，请联系作者")
	}
	for j, res := range resA {
		if dicts[j].Stat {
			for k, v := range resB[j].statData {
				if _, ok := res.statData[k]; !ok {
					res.statData[k] = v
				} else {
					res.statData[k].Count += v.Count
				}
			}
		}
		res.Basic.Commits += resB[j].Basic.Commits
		res.Basic.NotHanCount += resB[j].Basic.NotHanCount
		res.Basic.LackCount += resB[j].Basic.LackCount

		res.Words.Commits.Count += resB[j].Words.Commits.Count
		res.Words.Chars.Count += resB[j].Words.Chars.Count
		res.Words.FirstCount += resB[j].Words.FirstCount
		res.Collision.Commits.Count += resB[j].Collision.Commits.Count
		res.Collision.Chars.Count += resB[j].Collision.Chars.Count

		res.toTalEq10 += resB[j].toTalEq10
		res.Combs.Count += resB[j].Combs.Count
		res.Fingers.Same.Count += resB[j].Fingers.Same.Count

		res.Hands.LL.Count += resB[j].Hands.LL.Count
		res.Hands.LR.Count += resB[j].Hands.LR.Count
		res.Hands.RL.Count += resB[j].Hands.RL.Count
		res.Hands.RR.Count += resB[j].Hands.RR.Count

		res.Combs.DoubleHit.Count += resB[j].Combs.DoubleHit.Count
		res.Combs.TribleHit.Count += resB[j].Combs.TribleHit.Count
		res.Combs.SingleSpan.Count += resB[j].Combs.SingleSpan.Count
		res.Combs.MultiSpan.Count += resB[j].Combs.MultiSpan.Count
		res.Combs.LongFingersDisturb.Count += resB[j].Combs.LongFingersDisturb.Count
		res.Combs.LittleFingersDisturb.Count += resB[j].Combs.LittleFingersDisturb.Count

		for i := 33; i < 128; i++ {
			res.keysDist[i] += resB[j].keysDist[i]
		}
		for k := range resB[j].notHanMap {
			res.notHanMap[k] = struct{}{}
		}
		for k := range resB[j].lackMap {
			res.lackMap[k] = struct{}{}
		}

		for i, v := range resB[j].CodeLen.Dist {
			AddToVal(&res.CodeLen.Dist, i, v)
		}

		for i, v := range resB[j].Words.Dist {
			AddToVal(&res.Words.Dist, i, v)
		}
		for i, v := range resB[j].Collision.Dist {
			AddToVal(&res.Collision.Dist, i, v)
		}
	}
}
