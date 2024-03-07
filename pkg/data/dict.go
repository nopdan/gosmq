package data

import (
	"fmt"
	"sync"

	"github.com/nopdan/gosmq/pkg/matcher"
)

type Dict struct {
	// 码表文件
	Text *Text
	// default, jisu, duoduo|bingling, jidian
	Format string
	// 起顶码长
	Push int
	// 选重键
	SelectKeys string
	selectKeys []string
	// 是否只用码表里的单字
	Single bool
	// 匹配算法 greedy|ordered|dynamic
	Algorithm string
	// 空格按键方式 both|left|right
	SpacePref string
	// 转换码表是否覆盖
	Overwrite bool

	Matcher matcher.Matcher
	Length  int  // 词条数
	isInit  bool // 是否已经初始化
}

// 初始化 Dict
func (d *Dict) Init() error {
	if d.isInit {
		return nil
	}
	if d.Text == nil {
		return fmt.Errorf("无法初始化 Dict，Text 为空")
	}
	if !d.Text.isInit {
		err := d.Text.Init()
		if err != nil {
			return err
		}
	}
	d.isInit = true
	// 选重键
	d.selectKeys = make([]string, 0, 10)
	for i := range len(d.SelectKeys) {
		d.selectKeys = append(d.selectKeys, string(d.SelectKeys[i]))
	}
	// 空格偏好
	if d.SpacePref == "" {
		d.SpacePref = "both"
	}
	// 匹配算法
	if d.Single {
		d.Matcher = matcher.NewSingle()
	} else {
		switch d.Algorithm {
		case "greedy":
			d.Matcher = matcher.NewTrie(false)
		case "ordered":
			d.Matcher = matcher.NewTrie(true)
		case "dynamic":
			// TODO
			panic("未实现")
		default:
			d.Matcher = matcher.NewTrie(false)
		}
	}
	// 读取码表，构建 matcher
	var dict []*Entry
	switch d.Format {
	case "default", "":
		d.load()
	case "jisu", "js":
		dict = d.loadJisu()
	case "duoduo", "dd", "rime":
		dict = d.loadTSV(true)
	case "bingling", "bl":
		dict = d.loadTSV(false)
	case "xiaoxiao", "xx", "jidian", "jd":
		dict = d.loadXiao()
	default:
		panic("不支持的格式: " + d.Format)
	}

	if dict == nil || len(dict) == 0 {
		d.Matcher.Build()
		return nil
	}

	// 输出转换后的赛码表
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.output(dict)
	}()
	d.Matcher.Build()
	wg.Wait()
	return nil
}
