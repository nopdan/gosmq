package transformer

import (
	"bufio"
	"strconv"
	"strings"
)

type Jidian struct{}

func (j Jidian) Read(dict Dict) []Entry {
	ret := make([]Entry, 0, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), " ")
		if len(wc) < 2 {
			continue
		}
		// 单字模式修正
		revise := 0
		c := wc[0]
		for i := 1; i < len(wc); i++ {
			if dict.Single && len([]rune(wc[i])) != 1 {
				revise++
				continue
			}
			order := i - revise
			// 生成赛码表
			var code string
			if len(wc[0]) >= dict.PushStart && order == 1 {
			} else {
				if int(order) <= len(dict.SelectKeys) {
					code = c + string(dict.SelectKeys[order-1])
				} else {
					code = c + strconv.Itoa(int(order))
				}
			}
			ret = append(ret, Entry{wc[i], code, order})
		}
	}
	return ret
}
