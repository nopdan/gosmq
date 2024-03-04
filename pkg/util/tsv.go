package util

import (
	"bufio"
	"io"
	"strings"
)

type tsv struct {
	*bufio.Scanner
}

func NewTSV(rd io.Reader) *tsv {
	return &tsv{
		Scanner: bufio.NewScanner(rd),
	}
}

// 读取一行
func (t *tsv) ReadLine() (string, error) {
	if ok := t.Scan(); !ok {
		if t.Err() == nil {
			return "", io.EOF
		}
		return "", t.Err()
	}
	return t.Text(), nil
}

// 读取文件的一行，按 sep 分隔返回切片
func (t *tsv) Read(sep string) ([]string, error) {
	if ok := t.Scan(); !ok {
		if t.Err() == nil {
			return nil, io.EOF
		}
		return nil, t.Err()
	}
	return strings.Split(t.Text(), sep), nil
}
