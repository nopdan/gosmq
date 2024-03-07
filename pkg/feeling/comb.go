package feeling

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	_ "embed"
)

//go:embed assets/equivalent.txt
var equivalent []byte

//go:embed assets/fingering.txt
var fingering []byte

// 两键分布
type distrib struct {
	Equivalent float64 // 当量
	SingleSpan bool    // 小跨排
	MultiSpan  bool    // 大跨排
	Staggered  bool    // 错手
	Disturb    bool    // 小指干扰
}

// 1MB
var Combination [128][128]*distrib

func init() {

	rd := bytes.NewReader(equivalent)
	scan := bufio.NewScanner(rd)
	// 当量
	for scan.Scan() {
		line := strings.Split(scan.Text(), "\t")
		if len(line) != 2 {
			continue
		}
		code := line[0]
		dl, _ := strconv.ParseFloat(line[1], 64)
		if Combination[code[0]][code[1]] == nil {
			Combination[code[0]][code[1]] = new(distrib)
		}
		Combination[code[0]][code[1]].Equivalent = dl
	}

	rd.Reset(fingering)
	scan = bufio.NewScanner(rd)
	// 小跨排
	scan.Scan()
	line := strings.Split(scan.Text(), " ")
	for _, v := range line {
		Combination[v[0]][v[1]].SingleSpan = true
	}
	// 大跨排
	scan.Scan()
	line = strings.Split(scan.Text(), " ")
	for _, v := range line {
		Combination[v[0]][v[1]].MultiSpan = true
	}
	// 错手
	scan.Scan()
	line = strings.Split(scan.Text(), " ")
	for _, v := range line {
		Combination[v[0]][v[1]].Staggered = true
	}
	// 小指干扰
	scan.Scan()
	line = strings.Split(scan.Text(), " ")
	for _, v := range line {
		Combination[v[0]][v[1]].Disturb = true
	}
}
