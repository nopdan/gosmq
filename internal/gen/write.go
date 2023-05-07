package gen

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// 输出赛码表
func Write(dict []*Entry, path string) {
	var buf bytes.Buffer
	buf.Grow(len(dict))
	for _, entry := range dict {
		buf.WriteString(entry.Word)
		buf.WriteByte('\t')
		buf.WriteString(entry.Code)
		if entry.Pos != 1 {
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(entry.Pos))
		}
		buf.WriteByte('\n')
	}
	err := os.WriteFile(path, buf.Bytes(), 0666)
	if err != nil {
		fmt.Println("Warning! 输出赛码表失败：", err)
		return
	}
	fmt.Println("输出赛码表成功：", path)
}
