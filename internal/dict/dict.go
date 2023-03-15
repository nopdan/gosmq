package dict

import (
	"fmt"
	"io"
	"strings"

	"github.com/imetool/dtool/pkg/table"
	"github.com/imetool/gosmq/pkg/matcher"
	"github.com/imetool/goutil/util"
)

type Dict struct {
	Single       bool   // 单字模式
	Algorithm    string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	PressSpaceBy string // 空格按键方式 left|right|both
	Verbose      bool   // 输出详细数据

	Name   string // 码表名
	Length int    // 词条数
	// 初始化 Matcher
	Matcher matcher.Matcher
	Reader  io.Reader // 赛码表 io 流
}

// 从文件加载码表
func (dict *Dict) Load(path string) {
	rd, err := util.Read(path)
	if err != nil {
		fmt.Println("Warning! 读取文件失败：", err)
		return
	}
	if dict.Name == "" {
		dict.Name = util.GetFileName(path)
	}
	dict.Reader = rd
	dict.initialize()
}

// 从字符串加载码表
func (dict *Dict) LoadString(text, name string) {
	if text == "" {
		fmt.Println("Warning! 码表输入为空。")
		return
	}
	dict.Name = name
	dict.Reader = strings.NewReader(text)
	dict.initialize()
}

// 初始化 Dict
func (dict *Dict) initialize() {
	// 读取
	t := dict.read()
	dict.Length = len(t)
	// 添加符号
	PUNCTS := GetPuncts()
	for k, v := range PUNCTS {
		t = append(t, table.Entry{string(k), v, 1})
	}

	// 匹配算法
	if dict.Single {
		dict.Algorithm = "single"
	}
	if dict.Matcher == nil {
		dict.Matcher = matcher.New(dict.Algorithm)
	}
	m := dict.Matcher

	// fmt.Printf("%+v\n", dict)
	m.Build(t)
}
