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
	// 创建文件夹
	os.MkdirAll("result", os.ModePerm)
	// 输出分词结果
	var buf strings.Builder

	type details struct {
		CodePosCount
		word string
	}
	tmp := make(map[string]*CodePosCount, len(mr.wordSlice)/5)

	for i := range mr.wordSlice {
		buf.WriteString(fmt.Sprintf("%s\t%s\n", mr.wordSlice[i], mr.codeSlice[i]))

		if _, ok := tmp[mr.wordSlice[i]]; !ok {
			tmp[mr.wordSlice[i]] = &CodePosCount{mr.codeSlice[i], mr.pos[i], 1}
		} else {
			tmp[mr.wordSlice[i]].Count++
		}
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_分词结果.txt", textName, res.Name), []byte(buf.String()), 0666)

	// 输出词条数据
	buf.Reset()
	buf.WriteString("词条\t编码\t选重\t次数\n")
	tmp2 := make([]details, 0, len(tmp))
	for k, v := range tmp {
		tmp2 = append(tmp2, details{*v, k})
	}
	sort.Slice(tmp2, func(i, j int) bool {
		return tmp2[i].Count > tmp2[j].Count
	})
	for _, v := range tmp2 {
		buf.WriteString(v.word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Pos))
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Count))
		buf.WriteByte('\n')
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_词条数据.txt", textName, res.Name), []byte(buf.String()), 0666)
	// 输出 json 数据
	tmp3, _ := json.MarshalIndent(res, "", "  ")
	os.WriteFile(fmt.Sprintf("result/%s_%s.json", textName, res.Name), tmp3, 0666)
	fmt.Println("已输出详细数据，请查看 result 文件夹")
}
