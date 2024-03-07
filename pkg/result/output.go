package result

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
)

// 输出分词结果
func (res *Result) OutputSplit() {
	if len(res.segments) == 0 {
		return
	}
	slices.SortFunc(res.segments, func(i, j segment) int {
		return cmp.Compare(i.PartIdx, j.PartIdx)
	})
	// fmt.Printf("Segments: %+v\n", res.segments)
	// 创建文件夹
	dir := "02-分词结果"
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s_.txt", dir, res.Info.DictName, res.Info.TextName)
	f, _ := os.Create(fileName)
	defer f.Close()
	buf := bufio.NewWriterSize(f, 1024*1024)
	for i := range res.segments {
		for j := range res.segments[i].Segment {
			buf.WriteString(res.segments[i].Segment[j].Word)
			buf.WriteByte('\t')
			buf.WriteString(res.segments[i].Segment[j].Code)
			buf.WriteByte('\n')

			// fmt.Printf("%s\t%s\n", res.segments[i].Segment[j].Word, res.segments[i].Segment[j].Code)
		}
	}
	buf.Flush()
}

// 输出词条统计数据
func (res *Result) OutputStat() {
	if len(res.statData) == 0 {
		return
	}

	type detail struct {
		word string
		*CodePosCount
	}
	details := make([]detail, 0, len(res.statData))
	for k, v := range res.statData {
		details = append(details, detail{k, v})
	}
	slices.SortStableFunc(details, func(i, j detail) int {
		return cmp.Compare(j.Count, i.Count)
	})

	// 创建文件夹
	dir := "01-词条统计"
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s.txt", dir, res.Info.DictName, res.Info.TextName)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("create %s error: %v\n", fileName, err)
	}
	defer f.Close()
	buf := bufio.NewWriterSize(f, 1024*1024)
	buf.WriteString("词条\t编码\t选重\t次数\n")

	for _, v := range details {
		buf.WriteString(v.word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Pos))
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Count))
		buf.WriteByte('\n')
	}
	buf.Flush()
}
