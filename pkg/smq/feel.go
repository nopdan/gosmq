package smq

import (
	"github.com/nopdan/gosmq/pkg/feeling"
	"github.com/nopdan/gosmq/pkg/result"
)

type feel struct {
	mRes      *result.MatchRes
	spacePref string

	key        byte
	isLeft     bool
	finger     byte
	lastKey    byte
	lastIsLeft bool
	lastFinger byte
	last2Key   byte
}

func (f *feel) Invalid() {
	f.lastKey = 0
	f.lastIsLeft = false
	f.lastFinger = 0
}

func NewFeeling(target *result.MatchRes, spacePref string) *feel {
	return &feel{mRes: target, spacePref: spacePref}
}

// 处理当前按键，并更新状态。需要在 Process 退出前调用
func (f *feel) step() {
	currKey := f.key
	if currKey == '_' {
		switch f.spacePref {
		case "right":
			currKey = '+'
		case "both", "": // "both"
			// 如果上一个键是左手
			if f.lastFinger != 0 && f.lastIsLeft {
				currKey = '+'
			}
		}
	}
	f.mRes.Dist.Key[currKey]++
	f.last2Key, f.lastKey = f.lastKey, f.key
	f.lastIsLeft, f.lastFinger = f.isLeft, f.finger
}

// 传入的 key 必须为 a-z0-9,./;'[]-= 中的一个
//
// 特别的，传入大写字母自动转为小写，传入空格_，处理右手击键为+
func (f *feel) Process(key byte) {
	mRes := f.mRes
	// 跳过
	if key >= 128 {
		return
	}
	// magic: 将大写字母转为小写
	if 'A' <= key && key <= 'Z' {
		key |= 32
	}

	f.key = key
	f.isLeft, f.finger = feeling.KeyPos(f.key)
	// 如果当前键或者上一个键不合法(不在46键里)
	if f.lastKey == 0 || f.finger == 0 {
		f.step()
		return
	}

	comb := feeling.Combination[f.lastKey][f.key]
	// 当量表里找不到
	if comb == nil {
		mRes.Equivalent += 2.0
		mRes.Pair.Count++
		f.step()
		return
	}

	// 左右手分布
	if f.lastIsLeft {
		if f.isLeft {
			mRes.Pair.LeftToLeft++
		} else {
			mRes.Pair.LeftToRight++
		}
	} else {
		if f.isLeft {
			mRes.Pair.RightToLeft++
		} else {
			mRes.Pair.RightToRight++
		}
	}

	// 同指
	if f.finger == f.lastFinger {
		mRes.Pair.SameFinger++
	}
	// 同键、三连击
	if f.key == f.lastKey {
		mRes.Pair.DoubleHit++
		if f.key == f.last2Key {
			mRes.Pair.TribleHit++
		}
	}
	// 小跨排
	if comb.SingleSpan {
		mRes.Pair.SingleSpan++
	}
	// 大跨排
	if comb.MultiSpan {
		mRes.Pair.MultiSpan++
	}
	// 错手
	if comb.Staggered {
		mRes.Pair.Staggered++
	}
	// 小拇指干扰
	if comb.Disturb {
		mRes.Pair.Disturb++
	}

	mRes.Equivalent += comb.Equivalent
	mRes.Pair.Count++
	f.step()
	return
}
