package data

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nopdan/gosmq/pkg/util"
)

type Text struct {
	// 自定义名字
	Name string
	// 路径
	Path string
	// 上传的文件
	Bytes []byte
	// 剪切板文本
	String string

	reader *bufio.Reader
	size   int  // 文件大小
	isInit bool // 是否已经初始化
}

func (t *Text) Init() error {
	if t.isInit {
		return nil
	}
	t.isInit = true
	foo := func(rd io.Reader, size int) {
		t.reader = bufio.NewReaderSize(rd, 32*1024)
		t.size = size
	}
	if len(t.Path) > 0 {
		f, err := os.Open(t.Path)
		if err != nil {
			fmt.Printf("打开文件 %s 失败\n", t.Path)
			return err
		}
		fi, _ := f.Stat()
		if len(t.Name) == 0 {
			t.Name = fi.Name()
			t.Name = strings.TrimSuffix(t.Name, filepath.Ext(t.Name))
		}
		rd := util.ConvertReader(f)
		foo(rd, int(fi.Size()))
	} else if len(t.Bytes) > 0 {
		brd := bytes.NewReader(t.Bytes)
		rd := util.ConvertReader(brd)
		foo(rd, len(t.Bytes))
	} else if len(t.String) > 0 {
		rd := strings.NewReader(t.String)
		foo(rd, len(t.String))
	} else {
		return fmt.Errorf("无法初始化 Text")
	}
	if len(t.Name) == 0 {
		t.Name = "未命名"
	}
	return nil
}

// 重新初始化文本
func (t *Text) ReInit() {
	t.reader = nil
	t.size = 0
	t.isInit = false
	_ = t.Init()
}
