package gen

import (
	"fmt"
	"testing"
)

func TestFindSuffixInt(t *testing.T) {
	fmt.Println(FindSuffixInteger("aaa2"))
	fmt.Println(FindSuffixInteger("aaa0"))
	fmt.Println(FindSuffixInteger("aaa22"))
	fmt.Println(FindSuffixInteger("aaa20"))
	fmt.Println(FindSuffixInteger("aaa02"))
	fmt.Println(FindSuffixInteger("aaa_"))
	fmt.Println(FindSuffixInteger("aaa"))
}
