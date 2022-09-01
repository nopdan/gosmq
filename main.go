package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) == 1 {
		serve()
	} else {
		cli()
	}
}

func info() {
	fmt.Printf("gosmq v0.28 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("https://github.com/cxcn/gosmq/\n")
}
