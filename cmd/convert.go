package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nopdan/gosmq/pkg/data"
	"github.com/nopdan/gosmq/pkg/smq"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "转换赛码表",
	Run: func(cmd *cobra.Command, args []string) {
		convert()
	},
}

type Config struct {
	Path       string
	Format     string
	SelectKeys string
	PushStart  int
	Overwrite  bool
}

var config = &Config{}

func init() {
	convertCmd.Flags().StringVarP(&config.Path, "dict", "i", "", "待转换的码表")
	convertCmd.Flags().StringVarP(&config.Format, "format", "f", "jisu", "待转换码表的格式")
	convertCmd.Flags().StringVarP(&config.SelectKeys, "select", "k", "_;'", "自定义选重键")
	convertCmd.Flags().IntVarP(&config.PushStart, "push", "p", 4, "起顶码长")
	convertCmd.Flags().BoolVarP(&config.Overwrite, "overwrite", "o", false, "码表已存在时是否覆盖")
}

func convert() {
	// 命令行模式
	if config.Path != "" {
		config.convert()
		return
	}

	handle := func(err error) {
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
	}

	// 交互模式
	err := survey.AskOne(&survey.Input{
		Message: "待转换的码表:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &config.Path, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Select{
		Message: "码表格式:",
		Options: []string{"极速赛码表", "多多(Rime)", "冰凌", "小小(极点)"},
	}, &config.Format)
	handle(err)

	mFormat := make(map[string]string)
	mFormat["极速赛码表"] = "jisu"
	mFormat["多多(Rime)"] = "duoduo"
	mFormat["冰凌"] = "bingling"
	mFormat["小小(极点)"] = "xiaoxiao"
	config.Format = mFormat[config.Format]

	if config.Format != "jisu" {
		err = survey.AskOne(&survey.Input{
			Message: "起顶码长(码长大于等于此数，首选不会追加空格):",
			Default: "4",
		}, &config.PushStart)
		handle(err)
	}

	err = survey.AskOne(&survey.Input{
		Message: "自定义选重键:",
		Default: "_;'",
	}, &config.SelectKeys)
	handle(err)

	fmt.Println()
	fmt.Printf("config: %v\n", conf)
	config.convert()
}

func (c *Config) convert() {
	smq := &smq.Config{}
	dict := &data.Dict{
		Text:       &data.Text{Path: c.Path},
		Format:     c.Format,
		SelectKeys: c.SelectKeys,
		Push:       c.PushStart,
		Overwrite:  c.Overwrite,
	}
	smq.AddDict(dict)
}
