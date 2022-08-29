package transformer

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

type Jisu struct{}

func (j Jisu) Read(dict Dict) []Entry {
	ret := make([]Entry, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		c := wc[1]
		code := ""
		order := 0
		// a_ aa_
		if len(c)-1 > 0 && c[len(c)-1] == '_' {
			code = c[:len(c)-1]
			order = 1
		} else {
			re := regexp.MustCompile(`\d+$`)
			match := re.FindString(c)
			// a1 aa3
			if match != "" {
				code = c[:len(c)-len(match)]
				order, _ = strconv.Atoi(match)
			} else { // akdb ksdw
				code = c
				order = 0 // 和前面区分开
			}
		}
		// 生成赛码表
		// 自定义选重键
		if order == 0 {
			order = 1
		} else {
			if order <= len(dict.SelectKeys) {
				code += string(dict.SelectKeys[order-1])
			} else {
				code += strconv.Itoa(order)
			}
		}
		ret = append(ret, Entry{wc[0], c, order})
	}
	return ret
}
