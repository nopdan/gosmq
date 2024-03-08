package util

import (
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// 切片指定索引位置加1，若索引超出范围则扩容
func Increase(sli *[]int, idx int) {
	AddTo(1, sli, idx)
}

// 添加到切片指定索引位置，若索引超出范围则扩容
func AddTo(val int, sli *[]int, idx int) {
	for idx > len(*sli)-1 {
		*sli = append(*sli, 0)
	}
	(*sli)[idx] += val
}

// unsafe 强制转换 []byte 为 string
func UnsafeToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafe 强制转换 string 为 []byte
func UnsafeToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
			Len int
		}{s, len(s), len(s)},
	))
}

// 遍历文件夹
func WalkDir(dir string) []string {
	var files []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// 遍历文件夹 + 指定后缀
func WalkDirWithSuffix(dir string, suffix string) []string {
	var files []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files = append(files, path)
		}
		return nil
	})
	return files
}
