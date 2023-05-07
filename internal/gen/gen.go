package gen

import (
	"sort"
)

type Config struct {
	Path   string //待转换的码表路径
	Format string //待转换码表的格式

	SelectKeys    string //自定义选重键
	PushStart     int    //起顶码长
	SortByWordLen bool   // 按照词长重新排序
}

type Entry struct {
	Word string
	Code string
	Pos  int
}

func (c *Config) Gen() []*Entry {
	var dict []*Entry

	switch c.Format {
	case "jisu", "js":
		dict = c.LoadJisu()
	case "duoduo", "dd":
		dict = c.LoadTSV(true)
	case "bingling", "bl":
		dict = c.LoadTSV(false)
	default:
		panic("不支持的格式: " + c.Format)
	}

	if c.SortByWordLen {
		sort.SliceStable(dict, func(i, j int) bool {
			return len([]rune(dict[i].Word)) > len([]rune(dict[j].Word))
		})
	}
	return dict
}

// 专用，两位正数 1~99 byte 转 string
// func Itoa(b byte) string {
// 	if b < 10 {
// 		return string(b + '0')
// 	} else {
// 		return string([]byte{b / 10, b % 10})
// 	}
// }
