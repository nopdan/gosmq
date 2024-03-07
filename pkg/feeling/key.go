package feeling

// 按键分布，左手高位为 1
var keyDistrib [128]byte

func init() {
	// 手指
	keys := "12345qwertasdfgzxcvb_+67890-=yuiop[]hjkl;'nm,./"
	for i := range len(keys) {
		keyDistrib[keys[i]] = finger(keys[i])
	}
	// 左手空格"_" 右手空格"+"
	lk := "12345qwertasdfgzxcvb_"
	for i := range len(lk) {
		keyDistrib[lk[i]] |= 0b10000000
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
	case '4', 'r', 'f', 'v', '5', 't', 'g', 'b':
		return 4
	case '_':
		return 5
	case '+':
		return 6
	case '6', 'y', 'h', 'n', '7', 'u', 'j', 'm':
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
