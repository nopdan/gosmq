package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/imetool/gosmq/internal/serve"
	"github.com/imetool/gosmq/pkg/smq"
)

// 保存 html 结果
func OutputHTML(res []*smq.Result, flag bool) {
	if flag && len(res) != 0 {
		// 创建文件夹
		dir := "00-html"
		os.MkdirAll(dir, os.ModePerm)
		fileName := fmt.Sprintf("%s/%s.html", dir, res[0].TextName)
		h := serve.NewHTML()
		for _, v := range res {
			h.AddResult(v)
		}
		h.OutputHTMLFile(fileName)
		fmt.Println("已保存 html 结果")
	}
}

func printSep() {
	fmt.Println("----------------------")
}

func getFiles(fp string) []string {
	fi, err := os.Stat(fp)
	if err != nil {
		fmt.Println("找不到文件或文件夹", fp)
		panic(err)
	}
	if fi.IsDir() {
		ret := make([]string, 0)
		files, err := os.ReadDir(fp)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if !file.IsDir() {
				ret = append(ret, filepath.Join(fp, file.Name()))
			}
		}
		// fmt.Println(fp, ret)
		return ret
	}
	return []string{fp}
}

// 交换行和列索引
func transpose[T any](A [][]T) [][]T {
	result := make([][]T, len(A[0]))
	for i := range result {
		result[i] = make([]T, len(A))
	}
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0]); j++ {
			result[j][i] = A[i][j]
		}
	}
	return result
}
