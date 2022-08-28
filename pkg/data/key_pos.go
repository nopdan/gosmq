package data

type KeyPos struct {
	Key byte
	LoR bool // Left or Right
	Fin int  // 0-10, 10是错误的，左右大拇指_+
}

func GetKeyPos() map[byte]KeyPos {
	ret := make(map[byte]KeyPos)
	// 左手空格"_" 右手空格"+"
	lk := "_12345qwertasdfgzxcvb"  // left keys
	rk := "+67890yuiophjkl;'nm,./" // right keys
	for i := range lk {
		ret[lk[i]] = KeyPos{lk[i], false, 0}
	}
	for i := range rk {
		ret[rk[i]] = KeyPos{rk[i], true, 0}
	}
	// 手指
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol.0p;/'"
	fin := "1111222233334444444456777777778888999900000"
	for i := range keys {
		tmp := ret[keys[i]]
		tmp.Fin = int(fin[i] - 48)
		ret[keys[i]] = tmp
	}
	return ret
}
