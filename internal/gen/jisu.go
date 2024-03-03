package gen

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/nopdan/gosmq/pkg/util"
)

func (c *Config) LoadJisu() []*Entry {
	ret := make([]*Entry, 0, 1e5)
	rd, err := util.Read(c.Path)
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 2 {
			continue
		}
		word, code := wc[0], wc[1]

		// 带空格 a_ aa_
		if len(code) > 1 && code[len(code)-1] == '_' {
			ret = append(ret, &Entry{word, code, 1})
			continue
		}

		pre, suf := FindSuffixInteger(code)
		// 不带数字 akdb ksdw
		if suf == "" {
			ret = append(ret, &Entry{word, code, 1})
			continue
		}

		// 数字选重 a1 aa3
		pos, _ := strconv.Atoi(suf)
		if pos <= 0 {
			pos = 10
		}
		// 添加自定义选重键
		if pos <= len(c.SelectKeys) {
			code = pre + string(c.SelectKeys[pos-1])
		}
		// fmt.Println(wc[0], code, pos)
		ret = append(ret, &Entry{wc[0], code, pos})
	}
	return ret
}

// 查找末尾数字，返回前缀和后缀
func FindSuffixInteger(s string) (string, string) {
	var prefix, suffix string
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			suffix = string(s[i]) + suffix
		} else {
			prefix = s[:i+1]
			return prefix, suffix
		}
	}
	// 全是数字
	return s, ""
}
