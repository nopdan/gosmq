package util

import (
	"path/filepath"
	"strings"

	"golang.org/x/exp/constraints"
)

func Div[T constraints.Integer | constraints.Float](x, y T) float64 {
	return float64(x) / float64(y)
}

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

// 获取文件名，不含路径和后缀
func GetFileName(fp string) string {
	name := filepath.Base(fp)
	return strings.TrimSuffix(name, filepath.Ext(name))
}
