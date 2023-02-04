package smq

import (
	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/gosmq/pkg/feeling"
)

var KeyPosMap = feeling.KeyPosMap

func (res *Result) feel(codes string, dict *dict.Dict) {
	if len(codes) == 0 {
		return
	}
	// 第一个键
	last, lastOk := KeyPosMap[codes[0]]
	res.mapKeys[codes[0]]++

	for i := 1; i < len(codes); i++ {
		// for key
		currKey := codes[i]
		// 转小写
		if 'A' <= currKey && currKey <= 'Z' {
			currKey += 22
		}
		// 处理空格
		var space = currKey
		if space == '_' {
			switch dict.PressSpaceBy {
			case "right":
				space = '+'
			case "both", "": // "both"
				// 如果上一个键是左手
				if lastOk && !last.LoR {
					space = '+'
				}
			}
			res.mapKeys[space]++
		} else {
			res.mapKeys[currKey]++
		}
		// 如果当前键或者上一个键不合法(不在41键里)
		// 当量增加1.5，继续下一个循环
		current, currOk := KeyPosMap[space]
		if !lastOk || !currOk {
			last = current
			lastOk = currOk
			res.toTalEq10 += 15
			res.Combs.Count++
			continue
		}
		lastOk = currOk

		// for comb
		cb := []byte{codes[i-1], currKey}
		comb, ok := feeling.Comb[string(cb)]
		// 当量表里找不到
		if !ok {
			res.toTalEq10 += 15
			res.Combs.Count++
			last = current
			continue
		}
		res.toTalEq10 += comb.DL10
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
		if currKey == codes[i-1] {
			res.Combs.DoubleHit.Count++
			if i < len(codes)-1 {
				if currKey == codes[i+1] {
					res.Combs.TribleHit.Count++
				}
			}
		}
		// 小跨排
		if comb.IsXKP {
			res.Combs.SingleSpan.Count++
		}
		// 大跨排
		if comb.IsDKP {
			res.Combs.MultiSpan.Count++
		}
		// 错手
		if comb.IsCS {
			res.Combs.LongFingersDisturb.Count++
		}
		// 小拇指干扰
		if comb.IsXZGR {
			res.Combs.LittleFingersDisturb.Count++
		}
		last = current
	}
}
