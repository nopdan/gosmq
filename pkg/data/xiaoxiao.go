package data

import (
	"bufio"
	"strings"
)

// 小小|极点
func (d *Dict) loadXiao() []*Entry {
	var cap int = 1e5
	if d.Text.size > 0 {
		cap = d.Text.size / 32
	}
	ret := make([]*Entry, 0, cap)
	scan := bufio.NewScanner(d.Text.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), " ")
		if len(wc) < 2 {
			continue
		}
		code := wc[0]
		for i := 1; i < len(wc); i++ {
			word := wc[i]
			code = d.addSuffix(code, i)
			ret = append(ret, &Entry{word, code, i})
			d.insert(word, code, i)
		}
	}
	return ret
}
