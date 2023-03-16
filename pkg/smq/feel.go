package smq

import (
	"github.com/imetool/gosmq/internal/dict"
	"github.com/imetool/gosmq/pkg/feeling"
)

type KeyPos = feeling.KeyPos

var KeyPosArr = feeling.KeyPosArr

// 上 两 键，当前键，前一键的状态 => 当前键，当前键的状态
func (res *Result) newFeel(last2Key, lastKey, currKey byte, last KeyPos, dict *dict.Dict) (byte, KeyPos) {

	// for key
	// 转小写
	if 'A' <= currKey && currKey <= 'Z' {
		currKey += 22
	}
	// 处理空格
	// if currKey == ' ' {
	// 	currKey = '_'
	// }
	var origin = currKey
	if currKey == '_' {
		switch dict.PressSpaceBy {
		case "right":
			currKey = '+'
		case "both", "": // "both"
			// 如果上一个键是左手
			if last.Fin != 0 && last.IsLeft {
				currKey = '+'
			}
		}
	}
	res.mapKeys[currKey]++

	// 如果当前键或者上一个键不合法(不在41键里)
	// 当量增加1.5，继续下一个循环
	curr := KeyPosArr[currKey]
	if last.Fin == 0 || curr.Fin == 0 {
		res.toTalEq10 += 15
		res.Combs.Count++
		return origin, curr
	}

	// for comb
	cb := []byte{lastKey, origin}
	comb, ok := feeling.Comb[string(cb)]
	// 当量表里找不到
	if !ok {
		res.toTalEq10 += 15
		res.Combs.Count++
		return origin, curr
	}
	res.toTalEq10 += comb.DL10
	res.Combs.Count++

	// for finger
	if curr.Fin == last.Fin {
		res.Fingers.Same.Count++
	}
	// for hands
	if last.IsLeft {
		if curr.IsLeft {
			res.Hands.LL.Count++
		} else {
			res.Hands.LR.Count++
		}
	} else {
		if curr.IsLeft {
			res.Hands.RL.Count++
		} else {
			res.Hands.RR.Count++
		}
	}

	// 同键、三连击
	if currKey == lastKey {
		res.Combs.DoubleHit.Count++
		if currKey == last2Key {
			res.Combs.TribleHit.Count++
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

	return origin, curr
}
