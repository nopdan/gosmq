package feeling

import (
	"bytes"
	"strconv"

	_ "embed"

	"github.com/nopdan/gosmq/pkg/util"
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
var Comb [128][128]*distrib

func init() {

	rd := bytes.NewReader(equivalent)
	tsv := util.NewTSV(rd)

	// 当量
	for {
		line, err := tsv.Read("\t")
		if err != nil {
			break
		}
		if len(line) != 2 {
			continue
		}
		code := line[0]
		dl, _ := strconv.ParseFloat(line[1], 64)
		if Comb[code[0]][code[1]] == nil {
			Comb[code[0]][code[1]] = new(distrib)
		}
		Comb[code[0]][code[1]].Equivalent = dl
	}

	rd.Reset(fingering)
	tsv = util.NewTSV(rd)
	// 小跨排
	line, _ := tsv.Read(" ")
	for _, v := range line {
		Comb[v[0]][v[1]].SingleSpan = true
	}
	// 大跨排
	line, _ = tsv.Read(" ")
	for _, v := range line {
		Comb[v[0]][v[1]].MultiSpan = true
	}
	// 错手
	line, _ = tsv.Read(" ")
	for _, v := range line {
		Comb[v[0]][v[1]].Staggered = true
	}
	// 小指干扰
	line, _ = tsv.Read(" ")
	for _, v := range line {
		Comb[v[0]][v[1]].Disturb = true
	}
}
