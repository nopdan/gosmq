package smq

import (
	_ "embed"
	"strconv"
	"strings"
)

var zhKeysMap = map[rune]string{
	'，': ",",
	'。': ".",
	'、': "/",
	'；': ";",
	'‘': "'",
	'’': "'",
	'【': "[",
	'】': "]",
	'·': "`",

	'《': "=,",
	'》': "=.",
	'？': "=/",
	'：': "=;",
	'“': "='",
	'”': "='",
	'！': "=1",
	'￥': "=4",
	'（': "=9",
	'）': "=0",
}

var enKeysMap [128]string

func init() {
	baseKeys := ",./;[]-='"
	shiftKeys := "<>?:{}_+\""
	numPuncts := ")!@#$%^&*("
	for b := byte(33); b < 128; b++ {
		if 'A' <= b && b <= 'Z' {
			// magic: 将英文字符转换为小写
			enKeysMap[b] = string([]byte{'=', b | ' '})
		} else if idx := strings.IndexByte(shiftKeys, b); idx != -1 {
			// shift 符号
			enKeysMap[b] = string([]byte{'=', baseKeys[idx]})
		} else if idx = strings.IndexByte(numPuncts, b); idx != -1 {
			// shift+数字 符号
			enKeysMap[b] = "=" + strconv.Itoa(idx)
		} else {
			enKeysMap[b] = string(b)
		}
	}
}

// = 作为 shift
func convertPunct(char rune) string {
	// 中文标点符号
	if char >= 128 {
		if v, ok := zhKeysMap[char]; ok {
			return v
		}
		// ascii 内全角转半角
		if 0xFF01 <= char && char <= 0xFF5E {
			char -= 0xFEE0
		} else {
			return ""
		}
	}

	// 英文标点
	b := byte(char)
	return enKeysMap[b]
}
