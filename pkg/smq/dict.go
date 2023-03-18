package smq

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/imetool/gosmq/pkg/matcher"
	"github.com/imetool/goutil/util"
)

type Dict struct {
	Single       bool   // 单字模式
	Algorithm    string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	PressSpaceBy string // 空格按键方式 left|right|both
	Json         bool   // 输出 json 详细数据
	Stat         bool   // 输出词条数据
	Split        bool   // 输出分词数据

	Name   string // 码表名
	Length int    // 词条数
	Clean  bool   // 只统计词库中的词条

	Matcher matcher.Matcher // 初始化 Matcher
	Reader  io.Reader       // 赛码表 io 流
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
	// 匹配算法
	if dict.Single {
		dict.Algorithm = "single"
	}
	if dict.Matcher == nil {
		dict.Matcher = matcher.New(dict.Algorithm)
	}
	m := dict.Matcher

	// 读取码表，构建 matcher
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
		dict.Length++
		m.Insert(wc[0], wc[1], pos)
	}
	m.Build()
}
