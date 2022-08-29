package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {

	if len(os.Args) == 1 {
		web()
	} else if len(os.Args) == 2 && os.Args[1] == "serve" {
		web()
	} else {
		cli()
	}
}

func printInfo() {
	fmt.Printf("gosmq v0.26 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("https://github.com/cxcn/gosmq/\n")
}
