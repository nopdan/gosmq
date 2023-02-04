package feeling

import (
	"fmt"
	"strconv"
	"strings"

	_ "embed"
)

//go:embed assets/equivalent.txt
var equivalent string

//go:embed assets/fingering.txt
var fingering string

type combData struct {
	DL10   int
	DL     float64
	IsXKP  bool
	IsDKP  bool
	IsCS   bool
	IsXZGR bool
}

var Comb map[string]*combData

func init() {

	// equivalent = strings.ReplaceAll(equivalent, "\r\n", "\n")
	// fingering = strings.ReplaceAll(fingering, "\r\n", "\n")

	Comb = make(map[string]*combData)

	// 当量
	for _, v := range strings.Split(equivalent, "\n") {
		tmp := strings.Split(v, "\t")
		if len(tmp) != 2 {
			continue
		}
		c := new(combData)
		c.DL10, _ = strconv.Atoi(tmp[1])
		c.DL = float64(c.DL10) / 10
		Comb[tmp[0]] = c
	}

	// 指法
	fg := strings.Split(fingering, "\n")
	// 小跨排
	xkp := strings.Split(fg[0], " ")
	for _, v := range xkp {
		Comb[v].IsXKP = true
	}
	// 大跨排
	dkp := strings.Split(fg[1], " ")
	for _, v := range dkp {
		Comb[v].IsDKP = true
	}
	// 错手
	cs := strings.Split(fg[2], " ")
	for _, v := range cs {
		Comb[v].IsCS = true
	}
	// 小指干扰
	xzgr := strings.Split(fg[3], " ")
	for _, v := range xzgr {
		Comb[v].IsXZGR = true
	}
}

func Debug() {
	fmt.Println(Comb)
}
