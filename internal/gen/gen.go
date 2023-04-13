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

func (c *Config) Gen() rose.Table {
	var t rose.Table
	// 极速赛码表格式
	if c.Format == "jisu" {
		t = c.ReadJisu()
		return t
	}

	d := rose.Parse(c.Path, c.Format)
	t = d.GetTable()
	// 极点格式不需要处理候选位置
	if c.Format != "jidian" {
		d.GenPos()
	}

	for i := range t {
		if t[i].Pos <= 0 {
			t[i].Pos = 1
		}
		t[i].Code = c.addSuffix(t[i].Code, t[i].Pos)
	}
	if c.SortByWordLen {
		sort.SliceStable(t, func(i, j int) bool {
			return len([]rune(t[i].Word)) > len([]rune(t[j].Word))
		})
	}

	return t
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
