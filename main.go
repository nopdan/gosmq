package main

import (
	"os"

	"github.com/flowerime/gosmq/cmd"
)

func main() {
	os.MkdirAll("dict", os.ModePerm)
	os.MkdirAll("text", os.ModePerm)
	cmd.Execute()
}
