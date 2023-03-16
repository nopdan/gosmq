package feeling

type KeyPos struct {
	IsLeft bool

	// 0-10, 1对应1qaz, 10对应0p;/, 56对应左右大拇指_+, 0表示错误的键
	Fin byte
}

// 获取按键左右手以及用哪根手指击键
//
// 左右大拇指按空格分别用 _ + 表示
var KeyPosArr [256]KeyPos

func init() {
	// 左手空格"_" 右手空格"+"
	lk := "_12345qwertasdfgzxcvb"  // left keys
	rk := "+67890yuiophjkl;'nm,./" // right keys
	for i := range lk {
		KeyPosArr[lk[i]] = KeyPos{true, 0}
	}
	for i := range rk {
		KeyPosArr[rk[i]] = KeyPos{false, 0}
	}
	// 手指
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol."
	fins := "11112222333344444444567777777788889999"
	for i := range keys {
		tmp := KeyPosArr[keys[i]]
		tmp.Fin = fins[i] - '0'
		KeyPosArr[keys[i]] = tmp
	}
	keys = "0p;/'"
	for i := range keys {
		tmp := KeyPosArr[keys[i]]
		tmp.Fin = 10
		KeyPosArr[keys[i]] = tmp
	}
}
