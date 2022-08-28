package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cxcn/gosmq/web"
)

func main() {

	if len(os.Args) == 1 {
		web.Run()
	} else {
		cli()
	}
}

func printInfo() {
	fmt.Printf("gosmq v0.26 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("https://github.com/cxcn/gosmq/\n")
}
