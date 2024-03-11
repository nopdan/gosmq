package data

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

// 逻辑 CPU 数量（线程数）
var NUM_CPU = runtime.NumCPU()

type Text struct {
	// 自定义名字
	Name string
	// 路径
	Path string
	// 上传的文件
	Bytes []byte
	// 剪切板文本
	String string

	reader  *bufio.Reader
	size    int  // 文件大小
	IsInit  bool // 是否已经初始化
	bufSize int
}

func (t *Text) Init() {
	if t.IsInit {
		logger.Debug("文本已经初始化过了", "name", t.Name)
		return
	}
	if len(t.Path) > 0 {
		t.loadFile()
	} else if len(t.Bytes) > 0 {
		t.loadBytes()
	} else if len(t.String) > 0 {
		t.determineBufSize(len(t.String))
		t.reader = bufio.NewReader(strings.NewReader(t.String))
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

func (t *Text) determineBufSize(size int) {
	t.size = size
	if size > NUM_CPU*256*1024 {
		t.bufSize = 256 * 1024
	} else if size > NUM_CPU*64*1024 {
		t.bufSize = 64 * 1024
	} else if size > 16*1024 {
		t.bufSize = 16 * 1024
	} else {
		t.bufSize = 4 * 1024 // defaultBufSize
	}
}

func (t *Text) loadFile() {
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
	t.determineBufSize(int(fi.Size()))

	buf := make([]byte, 1024)
	_, _ = f.Read(buf)
	f.Seek(0, io.SeekStart)
	// 检测编码格式
	cs := t.detect(buf, f)
	if cs == nil {
		return
	}
	// 转换
	brd := bufio.NewReaderSize(f, t.bufSize)
	rd, _ := charset.NewReaderLabel(cs.Charset, brd)
	t.reader = bufio.NewReader(rd)
}

func (t *Text) loadBytes() {
	t.determineBufSize(len(t.Bytes))

	brd := bytes.NewReader(t.Bytes)
	buf := make([]byte, 1024)
	_, _ = brd.Read(buf)
	brd.Seek(0, io.SeekStart)
	// 检测编码格式
	cs := t.detect(buf, brd)
	if cs == nil {
		return
	}
	// 转换
	rd, _ := charset.NewReaderLabel(cs.Charset, brd)
	t.reader = bufio.NewReader(rd)
}

// 检测编码格式，返回 nil 不需要再处理
func (t *Text) detect(buf []byte, input io.Reader) *chardet.Result {
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf)
	if err != nil {
		t.reader = bufio.NewReaderSize(input, t.bufSize)
		return nil
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	// 删除 BOM 文件头
	boms := make(map[string][]byte)
	boms["UTF-16BE"] = []byte{0xfe, 0xff}
	boms["UTF-16LE"] = []byte{0xff, 0xfe}
	boms["UTF-8"] = []byte{0xef, 0xbb, 0xbf}
	if b, ok := boms[cs.Charset]; ok {
		if bytes.HasPrefix(buf, b) {
			_, _ = input.Read(b)
		}
	}
	if cs.Charset == "UTF-8" {
		t.reader = bufio.NewReaderSize(input, t.bufSize)
		return nil
	}
	return cs
}

// 重新初始化文本
func (t *Text) ReInit() {
	t.reader = nil
	t.size = 0
	t.IsInit = false
	t.Init()
}
