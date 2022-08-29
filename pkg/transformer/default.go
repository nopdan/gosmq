package transformer

import (
	"bufio"
	"strconv"
	"strings"
)

type Smb struct{}

func (s Smb) Read(dict Dict) []Entry {
	ret := make([]Entry, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) < 3 {
			continue
		}
		order, _ := strconv.Atoi(wc[2])
		if order == 0 {
			order = 1
		}
		ret = append(ret, Entry{wc[0], wc[1], order})
	}
	return ret
}
