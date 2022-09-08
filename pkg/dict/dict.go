package dict

import (
	"fmt"
	"strings"

	"github.com/cxcn/gosmq/pkg/data"
	"github.com/cxcn/gosmq/pkg/matcher"
	"github.com/cxcn/gosmq/pkg/util"
)

// 从文件加载码表
func (dict *Dict) Load(path string) {
	rd, err := util.Read(path)
	if err != nil {
		fmt.Println("Warning! 读取文件失败：", path)
		return
	}
	if dict.Name == "" {
		dict.Name = util.GetFileName(path)
	}
	dict.Reader = rd
	dict.Legal = true
}

// 从字符串加载码表
func (dict *Dict) LoadFromString(s string) {
	dict.Reader = strings.NewReader(s)
	dict.Legal = true
}

// 初始化赛码表
func (dict *Dict) Init() {
	if dict.SelectKeys == "" {
		dict.SelectKeys = "_;'"
	}
	if dict.PushStart == 0 {
		dict.PushStart = 4
	}

	// 读取码表
	if dict.Transformer == nil {
		dict.Transformer = NewTransformer(dict.Format)
	}
	t := dict.Transformer.Read(dict)
	dict.Length = len(t)

	// 匹配算法
	if dict.Matcher == nil {
		dict.Matcher = matcher.New(dict.Algorithm)
	}
	m := dict.Matcher
	for i := range t {
		if dict.Single && len([]rune(t[i].Word)) > 1 {
			dict.Length--
			continue
		}
		m.Insert(t[i].Word, t[i].Code, t[i].Order)
	}
	// 添加符号
	PUNCTS := data.GetPuncts()
	for k, v := range PUNCTS {
		m.Insert(k, v, 1)
	}
	if dict.OutputDict && dict.Format != "default" && dict.Format != "" {
		outputDict(t, dict.Name)
	}
}
