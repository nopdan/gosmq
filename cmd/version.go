package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "输出版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gosmq v1.0.0 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
