package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nopdan/gosmq/pkg/result"
)

func Output(data []*result.Result) {
	noColor := table.StyleColoredDark
	noColor.Color = table.ColorOptions{}

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
	fmt.Printf("%s\n\n", t.Render())

	fmt.Printf("非汉字：%v \n", data[0].Han.NotHan)
	for i, res := range data {
		if i == 0 && len(data) == 1 {
			fmt.Printf("缺字：  %s \n", res.Han.Lack)
			break
		}
		fmt.Printf("%d 缺字：  %s \n", i+1, res.Han.Lack)
	}
	fmt.Println()

	for i, res := range data {
		if len(data) == 1 {
			fmt.Printf("码长：")
		} else {
			fmt.Printf("%d %s：", i+1, "码长")
		}
		for j, v := range res.Dist.CodeLen {
			if v != 0 {
				fmt.Printf("%d:%d  ", j, v)
			}
		}
		fmt.Println()
		if len(data) == 1 {
			fmt.Printf("词长：")
		} else {
			fmt.Printf("%d %s：", i+1, "词长")
		}
		for j, v := range res.Dist.WordLen {
			if v != 0 {
				fmt.Printf("%d:%d  ", j, v)
			}
		}
		fmt.Println()
		if len(data) == 1 {
			fmt.Printf("选重：")
		} else {
			fmt.Printf("%d %s：", i+1, "选重")
		}
		for j, v := range res.Dist.Collision {
			if v != 0 {
				fmt.Printf("%d:%d  ", j, v)
			}
		}
		fmt.Printf("\n\n")
	}

	t = table.NewWriter()
	t.AppendHeader(table.Row{"总键数", "码长", "十击速度", "非汉字数", "计数", "缺字数", "计数"})
	for _, res := range data {
		t.AppendRow([]any{
			res.CodeLen.Total,
			fmt.Sprintf("%.4f", res.CodeLen.PerChar),
			fmt.Sprintf("%.2f", 600/res.CodeLen.PerChar),
			res.Han.NotHans, res.Han.NotHanCount,
			res.Han.Lacks, res.Han.LackCount,
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"首选词", "打词", "--", "打词字数", "--", "选重", "--", "选重字数", "--"})
	for _, res := range data {
		commitRate := func(x int) float64 {
			return div(x, res.Commit.Count)
		}
		t.AppendRow([]any{
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
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"左手", "右手", "左右", "右左", "左左", "右右"})
	for _, res := range data {
		t.AppendRow([]any{
			fmt.Sprintf("%.2f%%", 100*div(res.LeftHand, res.CodeLen.Total)),
			fmt.Sprintf("%.2f%%", 100*div(res.RightHand, res.CodeLen.Total)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.LeftToRight, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.RightToLeft, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.LeftToLeft, res.Pair.Count)),
			fmt.Sprintf("%.2f%%", 100*div(res.Pair.RightToRight, res.Pair.Count)),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"当量", "异手", "同指", "三连击", "两连击", "小跨排", "大跨排", "异指", "小指干扰", "错手"})
	for _, res := range data {
		pairRate := func(x int) float64 {
			return div(x, res.Pair.Count)
		}
		t.AppendRow([]any{
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
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "左拇指", "右拇指", "食指", "中指", "无名指", "小指"})
	for _, res := range data {
		newRow := []any{}
		for i := 1; i < 11; i++ {
			newRow = append(newRow, fmt.Sprintf("%.2f%%", 100*div(res.Dist.Finger[i], res.CodeLen.Total)))
		}
		t.AppendRow(newRow)
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	addKeys := func(res *result.Result, keys string) {
		header := []any{}
		row := []any{}
		for i := range len(keys) {
			header = append(header, string(keys[i]))
			row = append(row, fmt.Sprintf("%.2f%%", 100*res.Keys[string(keys[i])].Rate))
		}
		writer := table.NewWriter()
		writer.AppendHeader(header)
		writer.AppendRow(row)
		writer.SetStyle(table.StyleColoredBright)
		fmt.Printf("%s\n", writer.Render())
	}
	for _, res := range data {
		if len(data) != 1 {
			fmt.Printf("%s：\n", res.Info.DictName)
		}
		addKeys(res, "1234567890-=")
		addKeys(res, "qwertyuiop[]")
		addKeys(res, "asdfghjkl;'")
		addKeys(res, "zxcvbnm,./")
	}
	printSep()
}

func div(x, y int) float64 {
	return float64(x) / float64(y)
}
