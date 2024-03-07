package smq

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/result"
)

func TestPuncts(t *testing.T) {
	for k, v := range zhKeysMap {
		fmt.Printf("%s\t%s\n", string(k), v)
	}
	fmt.Println(enKeysMap)
}

func Print(res [][]*result.Result) {
	for _, v := range res {
		for _, vv := range v {
			data, _ := json.MarshalIndent(vv, "", "  ")
			fmt.Printf("result: %s\n", data)
		}
	}
}

func TestSmq(t *testing.T) {
	now := time.Now()
	s := &Config{}
	dict := &data.Dict{
		Text: &data.Text{
			Path: `D:\Code\go\gosmq\build\dict\091点儿2023春.txt`,
		},
	}
	s.AddDict(dict)
	fmt.Printf("载入码表耗时: %v\n", time.Since(now))
	text := &data.Text{Path: `D:\Code\go\gosmq\build\text\心情决定事情.txt`}
	s.AddText(text)
	res := s.Race()
	Print(res)
	fmt.Printf("耗时: %v\n", time.Since(now))

	now = time.Now()
	s.Reset()
	dict = &data.Dict{
		Text:   &data.Text{Path: `D:\Code\go\gosmq\build\dict\091点儿2023春.txt`},
		Single: true,
	}
	s.AddDict(dict)
	text.ReInit()
	s.AddText(text)
	res = s.Race()
	Print(res)
	fmt.Printf("耗时: %v\n", time.Since(now))
}

func BenchmarkSmq(b *testing.B) {
	s := &Config{}

	for i := 0; i < b.N; i++ {
		dict := &data.Dict{Text: &data.Text{Path: `D:\Code\go\gosmq\build\dict\091点儿2023春.txt`}}
		s.AddDict(dict)
	}
	text := &data.Text{Path: `D:\Code\go\gosmq\build\text\心情决定事情.txt`}
	s.AddText(text)

	text = &data.Text{Path: `D:\Code\go\gosmq\build\text\《红楼梦》-曹雪芹.txt`}
	s.AddText(text)

	res := s.Race()
	_ = res
}
