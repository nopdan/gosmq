package result

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/nopdan/gosmq/pkg/util"
)

var logger = util.Logger

// 输出分词结果
func (res *Result) OutputSplit() {
	if len(res.segments) == 0 {
		return
	}
	slices.SortFunc(res.segments, func(i, j segment) int {
		return cmp.Compare(i.PartIdx, j.PartIdx)
	})
	// 创建文件夹
	dir := "02-分词结果"
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s_.txt", dir, res.Info.DictName, res.Info.TextName)
	f, err := os.Create(fileName)
	if err != nil {
		logger.Warn("保存分词结果失败", "error", err)
		return
	}
	defer f.Close()
	buf := bufio.NewWriterSize(f, 1024*1024)
	for i := range res.segments {
		for j := range res.segments[i].Segment {
			buf.WriteString(res.segments[i].Segment[j].Word)
			buf.WriteByte('\t')
			buf.WriteString(res.segments[i].Segment[j].Code)
			buf.WriteByte('\n')
		}
	}
	buf.Flush()
	logger.Info("保存分词结果成功", "path", fileName)
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
		logger.Warn("保存词条统计数据失败", "error", err)
		return
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
	logger.Info("保存词条统计数据成功", "path", fileName)
}

// 输出 json 数据
func (res *Result) OutPutJson() {
	// 创建文件夹
	dir := "00-data"
	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s.json", dir, res.Info.DictName, res.Info.TextName)

	tmp, _ := json.MarshalIndent(res, "", "  ")
	err := os.WriteFile(fileName, tmp, 0666)
	if err != nil {
		logger.Warn("保存 json 数据失败", "error", err)
		return
	}
	logger.Info("保存 json 数据成功", "path", fileName)
}
