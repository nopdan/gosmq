package dict

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// 加上选重键
func (dict *Dict) getRealCode(c string, order int) string {
	if order != 1 || len(c) < dict.PushStart {
		if order <= len(dict.SelectKeys) {
			c += string(dict.SelectKeys[order-1])
		} else {
			c += strconv.Itoa(order)
		}
	}
	return c
}

// 输出赛码表
func outputDict(t []Entry, name string) {
	var buf bytes.Buffer
	buf.Grow(1e5)
	for i := range t {
		buf.WriteString(t[i].Word)
		buf.WriteByte('\t')
		buf.WriteString(t[i].Code)
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(t[i].Order))
		buf.WriteByte('\n')
	}
	path := "dict/_" + name + ".txt"
	err := os.WriteFile(path, buf.Bytes(), 0666)
	if err != nil {
		fmt.Println("Warning! 输出赛码表失败：", err)
		return
	}
	fmt.Println("输出赛码表成功：", path)
}
