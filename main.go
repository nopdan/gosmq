package main

import (
	"os"

	"github.com/nopdan/gosmq/cmd"
	"github.com/nopdan/gosmq/pkg/server"
)

func main() {
	_ = os.MkdirAll("dict", os.ModePerm)
	_ = os.MkdirAll("text", os.ModePerm)
	if len(os.Args) < 2 {
		server.Serve("7007", false)
	} else {
		cmd.Execute()
	}
}
