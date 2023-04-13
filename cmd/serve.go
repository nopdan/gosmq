package cmd

import (
	"github.com/flowerime/gosmq/internal/serve"
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
	serveCmd.Flags().StringVarP(&Port, "port", "p", "7172", "指定端口")
	serveCmd.Flags().BoolVarP(&Silent, "silent", "s", false, "静默启动")
}
