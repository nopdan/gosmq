package data

import (
	"fmt"
	"testing"
)

func TestFindSuffixInt(t *testing.T) {
	fmt.Println(findSuffixInteger("aaa2"))
	fmt.Println(findSuffixInteger("aaa0"))
	fmt.Println(findSuffixInteger("aaa22"))
	fmt.Println(findSuffixInteger("aaa20"))
	fmt.Println(findSuffixInteger("aaa02"))
	fmt.Println(findSuffixInteger("aaa_"))
	fmt.Println(findSuffixInteger("aaa"))
}

func TestAddSuffix(t *testing.T) {
	d := Dict{
		Text:       &Text{String: "test"},
		SelectKeys: "_;'",
	}
	d.Init()
	fmt.Println(string(d.addSuffix("test", 1)))
	fmt.Println(string(d.addSuffix("test", 2)))
	fmt.Println(string(d.addSuffix("test", 3)))
	fmt.Println(string(d.addSuffix("test", 4)))
	fmt.Println(string(d.addSuffix("test", 10)))
}

func TestSlice(t *testing.T) {
	a := make([]byte, 4, 8)
	b := a[:3]
	b = append(b, 1, 2, 3)
	fmt.Println(a) // [0 0 0 1]
	fmt.Println(b) // [0 0 0 1 2 3]
}
