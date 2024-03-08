package data

import (
	"bufio"
	"bytes"
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
	IsInit bool // 是否已经初始化
}

func (t *Text) Init() {
	if t.IsInit {
		logger.Debug("文本已经初始化过了", "name", t.Name)
		return
	}
	foo := func(rd io.Reader, size int) {
		t.reader = bufio.NewReaderSize(rd, 32*1024)
		t.size = size
	}
	if len(t.Path) > 0 {
		f, err := os.Open(t.Path)
		if err != nil {
			logger.Warn("文本初始化失败", "path", t.Path, "error", err)
			return
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
		logger.Warn("文本初始化失败", "name", t.Name)
		return
	}
	if len(t.Name) == 0 {
		t.Name = "未命名"
	}
	t.IsInit = true
	return
}

// 重新初始化文本
func (t *Text) ReInit() {
	t.reader = nil
	t.size = 0
	t.IsInit = false
	t.Init()
}
