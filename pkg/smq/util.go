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

func OutputDetail(textName string, res *Result, mr *matchRes) {

	type CodePosCount struct {
		Code  string
		Pos   int
		Count int
	}
	type detail struct {
		word string
		*CodePosCount
	}
	wordMap := make(map[string]*CodePosCount, len(mr.wordSlice)/5)

	// 输出分词结果
	var buf strings.Builder
	for i, word := range mr.wordSlice {
		buf.WriteString(fmt.Sprintf("%s\t%s\n", word, mr.codeSlice[i]))

		if _, ok := wordMap[word]; !ok {
			wordMap[word] = &CodePosCount{mr.codeSlice[i], mr.pos[i], 1}
		} else {
			wordMap[word].Count++
		}
	}
	// 创建文件夹
	os.MkdirAll("result", os.ModePerm)
	os.WriteFile(fmt.Sprintf("result/分词结果_%s_%s_.txt", res.Name, textName), []byte(buf.String()), 0666)

	// 输出词条数据
	buf.Reset()
	buf.WriteString("词条\t编码\t选重\t次数\n")
	details := make([]detail, 0, len(wordMap))
	for k, v := range wordMap {
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
	os.WriteFile(fmt.Sprintf("result/词条数据_%s_%s.txt", res.Name, textName), []byte(buf.String()), 0666)

	// 输出 json 数据
	tmp3, _ := json.MarshalIndent(res, "", "  ")
	os.WriteFile(fmt.Sprintf("result/data_%s_%s.json", res.Name, textName), tmp3, 0666)
	fmt.Println("已输出详细数据，请查看 result 文件夹")
}
