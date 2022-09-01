package smq

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
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
	SavePath    string      // 读取非默认码表格式时自动转换并保存的路径，默认保存在 dict 目录下
	SelectKeys  string      // 普通码表自定义选重键(默认为_;')
	PushStart   int         // 普通码表起顶码长(码长大于等于此数，首选不会追加空格)

	// 初始化 Matcher
	Algorithm string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	Matcher   Matcher

	PressSpaceBy string // 空格按键方式 left|right|both
	OutputDict   bool   // 输出转换后的码表
	OutputDetail bool   // 输出详细数据

	reader io.Reader // 赛码表 io 流
	length int       // 词条数
	legal  bool      // 合法输入
	trans  bool      // 是否转码表格式
}

// 从 io 流加载码表
func (dict *Dict) Load(rd io.Reader) {
	dict.reader = Tranformer(rd)
	dict.legal = true
}

// 从字符串流加载码表
func (dict *Dict) LoadFromString(s string) {
	dict.reader = readFromString(s)
	dict.legal = true
}

// 从文件加载码表
func (dict *Dict) LoadFromPath(path string) {
	rd, err := readFromPath(path)
	if err != nil {
		log.Println("Warning! 从文件读取码表失败，路径：", path)
		return
	}
	if dict.Name == "" {
		dict.Name = GetFileName(path)
	}
	dict.reader = rd
	dict.legal = true
}

func (dict *Dict) init() {
	// 读取码表
	if dict.SelectKeys == "" {
		dict.SelectKeys = "_;'"
	}
	if dict.PushStart == 0 {
		dict.PushStart = 4
	}
	dict.read()
}

func (dict *Dict) read() {
	dict.match()
	m := dict.Matcher

	dict.transform()
	d := toTD(dict)
	t := dict.Transformer.Read(d)
	dict.length = len(t)

	var buf bytes.Buffer
	buf.Grow(1e5)
	for i := 0; i < len(t); i++ {
		if dict.Single && len([]rune(t[i].Word)) > 1 {
			dict.length--
			continue
		}
		code := dict.getRealCode(t[i].Code, t[i].Order)
		m.Insert(t[i].Word, code, t[i].Order)
	}
	// 添加符号
	for k, v := range PUNCTS {
		m.Insert(k, v, 1)
	}
	m.Handle()

	// 输出赛码表
	if !dict.trans && dict.OutputDict {
		for i := 0; i < len(t); i++ {
			buf.WriteString(t[i].Word)
			buf.WriteByte('\t')
			buf.WriteString(t[i].Code)
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(t[i].Order))
			buf.WriteByte('\n')
		}
		err := os.WriteFile(dict.SavePath, buf.Bytes(), 0666)
		if err != nil {
			// SavePath 不对则保存在 dict 目录下
			os.Mkdir("dict", 0666)
			err = os.WriteFile("./dict/"+dict.Name+".txt", buf.Bytes(), 0666)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// 加上选重键
func (dict *Dict) getRealCode(c string, order int) string {
	if order != 1 || len(c) < dict.PushStart {
		if order <= len(dict.SelectKeys) {
			c += string(dict.SelectKeys[order-1])
		} else {
			c += strconv.Itoa(order)
		}
	}
	return c
}
