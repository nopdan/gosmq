package gen

import (
	"bufio"
	"strconv"
	"strings"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/rose"
)

func (c *Config) ReadJisu() rose.WordLibrary {
	ret := make(rose.WordLibrary, 0, 1e5)
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
		code := wc[1]
		// 带空格 a_ aa_
		if len(code)-1 > 0 && code[len(code)-1] == '_' {
			ret = append(ret, &rose.WubiEntry{wc[0], code, 1})
			continue
		}

		code, suf := FindSuffixInteger(code)
		// 不带数字 akdb ksdw
		if suf == "" {
			ret = append(ret, &rose.WubiEntry{wc[0], code, 1})
			continue
		}

		// 数字选重 a1 aa3
		pos, _ := strconv.Atoi(suf)
		if pos <= 0 {
			pos = 10
		}
		if len(c.SelectKeys) >= pos {
			code += string(c.SelectKeys[pos-1])
		}
		// fmt.Println(wc[0], code, pos)
		ret = append(ret, &rose.WubiEntry{wc[0], code, pos})
	}
	return ret
}

// 查找末尾数字，返回前缀和后缀
func FindSuffixInteger(s string) (string, string) {
	var preffix, suffix string
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			suffix = string(s[i]) + suffix
		} else {
			preffix = s[:i+1]
			return preffix, suffix
		}
	}
	// 全是数字
	return s, ""
}
