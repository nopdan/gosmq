package data

import (
	"bufio"
	"cmp"
	"os"
	"path/filepath"
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
func (d *Dict) output(dict []*Entry) {
	path := filepath.Join("dict", d.Text.Name+".txt")
	// 判断文件是否存在
	_, err := os.Stat(path)
	// 存在且不覆盖
	if err == nil && !d.Overwrite {
		logger.Warn("赛码表已经存在\n", "path", path)
		return
	} else if d.Overwrite {
		logger.Info("覆盖赛码表\n", "path", path)
	}
	// 按照词长排序
	slices.SortStableFunc(dict, func(i, j *Entry) int {
		return cmp.Compare(
			utf8.RuneCountInString(j.Word),
			utf8.RuneCountInString(i.Word),
		)
	})

	// 创建文件
	f, err := os.Create(path)
	if err != nil {
		logger.Warn("创建文件失败", "path", path, "error", err)
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
	logger.Info("输出赛码表", "path", path)
}
