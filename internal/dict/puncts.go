package dict

import (
	_ "embed"
	"strings"
)

//go:embed assets/puncts.txt
var puncts string

func GetPuncts() map[string]string {
	ret := make(map[string]string)
	puncts = strings.ReplaceAll(puncts, "\r\n", "\n")

	lines := strings.Split(puncts, "\n")
	for _, line := range lines {
		wc := strings.Split(line, "\t")
		if len(wc) != 2 {
			continue
		}
		ret[wc[0]] = wc[1]
	}
	return ret
}
