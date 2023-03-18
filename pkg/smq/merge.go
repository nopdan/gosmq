package smq

func MergeResults(results []*Result, stat bool) *Result {
	ret := newResult()
	ret.DictName = results[0].DictName
	ret.DictLen = results[0].DictLen
	ret.TextName = "总计"
	for _, res := range results {
		if stat {
			for k, v := range res.statData {
				if _, ok := ret.statData[k]; !ok {
					ret.statData[k] = v
				} else {
					ret.statData[k].Count += v.Count
				}
			}
		}
		ret.Basic.Commits += res.Basic.Commits
		ret.Basic.NotHanCount += res.Basic.NotHanCount
		ret.Basic.LackCount += res.Basic.LackCount

		ret.Words.Commits.Count += res.Words.Commits.Count
		ret.Words.Chars.Count += res.Words.Chars.Count
		ret.Words.FirstCount += res.Words.FirstCount
		ret.Collision.Commits.Count += res.Collision.Commits.Count
		ret.Collision.Chars.Count += res.Collision.Chars.Count

		ret.toTalEq10 += res.toTalEq10
		ret.Combs.Count += res.Combs.Count
		ret.Fingers.Same.Count += res.Fingers.Same.Count

		ret.Hands.LL.Count += res.Hands.LL.Count
		ret.Hands.LR.Count += res.Hands.LR.Count
		ret.Hands.RL.Count += res.Hands.RL.Count
		ret.Hands.RR.Count += res.Hands.RR.Count

		ret.Combs.DoubleHit.Count += res.Combs.DoubleHit.Count
		ret.Combs.TribleHit.Count += res.Combs.TribleHit.Count
		ret.Combs.SingleSpan.Count += res.Combs.SingleSpan.Count
		ret.Combs.MultiSpan.Count += res.Combs.MultiSpan.Count
		ret.Combs.LongFingersDisturb.Count += res.Combs.LongFingersDisturb.Count
		ret.Combs.LittleFingersDisturb.Count += res.Combs.LittleFingersDisturb.Count

		for i := 33; i < 128; i++ {
			ret.keysDist[i] += res.keysDist[i]
		}
		for k := range res.notHanMap {
			ret.notHanMap[k] = struct{}{}
		}
		for k := range res.lackMap {
			ret.lackMap[k] = struct{}{}
		}

		for i, v := range res.CodeLen.Dist {
			AddToVal(&ret.CodeLen.Dist, i, v)
		}

		for i, v := range res.Words.Dist {
			AddToVal(&ret.Words.Dist, i, v)
		}
		for i, v := range res.Collision.Dist {
			AddToVal(&ret.Collision.Dist, i, v)
		}
	}
	ret.stat()
	return ret
}
