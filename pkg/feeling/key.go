package feeling

// 按键分布，左手高位为 1
var keyDistrib [128]byte

func init() {
	// 左手空格"_" 右手空格"+"
	lk := "1qaz2wsx3edc4rfv5tgb_" // left keys
	// _ = "+6yhn7ujm8ik,9ol.0p;/'[]-=" // right keys
	for i := range lk {
		keyDistrib[lk[i]] += 128
	}
	// 手指
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol.0p;/'[]-="
	for i := range keys {
		keyDistrib[keys[i]] += finger(keys[i])
	}
}

// 按键左右手，左手为 true
//
// 手指分布，左右大拇指按空格分别用 _ + 表示
func KeyPos(key byte) (isLeft bool, finger byte) {
	b := keyDistrib[key]
	return b > 127, b & 0b01111111
}

func finger(key byte) byte {
	switch key {
	case '1', 'q', 'a', 'z':
		return 1
	case '2', 'w', 's', 'x':
		return 2
	case '3', 'e', 'd', 'c':
		return 3
	case '4', 'r', 'f', 'v':
		return 4
	case '5', 't', 'g', 'b', '_':
		return 5
	case '6', 'y', 'h', 'n', '+':
		return 6
	case '7', 'u', 'j', 'm':
		return 7
	case '8', 'i', 'k', ',':
		return 8
	case '9', 'o', 'l', '.':
		return 9
	case '0', 'p', ';', '/', '\'', '[', ']', '-', '=':
		return 10
	default:
		return 0
	}
}
