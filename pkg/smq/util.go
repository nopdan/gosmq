package smq

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func div(x, y int) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(x)/float64(y)), 64)
	return value
}

func AddTo(sli *[]int, pos int) {
	for pos > len(*sli)-1 {
		*sli = append(*sli, 0)
	}
	(*sli)[pos]++
}

func AddToVal(sli *[]int, pos int, val int) {
	for pos > len(*sli)-1 {
		*sli = append(*sli, 0)
	}
	(*sli)[pos] += val
}

const (
	S_SPLIT = 1 << iota
	S_STAT
	S_JSON
)

func (res *Result) Output(flag int) {

	// 输出分词结果
	if flag&S_SPLIT != 0 && len(res.wcIdxs) == 0 {
		// 创建文件夹
		dir := "02-分词结果"
		os.MkdirAll(dir, os.ModePerm)
		fileName := fmt.Sprintf("%s/%s_%s_.txt", dir, res.DictName, res.TextName)
		f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)

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
		fmt.Println("已输出分词结果")
	}

	// 输出词条统计数据
	if flag&S_STAT != 0 {
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
		fmt.Println("已输出词条统计数据")
	}

	// 输出 json 数据
	if flag&S_JSON != 0 {
		// 创建文件夹
		dir := "00-data"
		os.MkdirAll(dir, os.ModePerm)
		fileName := fmt.Sprintf("%s/%s_%s.json", dir, res.DictName, res.TextName)

		tmp, _ := json.MarshalIndent(res, "", "  ")
		os.WriteFile(fileName, tmp, 0666)
		fmt.Println("已输出 json 数据")
	}

}
