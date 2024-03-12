package server

import (
	"encoding/json"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/smq"
)

type Data struct {
	Merge bool `json:"merge"`
	Clean bool `json:"clean"`

	Text Text   `json:"text"`
	Dict []Dict `json:"dict"`
}

type Dict struct {
	Source string `json:"source"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	String string `json:"string"`

	Format string `json:"format"`
	Push   int    `json:"push"`
	Keys   string `json:"keys"`
	Single bool   `json:"single"`
	Algo   string `json:"algo"`
	Space  string `json:"space"`
}

type Text struct {
	Source string `json:"source"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	String string `json:"string"`
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
		t := &data.Text{
			Name: d.Text.Name,
		}
		switch d.Text.Source {
		case "local":
			t.Path = d.Text.Path
		case "clipboard":
			t.String = d.Text.String
		default:
			logger.Warn("不支持的数据源", "source", d.Text.Source)
		}
		smq.AddText(t)
	}

	for _, v := range d.Dict {
		t := &data.Text{
			Name: v.Name,
		}
		switch v.Source {
		case "local":
			t.Path = v.Path
		case "clipboard":
			t.String = v.String
		default:
			logger.Warn("不支持的数据源", "source", v.Source)
			continue
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
