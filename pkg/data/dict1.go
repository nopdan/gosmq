package data

import (
	"bufio"
	"strconv"
	"strings"
)

// 读取每行只有一个词条的码表
func (d *Dict) read1(format string) []*Entry {
	ret := make([]*Entry, 0, 1e5)
	// 统计编码出现的次数
	stat := make(map[string]int)
	// 逐行读取码表
	scan := bufio.NewScanner(d.Text.reader)
	var word, code string
	var pos int
	for scan.Scan() {
		values := strings.Split(scan.Text(), "\t")
		if len(values) < 2 {
			continue
		}
		switch format {
		case "duoduo":
			word, code = values[0], values[1]
		case "bingling":
			word, code = values[1], values[0]
		case "chai":
			if len(values) < 5 {
				continue
			}
			word, code = values[0], values[3]
			pos, _ = strconv.Atoi(values[4])
			code = d.addSuffix(code, pos)
		case "jisu":
			word, code = values[0], values[1]
			// 带空格 a_ aa_
			if len(code) > 1 && code[len(code)-1] == '_' {
				pos = 1
				break
			}
			code, pos = findSuffixInteger(code)
			// 不带数字 akdb ksdw
			if pos == 1 {
				break
			}
			// 数字选重 a1 aa3
			code = d.addSuffix(code, pos)
		}
		switch format {
		case "duoduo", "bingling":
			stat[code]++
			pos = stat[code]
			code = d.addSuffix(code, pos)
		}
		ret = append(ret, &Entry{word, code, pos})
		d.insert(word, code, pos)
	}
	return ret
}

// 为编码加上选重键，pos 是编码出现的次数，最小为 1
func (d *Dict) addSuffix(code string, pos int) string {
	// 选重
	if pos > 1 {
		// 缓存 选重数字转为字符串
		for pos > len(d.selectKeys)-2 {
			d.selectKeys = append(d.selectKeys, strconv.Itoa(len(d.selectKeys)+1))
		}
		return code + d.selectKeys[pos-1]
	}
	// 首选
	if d.pattern.MatchString(code) {
		return code
	}
	return code + "_"
}

// 查找末尾数字，返回前缀和后缀
func findSuffixInteger(s string) (string, int) {
	var prefix, suffix string
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			suffix = string(s[i]) + suffix
		} else {
			prefix = s[:i+1]
			break
		}
	}
	if suffix == "" {
		return s, 1
	}

	i, _ := strconv.Atoi(suffix)
	if i <= 0 {
		i = 10
	}
	return prefix, i
}
