package smq

import (
	"github.com/cxcn/gosmq/pkg/transformer"
)

type Transformer interface {
	Read(transformer.Dict) []transformer.Entry
}

// 转换赛码表
func (dict *Dict) transform() {
	if dict.Transformer == nil {
		switch dict.Format {
		case "jisu", "js":
			dict.Transformer = transformer.Jisu{}
		case "duoduo", "dd":
			dict.Transformer = transformer.Duoduo{}
		case "jidian", "jd":
			dict.Transformer = transformer.Jidian{}
		case "bingling", "bl":
			dict.Transformer = transformer.Duoduo{true}
		default:
			dict.Transformer = transformer.Duoduo{}
		}
	}
}

func toTD(dict *Dict) transformer.Dict {
	d := transformer.Dict{
		Name:       dict.Name,
		Reader:     dict.reader,
		PushStart:  dict.PushStart,
		SelectKeys: dict.SelectKeys,
	}
	return d
}
