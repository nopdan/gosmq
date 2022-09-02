package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) <= 1 {
		serve(false)
	} else if os.Args[1] == "serve" {
		serve(true)
	} else {
		cli()
	}
}

func info() {
	fmt.Printf("gosmq v1.0.0 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("https://github.com/cxcn/gosmq/\n")
}
