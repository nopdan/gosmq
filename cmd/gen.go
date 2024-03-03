package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nopdan/gosmq/pkg/util"
	"github.com/nopdan/gosmq/internal/gen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "转换赛码表",
	Run: func(cmd *cobra.Command, args []string) {
		_gen()
	},
}

var Config gen.Config

func init() {
	genCmd.Flags().StringVarP(&Config.Path, "dict", "i", "", "待转换的码表")
	genCmd.Flags().StringVarP(&Config.Format, "format", "f", "jisu", "待转换码表的格式")
	genCmd.Flags().StringVarP(&Config.SelectKeys, "select", "k", "_;'", "自定义选重键")
	genCmd.Flags().IntVarP(&Config.PushStart, "push", "p", 4, "起顶码长")
	genCmd.Flags().BoolVarP(&Config.SortByWordLen, "sort", "s", false, "按照词长重新排序")
}

func _gen() {
	// 命令行模式
	if Config.Path != "" {
		gen_write(Config)
		return
	}

	// 交互模式
	var conf gen.Config
	handle := func(err error) {
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
	}

	err := survey.AskOne(&survey.Input{
		Message: "待转换的码表:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &conf.Path, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Select{
		Message: "码表格式:",
		Options: []string{"极速赛码表", "多多", "冰凌"},
	}, &conf.Format)
	handle(err)

	mFormat := make(map[string]string)
	mFormat["极速赛码表"] = "jisu"
	mFormat["多多"] = "duoduo"
	mFormat["冰凌"] = "bingling"
	conf.Format = mFormat[conf.Format]

	if conf.Format != "jisu" {
		err = survey.AskOne(&survey.Input{
			Message: "起顶码长(码长大于等于此数，首选不会追加空格):",
			Default: "4",
		}, &conf.PushStart)
		handle(err)
	}

	err = survey.AskOne(&survey.Input{
		Message: "自定义选重键:",
		Default: "_;'",
	}, &conf.SelectKeys)
	handle(err)

	if conf.Format != "jisu" {
		err = survey.AskOne(&survey.Confirm{
			Message: "按照词长重新排序",
			Default: false,
		}, &conf.SortByWordLen)
		handle(err)
	}

	fmt.Printf("\nconf: %v\n", conf)
	gen_write(conf)
}

func gen_write(conf gen.Config) {
	table := conf.Gen()
	path := "dict/" + util.GetFileName(conf.Path) + ".txt"
	gen.Write(table, path)
}
