package cmd

import (
	"github.com/imetool/gosmq/internal/serve"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(Port, Silent)
	},
}

var Port string
var Silent bool

func init() {
	serveCmd.PersistentFlags().StringVarP(&Port, "port", "p", "7172", "指定端口")
	serveCmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "静默启动")
}
