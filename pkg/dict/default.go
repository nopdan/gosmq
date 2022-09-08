package dict

import (
	"bufio"
	"strconv"
	"strings"
)

type smb struct{}

func (s smb) Read(dict *Dict) []Entry {
	ret := make([]Entry, 0, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 3 {
			continue
		}
		order, _ := strconv.Atoi(wc[2])
		ret = append(ret, Entry{wc[0], wc[1], order})
	}
	return ret
}