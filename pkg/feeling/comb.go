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

const (
	IsXKP = uint16(1) << iota
	IsDKP
	IsCS
	IsXZGR
)

// DL10<<8 + IsXKP|IsDKP|IsCS|IsXZGR
var Comb [128][128]uint16

func init() {

	// equivalent = strings.ReplaceAll(equivalent, "\r\n", "\n")
	// fingering = strings.ReplaceAll(fingering, "\r\n", "\n")

	// 当量
	for _, v := range strings.Split(equivalent, "\n") {
		tmp := strings.Split(v, "\t")
		if len(tmp) != 2 {
			continue
		}
		DL10, _ := strconv.Atoi(tmp[1])
		// 10 <=DL10 < 32 = 2^5
		Comb[tmp[0][0]][tmp[0][1]] = uint16(DL10) << 8
	}

	// 指法
	fg := strings.Split(fingering, "\n")
	// 小跨排
	xkp := strings.Split(fg[0], " ")
	for _, v := range xkp {
		Comb[v[0]][v[1]] |= IsXKP
	}
	// 大跨排
	dkp := strings.Split(fg[1], " ")
	for _, v := range dkp {
		Comb[v[0]][v[1]] |= IsDKP
	}
	// 错手
	cs := strings.Split(fg[2], " ")
	for _, v := range cs {
		Comb[v[0]][v[1]] |= IsCS
	}
	// 小指干扰
	xzgr := strings.Split(fg[3], " ")
	for _, v := range xzgr {
		Comb[v[0]][v[1]] |= IsXZGR
	}
}

func Debug() {
	fmt.Println(Comb)
}
