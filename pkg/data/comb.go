package data

import (
	_ "embed"
	"strconv"
	"strings"
)

type Comb struct {
	Dl   map[string]int  // 当量*10
	Xkp  map[string]bool // 小跨排
	Dkp  map[string]bool // 大跨排
	Cs   map[string]bool // 错手
	Xzgr map[string]bool // 小指干扰
}

//go:embed assets/equivalent.txt
var equivalent string

//go:embed assets/fingering.txt
var fingering string

func GetComb() Comb {

	var ret Comb
	ret.Dl = make(map[string]int)
	ret.Xkp = make(map[string]bool)
	ret.Dkp = make(map[string]bool)
	ret.Cs = make(map[string]bool)
	ret.Xzgr = make(map[string]bool)

	// equivalent = strings.ReplaceAll(equivalent, "\r\n", "\n")
	// fingering = strings.ReplaceAll(fingering, "\r\n", "\n")

	// 当量
	for _, v := range strings.Split(equivalent, "\n") {
		tmp := strings.Split(v, "\t")
		if len(tmp) != 2 {
			continue
		}
		dl, _ := strconv.Atoi(tmp[1])
		ret.Dl[tmp[0]] = dl
	}

	// 指法
	fg := strings.Split(fingering, "\n")
	// 小跨排
	xkp := strings.Split(fg[0], " ")
	for _, v := range xkp {
		ret.Xkp[v] = true
	}
	// 大跨排
	dkp := strings.Split(fg[1], " ")
	for _, v := range dkp {
		ret.Dkp[v] = true
	}
	// 错手
	cs := strings.Split(fg[2], " ")
	for _, v := range cs {
		ret.Cs[v] = true
	}
	// 小指干扰
	xzgr := strings.Split(fg[3], " ")
	for _, v := range xzgr {
		ret.Xzgr[v] = true
	}

	return ret
}
