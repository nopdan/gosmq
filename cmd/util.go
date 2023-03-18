package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

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
