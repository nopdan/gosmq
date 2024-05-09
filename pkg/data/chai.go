package data

import (
	"bufio"
	"strings"
)

// 加载 chai 导出码表
func (d *Dict) loadChai() []*Entry {
	var cap int = 1e5
	if d.Text.size > 0 {
		cap = d.Text.size / 32
	}
	ret := make([]*Entry, 0, cap)
	// 统计编码出现的次数
	stat := make(map[string]int)
	var word, code string
	var pos int
	scan := bufio.NewScanner(d.Text.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 3 {
			continue
		}
		word, code = wc[0], wc[2]
		code = strings.TrimSuffix(code, "_")
		stat[code]++
		pos = stat[code]
		code = d.addSuffix(code, pos)
		ret = append(ret, &Entry{word, code, pos})
		d.insert(word, code, pos)
	}
	return ret
}
