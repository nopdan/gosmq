package smq

import (
	"fmt"
	"testing"
	"time"

	"github.com/nopdan/gosmq/pkg/dict"
	"github.com/nopdan/gosmq/pkg/text"
)

func TestPuncts(t *testing.T) {
	for k, v := range zhKeysMap {
		fmt.Printf("%s\t%s\n", string(k), v)
	}
	fmt.Println(enKeysMap)
}

func TestSmq(t *testing.T) {
	now := time.Now()

	// s := New(WithSplit())
	s := New(WithSplit(), WithStat())

	// t1 := text.New("testdict", text.WithPath(`D:\Code\go\gosmq\build\dict\091点儿2023春.txt`))
	t1 := text.New("testdict", text.WithPath(`D:\Code\go\gosmq\build\dict\091五笔.txt`))
	d := dict.New(t1, dict.WithAlgorithm("ordered"))
	s.AddDict(d)
	fmt.Printf("载入码表耗时: %v\n", time.Since(now))

	// t1 = text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\心情决定事情.txt`))
	t1 = text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\《红楼梦》-曹雪芹.txt`))
	s.AddText(t1)

	res := s.Race()
	fmt.Println(res)

	fmt.Printf("耗时: %v\n", time.Since(now))
}

func TestSmqSingle(t *testing.T) {
	now := time.Now()

	// s := New(WithSplit())
	s := New(WithClean(), WithSplit(), WithStat())

	t1 := text.New("testdict", text.WithPath(`D:\Code\go\gosmq\build\dict\091点儿2023春.txt`))
	d := dict.New(t1, dict.WithSingle())
	s.AddDict(d)
	fmt.Printf("载入码表耗时: %v\n", time.Since(now))

	// t1 = text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\心情决定事情.txt`))
	t1 = text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\《红楼梦》-曹雪芹.txt`))
	s.AddText(t1)

	_ = s.Race()

	fmt.Printf("耗时: %v\n", time.Since(now))
}

func BenchmarkSmq(b *testing.B) {
	s := New()

	for i := 0; i < b.N; i++ {
		t1 := text.New("testdict", text.WithPath(`D:\Code\go\gosmq\build\dict\091点儿2023春.txt`))
		d := dict.New(t1, dict.WithSingle())
		s.AddDict(d)
	}

	// t1 := text.New("testdict", text.WithPath(`D:\Code\go\gosmq\build\dict\091点儿2023春.txt`))
	// d := dict.New(t1)
	// s.AddDict(d)

	t2 := text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\心情决定事情.txt`))
	s.AddText(t2)
	t3 := text.New("test", text.WithPath(`D:\Code\go\gosmq\build\text\《红楼梦》-曹雪芹.txt`))
	s.AddText(t3)

	res := s.Race()
	_ = res
}
