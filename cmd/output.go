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
	if len(data) == 1 {
		t.AppendRow(table.Row{"文本", data[0].Info.TextName, "字数", data[0].Info.TextLen}, table.RowConfig{})
		for _, res := range data {
			t.AppendRow(table.Row{"码表", res.Info.DictName, "词条数", res.Info.DictLen})
		}
		t.SetStyle(noColor)
		fmt.Printf("%s\n\n", t.Render())
	} else {
		tmpRow := table.Row{"文本"}
		tmpRow = append(tmpRow, data[0].Info.TextName)
		tmpRow = append(tmpRow, "|  码表")
		for _, res := range data {
			tmpRow = append(tmpRow, res.Info.DictName)
		}
		t.AppendRow(tmpRow)
		tmpRow = table.Row{"字数"}
		tmpRow = append(tmpRow, data[0].Info.TextLen)
		tmpRow = append(tmpRow, "|  词条数 ")
		for _, res := range data {
			tmpRow = append(tmpRow, res.Info.DictLen)
		}
		t.AppendRow(tmpRow)
		t.SetStyle(noColor)
		fmt.Printf("%s\n\n", t.Render())
	}

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
			res.Keys.Count,
			fmt.Sprintf("%.4f", res.Keys.CodeLen),
			fmt.Sprintf("%.2f", 600/res.Keys.CodeLen),
			res.Han.NotHans, res.Han.NotHanCount,
			res.Han.Lacks, res.Han.LackCount,
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"首选词", "打词", "--", "打词字数", "--", "选重", "--", "选重字数", "--"})
	for _, res := range data {
		t.AppendRow([]any{
			res.Commit.WordFirst,
			res.Commit.Word,
			commitRate(res.Commit.Word, res),
			res.Char.Word,
			charRate(res.Char.Word, res),
			res.Commit.Collision,
			commitRate(res.Commit.Collision, res),
			res.Char.Collision,
			charRate(res.Char.Collision, res),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"左手", "右手", "左右", "右左", "左左", "右右"})
	for _, res := range data {
		t.AppendRow([]any{
			keyRate(res.Keys.LeftHand, res),
			keyRate(res.Keys.RightHand, res),
			pairRate(res.Pair.LeftToRight, res),
			pairRate(res.Pair.RightToLeft, res),
			pairRate(res.Pair.LeftToLeft, res),
			pairRate(res.Pair.RightToRight, res),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"当量", "异手", "同指", "三连击", "两连击", "小跨排", "大跨排", "异指", "小指干扰", "错手"})
	for _, res := range data {
		t.AppendRow([]any{
			fmt.Sprintf("%.4f", float64(res.Pair.Equivalent)/float64(res.Pair.Count)),
			pairRate(res.Pair.DiffHand, res),
			pairRate(res.Pair.SameFinger, res),
			pairRate(res.Pair.TribleHit, res),
			pairRate(res.Pair.DoubleHit, res),
			pairRate(res.Pair.SingleSpan, res),
			pairRate(res.Pair.MultiSpan, res),
			pairRate(res.Pair.DiffFinger, res),
			pairRate(res.Pair.Disturb, res),
			pairRate(res.Pair.Staggered, res),
		})
	}
	t.SetStyle(table.StyleColoredBright)
	fmt.Printf("%s\n\n", t.Render())

	t = table.NewWriter()
	t.AppendHeader(table.Row{"小指", "无名指", "中指", "食指", "左拇指", "右拇指", "食指", "中指", "无名指", "小指"})
	for _, res := range data {
		newRow := []any{}
		for i := 1; i < 11; i++ {
			newRow = append(newRow, keyRate(res.Dist.Finger[i], res))
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
			row = append(row, keyRate(res.Dist.Key[string(keys[i])], res))
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
	fmt.Println("----------------------")
}

func commitRate(count int, res *result.Result) string {
	rate := float64(count) / float64(res.Commit.Count)
	return fmt.Sprintf("%.2f%%", 100*rate)
}

func charRate(count int, res *result.Result) string {
	rate := float64(count) / float64(res.Char.Count)
	return fmt.Sprintf("%.2f%%", 100*rate)
}

func keyRate(count int, res *result.Result) string {
	rate := float64(count) / float64(res.Keys.Count)
	return fmt.Sprintf("%.2f%%", 100*rate)
}

func pairRate(count int, res *result.Result) string {
	rate := float64(count) / float64(res.Pair.Count)
	return fmt.Sprintf("%.2f%%", 100*rate)
}
