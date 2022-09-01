package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	_ "embed"

	"github.com/cxcn/gosmq/pkg/smq"
)

//go:embed assets/tmpl.html
var tmpl string

// 赛码结果
type Result struct {
	smq.Result
	KeyHeatMap [][]template.HTML
	FinHeatMap [10]template.HTML
}

// 供模版使用的数据
type TmplData struct {
	TextName    string
	TextLen     int
	NotHanCount int
	Results     []*Result
}

// 初始化，接收文本名
func NewHTML(s string) *TmplData {
	ret := new(TmplData)
	if strings.ContainsRune(s, '《') {
		ret.TextName = s
	} else {
		ret.TextName = "《" + s + "》"
	}
	return ret
}

// 添加一个结果
func (d *TmplData) AddResult(res *smq.Result) {

	d.TextLen = res.Basic.TextLen
	d.NotHanCount = res.Basic.NotHanCount

	tmp := new(Result)
	tmp.Result = *res

	tmp.genKeyHeatMap()
	tmp.genFinHeatMap()
	d.Results = append(d.Results, tmp)
}

// 输出 html 文件
func (d *TmplData) OutputHTMLFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	d.OutputHTML(file)
}

func (d *TmplData) OutputHTML(w io.Writer) {
	funcMap := template.FuncMap{"toPer": toPercentage}
	t := template.New("tmpl.html").Funcs(funcMap)
	_, err := t.Parse(tmpl)
	if err != nil {
		panic(err)
	}
	t.Execute(w, d)
}

// float64转换成百分数
func toPercentage(src float64) string {
	return fmt.Sprintf("%.2f%%", src*100.0)
}

// 生成按键热力图
func (res *Result) genKeyHeatMap() {
	keys := res.Keys
	max := 0.07
	res.KeyHeatMap = make([][]template.HTML, 4)
	line1 := "1234567890 "
	line2 := "qwertyuiop "
	line3 := "asdfghjkl;'"
	line4 := "zxcvbnm,./ "
	for j := 0; j < 11; j++ {
		res.KeyHeatMap[0] = append(res.KeyHeatMap[0], genKeyHeatCode(keys[string(line1[j])], max, line1[j]))
		res.KeyHeatMap[1] = append(res.KeyHeatMap[1], genKeyHeatCode(keys[string(line2[j])], max, line2[j]))
		res.KeyHeatMap[2] = append(res.KeyHeatMap[2], genKeyHeatCode(keys[string(line3[j])], max, line3[j]))
		res.KeyHeatMap[3] = append(res.KeyHeatMap[3], genKeyHeatCode(keys[string(line4[j])], max, line4[j]))
	}
}

// 按键颜色代码片段
func genKeyHeatCode(sk *smq.CaR, max float64, key byte) template.HTML {
	if key == ' ' {
		return template.HTML("")
	}
	var freq float64
	if sk != nil {
		freq = sk.Rate
	}
	return template.HTML(fmt.Sprintf(
		`<td class="key" style="background-color: rgba(255,0,0,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
		freq/max*0.6, string(key), freq*100))
}

// 生成手指热力图
func (res *Result) genFinHeatMap() {
	src := res.Fingers.Dist
	max := 0.25
	fins := []string{"左小", "左无", "左中", "左食", "大拇指", "右食", "右中", "右无", "右小"}
	for i := 0; i < 9; i++ {
		res.FinHeatMap[i] = genFinHeatCode(src[i+1].Rate, max, i, fins[i])
	}
}

// 手指颜色代码片段
func genFinHeatCode(freq, max float64, id int, fin string) template.HTML {
	if id == 4 {
		return template.HTML(fmt.Sprintf(
			`<td class="key fin" colspan="2" style="background-color: rgba(0,0,255,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
			freq/max*0.6, fin, freq*100))
	}
	return template.HTML(fmt.Sprintf(
		`<td class="key fin" style="background-color: rgba(0,0,255,%.4f);">%s <div class="heatMapRate">%.2f</div></td>`,
		freq/max*0.6, fin, freq*100))
}
