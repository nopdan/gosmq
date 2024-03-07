package dict

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/nopdan/gosmq/pkg/matcher"
)

type Entry struct {
	Word string
	Code string
	Pos  int
}

// 初始化 Dict
func (d *Dict) init() {
	// 匹配算法
	if d.Single {
		d.Matcher = matcher.NewSingle()
	} else {
		switch d.algorithm {
		case "greedy", "":
			d.Matcher = matcher.NewTrie(false)
		case "ordered":
			d.Matcher = matcher.NewTrie(true)
		case "dynamic":
			// TODO
			fallthrough
		default:
			panic("不支持的匹配算法: " + d.algorithm)
		}
	}

	var dict []*Entry
	// 读取码表，构建 matcher
	switch d.format {
	case "default", "":
		d.load()
	case "jisu", "js":
		dict = d.loadJisu()
	case "duoduo", "dd", "rime":
		dict = d.loadTSV(true)
	case "bingling", "bl":
		dict = d.loadTSV(false)
	case "xiaoxiao", "xx", "jidian", "jd":
		dict = d.loadXiao()
	default:
		panic("不支持的格式: " + d.format)
	}
	// 输出转换后的赛码表
	var wg sync.WaitGroup
	if dict != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Output(dict, filepath.Join("dict",
				strings.TrimSuffix(d.Name, ".txt")+".txt"),
			)
		}()
	}
	d.Matcher.Build()
	if dict != nil {
		wg.Wait()
	}
}

// 默认格式
func (d *Dict) load() {
	scan := bufio.NewScanner(d.Reader)
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
func Output(dict []*Entry, path string) {
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
	var buf bytes.Buffer
	buf.Grow(len(dict))
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

	err = os.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("输出赛码表失败：", err)
	} else {
		fmt.Println("输出赛码表成功：", path)
	}
}
