package gen

import (
	"sort"
	"strconv"

	"github.com/flowerime/rose/pkg/rose"
)

type Config struct {
	Path   string //待转换的码表路径
	Format string //待转换码表的格式

	SelectKeys    string //自定义选重键
	PushStart     int    //起顶码长
	SortByWordLen bool   // 按照词长重新排序
}

func (c *Config) Gen() rose.WubiTable {
	var wl rose.WordLibrary
	var ct rose.CodeTable

	// 极速赛码表格式
	if c.Format == "jisu" {
		wl = c.ReadJisu()
		ct = wl.ToCodeTable()
	} else {
		d := rose.Parse(c.Path, c.Format)
		ct = d.ToCodeTable()
	}
	wt := ct.ToWubiTable()

	for i := range wt {
		if wt[i].Pos <= 0 {
			wt[i].Pos = 1
		}
		wt[i].Code = c.addSuffix(wt[i].Code, wt[i].Pos)
	}
	if c.SortByWordLen {
		sort.SliceStable(wt, func(i, j int) bool {
			return len([]rune(wt[i].Word)) > len([]rune(wt[j].Word))
		})
	}
	return wt
}

// 加上选重键
func (c Config) addSuffix(s string, pos int) string {
	if pos != 1 || len(s) < c.PushStart {
		if int(pos) <= len(c.SelectKeys) {
			s += string(c.SelectKeys[pos-1])
		} else {
			s += strconv.Itoa(pos)
		}
	}
	return s
}

// 专用，两位正数 1~99 byte 转 string
// func Itoa(b byte) string {
// 	if b < 10 {
// 		return string(b + '0')
// 	} else {
// 		return string([]byte{b / 10, b % 10})
// 	}
// }
