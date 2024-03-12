package server

import (
	"encoding/json"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/smq"
)

type Data struct {
	Source string `json:"source"`
	Path   string `json:"path"`
	Text   string `json:"text"`
	Merge  bool   `json:"merge"`
	Clean  bool   `json:"clean"`

	Dict []Dict `json:"dict"`
}

type Dict struct {
	Path   string `json:"path"`
	Format string `json:"format"`
	Push   int    `json:"push"`
	Keys   string `json:"keys"`
	Single bool   `json:"single"`
	Algo   string `json:"algo"`
	Space  string `json:"space"`
}

func (d *Data) Race() []byte {
	smq := &smq.Config{
		Merge: d.Merge,
		Clean: d.Clean,
	}
	if d.Merge {
		for _, text := range textList {
			t := &data.Text{
				Path: text,
			}
			smq.AddText(t)
		}
	} else {
		t := &data.Text{}
		switch d.Source {
		case "local":
			t.Path = d.Path
		case "clipboard":
			t.String = d.Text
			t.Name = "剪贴板"
		default:
			logger.Warn("不支持的数据源", "source", d.Source)
		}
		smq.AddText(t)
	}

	for _, v := range d.Dict {
		t := &data.Text{
			Path: v.Path,
		}
		d := &data.Dict{
			Text:       t,
			Format:     v.Format,
			Push:       v.Push,
			SelectKeys: v.Keys,
			Single:     v.Single,
			Algorithm:  v.Algo,
			SpacePref:  v.Space,
		}
		smq.AddDict(d)
	}
	res := smq.Race()
	if len(res) < 1 {
		return []byte{}
	}
	data, err := json.Marshal(res[0])
	if err != nil {
		logger.With("error", err).Error("json marshal")
		return []byte{}
	}
	return data
}
