package smq

func (res *Result) feel(codes string, dict *Dict) {
	if len(codes) == 0 {
		return
	}
	res.mapKeys[codes[0]]++
	last, ok := KEY_POS[codes[0]]
	if !ok {
		last.Key = codes[0]
	}
	for i := 1; i < len(codes); i++ {
		currKey := codes[i]
		// for key
		// 转小写
		if 65 <= currKey && currKey <= 90 {
			currKey += 22
		}
		if currKey == '_' {
			switch dict.PressSpaceBy {
			case "left":
			case "right":
				currKey = '+'
			default: // "both"
				// 如果上一个键是左手
				if !last.LoR {
					currKey = '+'
				}
			}
		}
		res.mapKeys[currKey]++

		current, ok := KEY_POS[currKey]
		if !ok {
			last.Key = currKey
			continue
		}

		// for comb
		comb := string([]byte{last.Key, currKey})
		dl10, ok := COMB.Dl[comb]
		if !ok {
			// log.Printf(`comb nil"%v"%v"%v`, last.key, current, comb)
			last = current
			continue
		}
		res.toTalEq10 += dl10
		res.Combs.Count++

		// for finger
		if current.Fin == last.Fin {
			res.Fingers.Same.Count++
		}
		// for hands
		if current.LoR {
			if last.LoR { // RR
				res.Hands.RR.Count++
			} else { // RL
				res.Hands.RL.Count++
			}
		} else {
			if last.LoR { // LR
				res.Hands.LR.Count++
			} else { // LL
				res.Hands.LL.Count++
			}
		}

		// 同键、三连击
		if currKey == last.Key {
			res.Combs.DoubleHit.Count++
			if i < len(codes)-1 {
				if currKey == codes[i+1] {
					res.Combs.TribleHit.Count++
				}
			}
		}
		// 小跨排
		if COMB.Xkp[comb] {
			res.Combs.SingleSpan.Count++
		}
		// 大跨排
		if COMB.Dkp[comb] {
			res.Combs.MultiSpan.Count++
		}
		// 错手
		if COMB.Cs[comb] {
			res.Combs.LongFingersDisturb.Count++
		}
		// 小拇指干扰
		if COMB.Xzgr[comb] {
			res.Combs.LittleFingersDisturb.Count++
		}
		last = current
	}
}
