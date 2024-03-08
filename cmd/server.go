package cmd

import (
	"github.com/nopdan/gosmq/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve(Port, Silent)
	},
}

var Port string
var Silent bool

func init() {
	serverCmd.Flags().StringVarP(&Port, "port", "p", "7172", "指定端口")
	serverCmd.Flags().BoolVarP(&Silent, "silent", "s", false, "静默启动")
}
