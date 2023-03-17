package gen

import (
	"sort"
	"strconv"

	"github.com/imetool/dtool/pkg/table"
)

type Config struct {
	Path   string //待转换的码表路径
	Format string //待转换码表的格式

	SelectKeys    string //自定义选重键
	PushStart     int    //起顶码长
	SortByWordLen bool   // 按照词长重新排序
}

func (c Config) Gen() table.Table {
	var d table.Table
	// 极速赛码表格式
	if c.Format == "jisu" {
		d = c.ReadJisu()
		return d
	}

	d = table.Parse(c.Format, c.Path)
	// 极点格式不需要处理候选位置
	if c.Format != "jidian" {
		d.GenPos()
	}

	for i := range d {
		if d[i].Pos <= 0 {
			d[i].Pos = 1
		}
		d[i].Code = c.addSuffix(d[i].Code, d[i].Pos)
	}
	if c.SortByWordLen {
		sort.SliceStable(d, func(i, j int) bool {
			return len([]rune(d[i].Word)) > len([]rune(d[j].Word))
		})
	}

	return d
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
