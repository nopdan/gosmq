package data

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/nopdan/gosmq/pkg/matcher"
	"github.com/nopdan/gosmq/pkg/util"
)

var logger = util.Logger

type Dict struct {
	// 码表文件
	Text *Text
	// default, jisu, duoduo|bingling, jidian
	Format string
	// 顶屏模式
	Pattern string
	pattern *regexp.Regexp
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
	IsInit  bool // 是否已经初始化
}

// 初始化 Dict
func (d *Dict) Init() {
	if d.IsInit {
		logger.Debug("码表已经初始化过了", "name", d.Text.Name)
		return
	}
	if d.Text == nil {
		logger.Warn("码表未指定", "name", d.Text.Name)
		return
	}
	if !d.Text.IsInit {
		d.Text.Init()
		if !d.Text.IsInit {
			logger.Warn("码表初始化失败", "name", d.Text.Name)
			return
		}
	}
	now := time.Now()
	// 顶屏模式
	var err error
	d.pattern, err = regexp.Compile(d.Pattern)
	if err != nil {
		d.pattern, _ = regexp.Compile("^.{4,}$")
	}
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
			logger.Warn("还未实现此算法")
			fallthrough
		default:
			d.Matcher = matcher.NewTrie(false)
		}
	}
	// 读取码表，构建 matcher
	var dict []*Entry
	switch d.Format {
	case "default", "":
		d.load()
	case "duoduo", "dd", "rime":
		dict = d.read1("duoduo")
	case "bingling", "bl":
		dict = d.read1("bingling")
	case "chai":
		dict = d.read1("chai")
	case "jisu", "js":
		dict = d.read1("jisu")
	case "xiaoxiao", "xx", "jidian", "jd":
		dict = d.read2()
	default:
		logger.Fatal("码表格式不正确", "format", d.Format)
	}
	if d.Length == 0 {
		logger.Warn("码表为空", "name", d.Text.Name, "path", d.Text.Path)
		return
	}

	d.IsInit = true
	if dict == nil || len(dict) == 0 {
		d.Matcher.Build()
		logger.Info(fmt.Sprintf("已载入码表: %s", d.Text.Name), "耗时", time.Since(now))
		return
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
	return
}
