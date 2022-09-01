package transformer

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Jisu struct{}

func (j Jisu) Read(dict Dict) []Entry {
	ret := make([]Entry, 0, 1e5)
	scan := bufio.NewScanner(dict.Reader)

	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		if len(wc) != 2 {
			continue
		}
		c := wc[1]
		var code, match string
		var order int
		// a_ aa_
		if len(c)-1 > 0 && c[len(c)-1] == '_' {
			order = 1
			code = c
			if len(dict.SelectKeys) >= 1 && dict.SelectKeys[0] != '_' {
				code = c[:len(c)-1] + string(dict.SelectKeys[0])
			}
			ret = append(ret, Entry{wc[0], code, order})
			continue
		}

		re := regexp.MustCompile(`\d+$`)
		match = re.FindString(c)
		// akdb ksdw
		if match == "" {
			ret = append(ret, Entry{wc[0], c, 1})
			continue
		}

		// a1 aa3
		code = c[:len(c)-len(match)]
		order, err := strconv.Atoi(match)
		if err != nil {
			fmt.Println(match, err)
			order = 1
		}
		if order == 0 {
			order = 10
		}
		if len(dict.SelectKeys) >= order {
			code += string(dict.SelectKeys[order-1])
		}
		ret = append(ret, Entry{wc[0], code, order})
	}
	return ret
}
