package transformer

import (
	"io"
	"strconv"
)

type Dict struct {
	Name       string
	Reader     io.Reader
	PushStart  int
	SelectKeys string
}

type Entry struct {
	Word  string
	Code  string
	Order int
}

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
