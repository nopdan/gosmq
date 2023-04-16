package gen

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/flowerime/rose/pkg/rose"
)

// 输出赛码表
func Write(t rose.WubiTable, path string) {
	var buf bytes.Buffer
	buf.Grow(len(t))
	for i := range t {
		buf.WriteString(t[i].Word)
		buf.WriteByte('\t')
		buf.WriteString(t[i].Code)
		if t[i].Pos != 1 {
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(t[i].Pos))
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
