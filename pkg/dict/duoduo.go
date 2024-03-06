package dict

import (
	"bufio"
	"slices"
	"strings"

	"github.com/nopdan/gosmq/pkg/util"
)

// 加载多多或者冰凌码表
func (d *Dict) loadTSV(wordFirst bool) []*Entry {
	var cap int = 1e5
	if d.Size > 0 {
		cap = d.Size / 32
	}
	ret := make([]*Entry, 0, cap)
	// 统计编码出现的次数
	stat := make(map[string]int)
	var word, code string
	var pos int
	scan := bufio.NewScanner(d.Reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		if wordFirst {
			word, code = wc[0], wc[1]
		} else {
			word, code = wc[1], wc[0]
		}
		stat[code]++
		pos = stat[code]
		code = d.addSuffix(code, pos)
		ret = append(ret, &Entry{word, code, pos})
		d.insert(word, code, pos)
	}
	return ret
}

// 加上选重键，pos 是编码出现的次数，最小为 1
func (d *Dict) addSuffix(code string, pos int) string {
	// 大于等于起顶码长，首选不用添加空格 _
	if pos == 1 && len(code) >= d.push {
		return code
	}
	// 添加自定义选重键
	tmp := util.UnsafeToBytes(code)
	tmp = slices.Concat(tmp, d.getSelectKey(pos))
	return util.UnsafeToString(tmp)
}
