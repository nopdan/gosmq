package cmd

import (
	"github.com/nopdan/gosmq/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(port, silent, prefix)
	},
}

var port int
var silent bool
var prefix string

func init() {
	serverCmd.Flags().IntVarP(&port, "port", "p", 7007, "指定端口")
	serverCmd.Flags().BoolVarP(&silent, "silent", "s", false, "静默启动")
	serverCmd.Flags().StringVarP(&prefix, "prefix", "d", "", "工作目录")
}
