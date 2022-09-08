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

func OutputDetail(textName string, res *Result) {
	// 创建文件夹
	os.Mkdir("result", 0666)
	// 输出分词结果
	var buf strings.Builder
	for i := 0; i < len(res.Data.CodeSlice); i++ {
		buf.WriteString(fmt.Sprintf("%s\t%s\n", res.Data.WordSlice[i], string(res.Data.CodeSlice[i])))
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_分词结果.txt", textName, res.Name), []byte(buf.String()), 0666)
	// 输出词条数据
	buf.Reset()
	buf.WriteString("词条\t编码\t顺序\t次数\n")
	type details struct {
		CoC
		word string
	}
	tmp := make([]details, 0, len(res.Data.Details))
	for k, v := range res.Data.Details {
		tmp = append(tmp, details{*v, k})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Count > tmp[j].Count
	})
	for _, v := range tmp {
		buf.WriteString(v.word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Order))
		buf.WriteByte('\t')
		buf.WriteString(strconv.Itoa(v.Count))
		buf.WriteByte('\n')
	}
	os.WriteFile(fmt.Sprintf("result/%s_%s_词条数据.txt", textName, res.Name), []byte(buf.String()), 0666)
	// 输出 json 数据
	res.Data.CodeSlice = []string{}
	res.Data.WordSlice = []string{}
	res.Data.Details = make(map[string]*CoC)
	tmp2, _ := json.MarshalIndent(res, "", "  ")
	os.WriteFile(fmt.Sprintf("result/%s_%s.json", textName, res.Name), tmp2, 0666)
	fmt.Println("已输出详细数据，请查看 result 文件夹")
}
