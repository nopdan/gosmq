package dict

import (
	"github.com/nopdan/gosmq/pkg/matcher"
	"github.com/nopdan/gosmq/pkg/text"
)

type Dict struct {
	*text.Text
	// default, jisu, duoduo|bingling, jidian
	format string
	// 起顶码长
	push int
	// 选重键
	selectKeys []string
	// 是否只用码表里的单字
	Single bool
	// 匹配算法 greedy|ordered|dynamic
	algorithm string
	// 空格按键方式 both|left|right
	SpacePref string

	Matcher matcher.Matcher
	Length  int // 词条数
}

type DictOption func(*Dict)

func New(text *text.Text, opts ...DictOption) *Dict {
	dict := &Dict{
		Text:      text,
		format:    "default",
		algorithm: "greedy",
		SpacePref: "both",
	}
	for _, opt := range opts {
		opt(dict)
	}
	dict.init()
	return dict
}

func WithFormat(format string) DictOption {
	return func(opt *Dict) {
		opt.format = format
	}
}

func WithPush(push int) DictOption {
	return func(opt *Dict) {
		opt.push = push
	}
}

func WithSelectKeys(keys string) DictOption {
	return func(opt *Dict) {
		res := make([]string, 0, 10)
		for i := range len(keys) {
			res = append(res, string(keys[i]))
		}
		opt.selectKeys = res
	}
}

func WithSingle() DictOption {
	return func(opt *Dict) {
		opt.Single = true
	}
}

func WithAlgorithm(algorithm string) DictOption {
	return func(opt *Dict) {
		opt.algorithm = algorithm
	}
}

func WithSpacePref(spacePref string) DictOption {
	return func(opt *Dict) {
		opt.SpacePref = spacePref
	}
}
