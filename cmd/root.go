package cmd

import (
	"fmt"
	"os"

	"github.com/imetool/gosmq/internal/serve"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smq",
	Short: "gosmq 是一个非常快的赛码器。",
	Long:  "对基于码表的输入法针对特定文章进行测评\nhttps://github.com/imetool/gosmq",
	Run: func(cmd *cobra.Command, args []string) {
		goWithSurvey()
	},
}

func init() {
	rootCmd.AddCommand(goCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(multiCmd)
}

func Execute() {
	if len(os.Args) <= 1 {
		serve.Serve("7172", false)
		return
	}
	if len(os.Args) == 2 && os.Args[1] == "gen" {
		fmt.Println("交互模式")
		genWithSurvey()
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printSep() {
	fmt.Println("----------------------")
}
