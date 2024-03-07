package result

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 输出分词结果
func (res *Result) OutputSplit() {
	if len(res.wcIdxs) == 0 {
		return
	}
	// 创建文件夹
	dir := "02-分词结果"
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s_.txt", dir, res.DictName, res.TextName)
	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)

	sort.Slice(res.wcIdxs, func(j, k int) bool {
		return res.wcIdxs[j].idx < res.wcIdxs[k].idx
	})
	for i := range res.wcIdxs {
		var buf strings.Builder
		for j := range res.wcIdxs[i].wordSli {
			buf.WriteString(res.wcIdxs[i].wordSli[j])
			buf.Write([]byte{'\t'})
			buf.WriteString(res.wcIdxs[i].codeSli[j])
			buf.Write([]byte{'\n'})
		}
		f.WriteString(buf.String())
	}
	f.Close()
	// 清空 wcIdxs
	res.wcIdxs = make([]wcIdx, 0)
}

// 输出词条统计数据
func (res *Result) OutputStat() {
	if len(res.statData) == 0 {
		return
	}
	// 创建文件夹
	dir := "01-词条统计"
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf("%s/%s_%s.txt", dir, res.DictName, res.TextName)

	type detail struct {
		word string
		*CodePosCount
	}
	var buf strings.Builder
	buf.WriteString("词条\t编码\t选重\t次数\n")
	details := make([]detail, 0, len(res.statData))
	for k, v := range res.statData {
		details = append(details, detail{k, v})
	}
	sort.Slice(details, func(i, j int) bool {
		return details[i].Count > details[j].Count
	})
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
	os.WriteFile(fileName, []byte(buf.String()), 0666)
	res.statData = make(map[string]*CodePosCount)
}
