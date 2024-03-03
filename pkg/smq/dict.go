package smq

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nopdan/gosmq/pkg/matcher"
	"github.com/nopdan/gosmq/pkg/util"
)

type Dict struct {
	Name    string // 码表名
	Single  bool   // 单字模式
	Stable  bool   // 按照码表顺序
	UseTail bool   // 压缩 tail
	Clean   bool   // 只统计词库中的词条

	// 空格按键方式 left|right|both
	PressSpaceBy string

	Split bool // 统计分词数据并输出
	Stat  bool // 统计词条数据并输出

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
		dict.matcher = matcher.NewSingle()
	}
	if dict.matcher == nil {
		dict.matcher = matcher.NewTrie(dict.Stable, dict.UseTail)
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
