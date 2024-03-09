package server

import (
	"encoding/json"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/smq"
)

type Data struct {
	Text  []Text `json:"text"`
	Dict  []Dict `json:"dict"`
	Clean bool   `json:"clean"`
}

type Dict struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Index  int    `json:"index"`
	Text   string `json:"text"`

	Format string `json:"format"`
	Push   int    `json:"push"`
	Keys   string `json:"keys"`
	Single bool   `json:"single"`
	Algo   string `json:"algo"`
	Space  string `json:"space"`
}

type Text struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Index  int    `json:"index"`
	Text   string `json:"text"`
}

func (d *Data) Race() []byte {
	smq := &smq.Config{
		Clean: d.Clean,
	}
	for _, v := range d.Text {
		t := &data.Text{
			Name: v.Name,
		}
		switch v.Source {
		case "local":
			t.Path = v.Path
		case "upload":
			t.Bytes = files[v.Index]
		case "clipboard":
			t.String = v.Text
		default:
			logger.Warn("不支持的数据源", "source", v.Source)
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
		case "upload":
			t.Bytes = files[v.Index]
		case "clipboard":
			t.String = v.Text
		default:
			logger.Warn("不支持的数据源", "source", v.Source)
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
	data, err := json.Marshal(res)
	if err != nil {
		logger.With("error", err).Error("json marshal")
		return []byte{}
	}
	return data
}
