package util

import (
	"fmt"
	"testing"
)

func TestWalkDir(t *testing.T) {
	res := WalkDir(`D:\Code\go\gosmq\build\text`)
	fmt.Printf("res: %v\n", res)

	res = WalkDir(`..\..\sample`)
	fmt.Printf("res: %v\n", res)
}

func TestWalkDirWithSuffix(t *testing.T) {
	res := WalkDirWithSuffix(`..\`, "_test.go")
	fmt.Printf("res: %v\n", res)
}
