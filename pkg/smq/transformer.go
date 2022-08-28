package smq

import (
	"bytes"
	"log"
	"os"

	"github.com/cxcn/gosmq/pkg/transformer"
)

type Transformer interface {
	Read(transformer.Dict) []byte
}

// 转换赛码表
func (dict *Dict) transform() {
	if dict.Transformer == nil {
		switch dict.Format {
		case "jisu", "js":
			dict.Transformer = &transformer.Jisu{}
		case "duoduo", "dd":
			dict.Transformer = &transformer.Duoduo{}
		case "jidian", "jd":
			dict.Transformer = &transformer.Jidian{}
		case "bingling", "bl":
			dict.Transformer = &transformer.Duoduo{true}
		}
	}
	d := toTD(dict)
	// 输出赛码表
	if dict.Transformer != nil {
		newBytes := dict.Transformer.Read(d)
		err := os.WriteFile(dict.SavePath, newBytes, 0666)
		if err != nil {
			// SavePath 不对则保存在 dict 目录下
			os.Mkdir("dict", 0666)
			err = os.WriteFile("./dict/"+dict.Name+".txt", newBytes, 0666)
			if err != nil {
				log.Println(err)
			}
		}
		dict.reader = bytes.NewReader(newBytes)
	}
}

func toTD(dict *Dict) transformer.Dict {
	d := transformer.Dict{
		SavePath:   dict.SavePath,
		Name:       dict.Name,
		PushStart:  dict.PushStart,
		SelectKeys: dict.SelectKeys,
		Single:     dict.Single,
	}
	return d
}
