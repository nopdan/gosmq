package transformer

import (
	"bufio"
	"strings"
)

type Duoduo struct {
	Reverse bool
}

func (d Duoduo) Read(dict Dict) []Entry {
	ret := make([]Entry, 0, 1e5)
	mapOrder := make(map[string]int)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 2 {
			continue
		}
		var w, c string
		if d.Reverse {
			w, c = wc[1], wc[0]
		} else {
			w, c = wc[0], wc[1]
		}

		mapOrder[c]++
		order := mapOrder[c]
		ret = append(ret, Entry{w, c, order})
	}
	return ret
}
