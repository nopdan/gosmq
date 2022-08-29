package transformer

import (
	"bufio"
	"strconv"
	"strings"
)

type Duoduo struct {
	Reverse bool
}

func (d Duoduo) Read(dict Dict) []Entry {
	ret := make([]Entry, 1e5)
	mapOrder := make(map[string]int)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		var w, c string
		if d.Reverse {
			w, c = wc[1], wc[0]
		} else {
			w, c = wc[0], wc[1]
		}

		mapOrder[c]++
		order := mapOrder[c]
		// 生成赛码表
		if len(c) >= dict.PushStart && order == 1 {
		} else {
			if order <= len(dict.SelectKeys) {
				c += string(dict.SelectKeys[order-1])
			} else {
				c += strconv.Itoa(order)
			}
		}
		ret = append(ret, Entry{w, c, order})
	}
	return ret
}
