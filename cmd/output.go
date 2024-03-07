package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nopdan/gosmq/pkg/result"
)

func Output(data []*result.Result) {
	// r, _ := json.MarshalIndent(res, "", "  ")
	// fmt.Printf(string(r))

	noColor := table.StyleColoredDark
	noColor.Color = table.ColorOptions{}

	out := ""

	t := table.NewWriter()
	tmpRow := table.Row{"文本名"}
	tmpRow = append(tmpRow, data[0].Info.TextName)
	tmpRow = append(tmpRow, "|")
	tmpRow = append(tmpRow, "方案名")
	for _, res := range data {
		tmpRow = append(tmpRow, res.Info.DictName)
	}
	t.AppendRow(tmpRow)
	tmpRow = table.Row{"字数"}
	tmpRow = append(tmpRow, data[0].Info.TextLen)
	tmpRow = append(tmpRow, "|")
	tmpRow = append(tmpRow, "词条数 ")
	for _, res := range data {
		tmpRow = append(tmpRow, res.Info.DictLen)
	}
	t.AppendRow(tmpRow)
	t.SetStyle(noColor)
	out += t.Render() + "\n\n"

	out += fmt.Sprintf("非汉字：%v \n", data[0].Han.NotHan)
	for i, res := range data {
		if i == 0 && len(data) == 1 {
			out += fmt.Sprintf("缺字：  %s \n", res.Han.Lack)
			break
		}
		out += fmt.Sprintf("%d 缺字：  %s \n", i+1, res.Han.Lack)
	}
	out += "\n"

	for i, res := range data {
		if len(data) == 1 {
			out += "码长："
		} else {
			out += fmt.Sprintf("%d %s：", i+1, "码长")
		}
		for j, v := range res.Dist.CodeLen {
			if v != 0 {
				out += fmt.Sprintf("%d:%d  ", j, v)
			}
		}
		out += "\n"
		if len(data) == 1 {
			out += "词长："
		} else {
			out += fmt.Sprintf("%d %s：", i+1, "词长")
		}
		for j, v := range res.Dist.WordLen {
			if v != 0 {
				out += fmt.Sprintf("%d:%d  ", j, v)
			}
		}
		out += "\n"
		if len(data) == 1 {
			out += "选重："
		} else {
			out += fmt.Sprintf("%d %s：", i+1, "选重")
		}
		for j, v := range res.Dist.Collision {
			if v != 0 {
				out += fmt.Sprintf("%d:%d  ", j, v)
			}
		}
		out += "\n\n"
	}

	t = table.NewWriter()
	t.AppendHeader(table.Row{"总键数", "码长", "十击速度", "非汉字数", "计数", "缺字数", "计数"})
	for _, res := range data {
		t.AppendRow([]interface{}{
			res.CodeLen.Total,
			fmt.Sprintf("%.4f", res.CodeLen.PerChar),
			fmt.Sprintf("%.2f", 600/res.CodeLen.PerChar),
			res.Han.NotHans, res.Han.NotHanCount,
			res.Han.Lacks, res.Han.LackCount,
		})
	}
	t.SetStyle(table.StyleColoredBright)
	out += t.Render() + "\n\n"

	t = table.NewWriter()
	t.AppendHeader(table.Row{"首选词", "打词", "--", "打词字数", "--", "选重", "--", "选重字数", "--"})
	for _, res := range data {
		commitRate := func(x int) float64 {
			return div(x, res.Commit.Count)
		}
		t.AppendRow([]interface{}{
			res.Commit.WordFirst,
			res.Commit.Word,
			fmt.Sprintf("%.2f%%", 100*commitRate(res.Commit.Word)),
			res.Commit.WordChars,
			fmt.Sprintf("%.2f%%", 100*commitRate(res.Commit.WordChars)),
			res.Commit.Collision,
			fmt.Sprintf("%.2f%%", 100*commitRate(res.Commit.Collision)),
			res.Commit.CollisionChars,
			fmt.Sprintf("%.2f%%", 100*commitRate(res.Commit.CollisionChars)),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	out += t.Render() + "\n\n"

	t = table.NewWriter()
	t.AppendHeader(table.Row{"左手", "右手", "左右", "右左", "左左", "右右"})
	for _, res := range data {
		t.AppendRow([]interface{}{
			fmt.Sprintf("%.2f%%", 100*div(res.LeftHand, res.CodeLen.Total)),
			fmt.Sprintf("%.2f%%", 100*div(res.RightHand, res.CodeLen.Total)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.LeftToRight, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.RightToLeft, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.LeftToLeft, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.RightToRight, res.Pair.Count)),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	out += t.Render() + "\n\n"

	t = table.NewWriter()
	t.AppendHeader(table.Row{"当量", "异手", "同指", "三连击", "两连击", "小跨排", "大跨排", "异指", "小指干扰", "错手"})
	for _, res := range data {
		pairRate := func(x int) float64 {
			return div(x, res.Pair.Count)
		}
		t.AppendRow([]interface{}{
			fmt.Sprintf("%.4f", pairRate(int(res.Equivalent))),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.DiffHand)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.SameFinger)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.TribleHit)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.DoubleHit)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.SingleSpan)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.MultiSpan)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.DiffFinger)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.Disturb)),
			fmt.Sprintf("%.2f%%", 100*pairRate(res.Pair.Staggered)),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	out += t.Render() + "\n\n"

	t = table.NewWriter()
	t.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "左拇指", "右拇指", "食指", "中指", "无名指", "小指"})
	for _, res := range data {
		newRow := []interface{}{}
		for i := 1; i < 11; i++ {
			newRow = append(newRow, fmt.Sprintf("%.2f%%", 100*div(res.Dist.Finger[i], res.CodeLen.Total)))
		}
		t.AppendRow(newRow)
	}
	t.SetStyle(table.StyleColoredBright)
	out += t.Render() + "\n\n"

	for _, res := range data {
		if len(data) != 1 {
			out += res.Info.DictName + "：\n"
		}
		t6_1_header := []interface{}{}
		t6_2_header := []interface{}{}
		t6_3_header := []interface{}{}
		t6_4_header := []interface{}{}
		t6_1_row := []interface{}{}
		t6_2_row := []interface{}{}
		t6_3_row := []interface{}{}
		t6_4_row := []interface{}{}
		// keys := "1234567890qwertyuiopasdfghjkl;zxcvbnm,./"
		keys1 := "1234567890"
		keys2 := "qwertyuiop"
		keys3 := "asdfghjkl;" // 单引号单独加
		keys4 := "zxcvbnm,./"
		for i := 0; i < 10; i++ {
			a, b, c, d := keys1[i], keys2[i], keys3[i], keys4[i]
			t6_1_header = append(t6_1_header, string(a))
			t6_2_header = append(t6_2_header, string(b))
			t6_3_header = append(t6_3_header, string(c))
			t6_4_header = append(t6_4_header, string(d))
			t6_1_row = append(t6_1_row, fmt.Sprintf("%.2f%%", 100*res.Keys[string(a)].Rate))
			t6_2_row = append(t6_2_row, fmt.Sprintf("%.2f%%", 100*res.Keys[string(b)].Rate))
			t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*res.Keys[string(c)].Rate))
			t6_4_row = append(t6_4_row, fmt.Sprintf("%.2f%%", 100*res.Keys[string(d)].Rate))
		}
		t6_3_header = append(t6_3_header, "'")
		t6_3_row = append(t6_3_row, fmt.Sprintf("%.2f%%", 100*res.Keys["'"].Rate))

		t6_1 := table.NewWriter()
		t6_1.AppendHeader(t6_1_header)
		t6_1.AppendRow(t6_1_row)
		t6_1.SetStyle(table.StyleColoredBright)
		out += fmt.Sprintln(t6_1.Render())

		t6_2 := table.NewWriter()
		t6_2.AppendHeader(t6_2_header)
		t6_2.AppendRow(t6_2_row)
		t6_2.SetStyle(table.StyleColoredBright)
		out += fmt.Sprintln(t6_2.Render())

		t6_3 := table.NewWriter()
		t6_3.AppendHeader(t6_3_header)
		t6_3.AppendRow(t6_3_row)
		t6_3.SetStyle(table.StyleColoredBright)
		out += fmt.Sprintln(t6_3.Render())

		t6_4 := table.NewWriter()
		t6_4.AppendHeader(t6_4_header)
		t6_4.AppendRow(t6_4_row)
		t6_4.SetStyle(table.StyleColoredBright)
		out += fmt.Sprintln(t6_4.Render())
	}
	fmt.Print(out)
	printSep()
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}
