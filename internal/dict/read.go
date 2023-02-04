package dict

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/imetool/dtool/pkg/table"
)

func (dict *Dict) read() table.Table {
	ret := make(table.Table, 0, 1e5)
	scan := bufio.NewScanner(dict.Reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		pos := 1
		if len(wc) == 3 {
			pos, _ = strconv.Atoi(wc[2])
		} else if len(wc) != 2 {
			continue
		}
		if dict.Single {
			if len([]rune(wc[0])) != 1 {
				continue
			}
		}
		ret = append(ret, table.Entry{wc[0], wc[1], pos})
	}
	return ret
}

// 专用，两位正数 1~99 string 转 byte
// func Atoi(s string) byte {
// 	if len(s) == 1 {
// 		return s[0] - '0'
// 	} else if len(s) == 2 {
// 		a, b := s[0]-'0', s[1]-'0'
// 		return a*10 + b
// 	}
// 	return 1
// }
