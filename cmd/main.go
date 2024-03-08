package cmd

import (
	"os"

	"github.com/nopdan/gosmq/pkg/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "这是最快的赛码器\n用于对基于码表的输入法针对特定文章进行测评\nhttps://github.com/nopdan/gosmq",
	Run: func(cmd *cobra.Command, args []string) {
		_root()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(convertCmd)
}

var logger = util.Logger

func Execute() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "server")
	}
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
