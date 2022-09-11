package dict

import (
	"bufio"
	"strings"
)

type jidian struct{}

func (jidian) Read(dict *Dict) []Entry {
	ret := make([]Entry, 0, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), " ")
		if len(wc) < 2 {
			continue
		}
		code := wc[0]
		for i := 1; i < len(wc); i++ {
			code = dict.getRealCode(code, i)
			ret = append(ret, Entry{wc[i], code, i})
		}
	}
	return ret
}
