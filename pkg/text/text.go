package text

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nopdan/gosmq/pkg/util"
)

type Text struct {
	// 文件名
	Name string
	// 来源 local|upload|clipboard
	source string

	// 绝对路径
	path string
	// 上传的文件
	data []byte
	// 剪切板文本
	plainText string

	Reader *bufio.Reader
	Size   int // 文件大小
}

type TextOption func(*Text)

func New(name string, opts ...TextOption) *Text {
	text := &Text{
		Name:   name,
		source: "local",
	}
	for _, opt := range opts {
		opt(text)
	}
	err := text.init()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return text
}

func WithSource(source string) TextOption {
	return func(opt *Text) {
		opt.source = source
	}
}

func WithPath(path string) TextOption {
	return func(opt *Text) {
		opt.path = path
	}
}

func WithData(data []byte) TextOption {
	return func(opt *Text) {
		opt.data = data
	}
}

func WithText(text string) TextOption {
	return func(opt *Text) {
		opt.plainText = text
	}
}

func (t *Text) init() error {
	foo := func(rd io.Reader, size int) {
		t.Reader = bufio.NewReaderSize(rd, 32*1024)
		t.Size = size
	}
	switch t.source {
	case "local":
		// 读取本地文件
		f, err := os.Open(t.path)
		if err != nil {
			return fmt.Errorf("text.New(): %w", err)
		}
		fi, _ := f.Stat()
		rd := util.ConvertReader(f)
		foo(rd, int(fi.Size()))
	case "upload":
		if len(t.data) == 0 {
			return fmt.Errorf("text.New(): data is empty")
		}
		brd := bytes.NewReader(t.data)
		rd := util.ConvertReader(brd)
		foo(rd, len(t.data))
	case "clipboard":
		if len(t.plainText) == 0 {
			return fmt.Errorf("text.New(): plainText is empty")
		}
		rd := strings.NewReader(t.plainText)
		foo(rd, len(t.plainText))
	}
	return nil
}
