package smq

import (
	"bufio"
	"io"
)

// 分割文本
func SplitStep(brd *bufio.Reader, bufLen int) ([]byte, error) {
	var text []byte

	buffer := make([]byte, bufLen)
	n, err := io.ReadFull(brd, buffer)
	buffer = buffer[:n]
	for {
		b, err := brd.ReadByte()
		// 控制字符 直接切分
		if b < 33 {
			text = buffer
			break
		}
		// utf-8 前缀
		if b >= 0b11000000 {
			brd.UnreadByte()
		} else {
			buffer = append(buffer, b)
		}
		// EOF
		if err != nil {
			text = buffer
			break
		}
		// 读到合法字符，开始读 rune
		if b < 128 || b >= 0b11000000 {
			text = buffer
		OUT:
			// 超过限制读不到分割符直接 break
			for lim := 0; lim < bufLen; lim++ {
				rn, _, err := brd.ReadRune()
				if rn < 33 {
					break
				}
				text = append(text, []byte(string(rn))...)
				if err != nil {
					break
				}
				switch rn {
				case '。', '？', '！', '》':
					break OUT
				}
			}
			break
		}
	}
	return text, err
}
