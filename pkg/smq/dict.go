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
	Name         string // 码表名
	Single       bool   // 单字模式
	Algorithm    string // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
	PressSpaceBy string // 空格按键方式 left|right|both
	Clean        bool   // 只统计词库中的词条
	Verbose      bool   // 详细

	matcher matcher.Matcher // 初始化 Matcher
	reader  io.Reader       // 赛码表 io 流
	length  int             // 词条数
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
	dict.reader = rd
	dict.init()
}

// 从字符串加载码表
func (dict *Dict) LoadString(text, name string) {
	if text == "" {
		fmt.Println("Warning! 码表输入为空。")
		return
	}
	dict.Name = name
	dict.reader = strings.NewReader(text)
	dict.init()
}

// 初始化 Dict
func (dict *Dict) init() {
	// 匹配算法
	if dict.Single {
		dict.Algorithm = "single"
	}
	if dict.matcher == nil {
		dict.matcher = matcher.New(dict.Algorithm)
	}
	m := dict.matcher

	// 读取码表，构建 matcher
	scan := bufio.NewScanner(dict.reader)
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
		dict.length++
		m.Insert(wc[0], wc[1], pos)
	}
	m.Build()
}
