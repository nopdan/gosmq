package gen

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/nopdan/ku"
)

// 加载多多码表
func (c *Config) LoadTSV(wordFirst bool) []*Entry {
	ret := make([]*Entry, 0, 1e5)
	rd, err := ku.Read(c.Path)
	if err != nil {
		panic(err)
	}
	// 统计编码出现的次数
	stat := make(map[string]int)
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		word, code := wc[0], wc[1]
		if !wordFirst {
			word, code = code, word
		}
		stat[code]++
		pos := stat[code]
		code = c.addSuffix(code, pos)
		ret = append(ret, &Entry{word, code, pos})
	}
	return ret
}

// 加上选重键，pos 是编码出现的次数，最小为 1
func (c Config) addSuffix(code string, pos int) string {
	// 大于等于起顶码长，首选不用添加空格 _
	if len(code) >= c.PushStart {
		if pos == 1 {
			return code
		}
	}

	// 添加自定义选重键
	if pos <= len(c.SelectKeys) {
		code += string(c.SelectKeys[pos-1])
	} else {
		code += strconv.Itoa(pos)
	}
	return code
}
