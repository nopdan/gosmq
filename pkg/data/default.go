package data

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Entry struct {
	Word string
	Code string
	Pos  int
}

// 默认格式
func (d *Dict) load() {
	scan := bufio.NewScanner(d.Text.reader)
	for scan.Scan() {
		wc := strings.Split(scan.Text(), "\t")
		pos := 1
		if len(wc) >= 3 {
			pos, _ = strconv.Atoi(wc[2])
		} else if len(wc) < 2 {
			continue
		}
		d.insert(wc[0], wc[1], pos)
	}
}

// 向 matcher 中添加一个词条
func (d *Dict) insert(word, code string, pos int) {
	if d.Single && utf8.RuneCountInString(word) != 1 {
		return
	}
	d.Matcher.Insert(word, code, pos)
	d.Length++
}

// 输出赛码表
func output(dict []*Entry, path string) {
	// 判断文件是否存在，若存在则直接退出
	_, err := os.Stat(path)
	if err == nil {
		fmt.Printf("赛码表已经存在：%s\n", path)
		return
	}
	// 按照词长排序
	slices.SortStableFunc(dict, func(i, j *Entry) int {
		return cmp.Compare(
			utf8.RuneCountInString(i.Word),
			utf8.RuneCountInString(j.Word),
		)
	})

	// 创建文件
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("创建赛码表失败：", err)
		return
	}
	defer f.Close()

	buf := bufio.NewWriterSize(f, 1024*1024)
	for _, entry := range dict {
		buf.WriteString(entry.Word)
		buf.WriteByte('\t')
		buf.WriteString(entry.Code)
		if entry.Pos != 1 {
			buf.WriteByte('\t')
			buf.WriteString(strconv.Itoa(entry.Pos))
		}
		buf.WriteByte('\n')
	}
	buf.Flush()
	fmt.Println("输出赛码表成功：", path)
}
