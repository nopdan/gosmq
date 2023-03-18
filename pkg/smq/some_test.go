package smq

import (
	"fmt"
	"testing"
)

func TestPuncts(t *testing.T) {
	for k, v := range zhKeysMap {
		fmt.Printf("%s\t%s\n", string(k), v)
	}
	fmt.Println(enKeysMap)
}
