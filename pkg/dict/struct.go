package dict

import (
	"io"

	"github.com/cxcn/gosmq/pkg/matcher"
)

type Dict struct {
	Name   string // 码表名
	Single bool   // 单字模式

	Format string /* 码表格式
	default: 默认 本程序赛码表 词\t编码选重\t选重
	jisu:js 极速赛码表 词\t编码选重
	duoduo:dd 多多格式码表 词\t编码
	jidian:jd 极点格式 编码\t词1 词2 词3
	bingling:bl 冰凌格式码表 编码\t词
	*/
	Transformer Transformer // 自定义码表格式转换
	SelectKeys  string      // 普通码表自定义选重键(默认为_;')
	PushStart   int         // 普通码表起顶码长(码长大于等于此数，首选不会追加空格)

	// 初始化 Matcher
	Algorithm string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	Matcher   matcher.Matcher

	PressSpaceBy string // 空格按键方式 left|right|both
	OutputDict   bool   // 输出转换后的码表
	OutputDetail bool   // 输出详细数据

	Reader io.Reader // 赛码表 io 流
	Length int       // 词条数
	Legal  bool      // 合法输入
}

type Entry struct {
	Word  string
	Code  string
	Order int
}
