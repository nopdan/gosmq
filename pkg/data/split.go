package data

import (
	"io"

	"github.com/nopdan/gosmq/pkg/util"
)

// 从Text读取器中读取，直到达到一个控制字符、特定的UTF-8字符，或者超过缓冲区的容量。
//
// 最后一次读取返回 io.EOF
func (t *Text) Iter() ([]byte, error) {
	if t.reader == nil {
		return nil, io.EOF
	}
	if t.size < t.bufSize {
		return io.ReadAll(t.reader)
	}
	buffer := make([]byte, t.bufSize, t.bufSize+4*1024)
	n, _ := io.ReadFull(t.reader, buffer)
	buffer = buffer[:n]

	for {
		b, err := t.reader.ReadByte()
		// EOF
		if err != nil {
			return buffer, io.EOF
		}
		// 防止切到 utf-8 编码中间
		if b < 33 {
			// 控制字符 直接切分
			return buffer, nil
		} else if b < 0b11000000 { // 0b0xxxxxxx || 0b10xxxxxx
			// ascii 码，或者 utf-8 编码后几位
			buffer = append(buffer, b)
			continue
		} else { // 0b11xxxxxx
			// utf-8 前缀，回退，之后读取 rune
			_ = t.reader.UnreadByte()
			break
		}
	}

	for {
		r, _, _ := t.reader.ReadRune()
		// 控制字符 直接切分
		if r < 33 {
			return buffer, nil
		}
		switch r {
		case '“', '‘', '：', '《':
			_ = t.reader.UnreadRune()
			return buffer, nil
		}
		b := util.UnsafeToBytes(string(r))
		// 超过 buffer 容量直接返回，减少切片扩容
		if len(buffer)+len(b) > cap(buffer) {
			_ = t.reader.UnreadRune()
			return buffer, nil
		}
		buffer = append(buffer, b...)
		switch r {
		case '。', '？', '！', '》':
			return buffer, nil
		}
	}
}
