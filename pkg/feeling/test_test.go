package feeling

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	a := Comb["cf"]
	fmt.Println(a.DL)
	a = Comb["fr"]
	fmt.Println(a.IsXKP)
	a = Comb["bt"]
	fmt.Println(a.IsDKP)
	a = Comb["xe"]
	fmt.Println(a.IsCS)
	a = Comb["aw"]
	fmt.Println(a.IsXZGR)
}

func TestKeyPos(t *testing.T) {
	fmt.Println(KeyPosMap['g'])
	fmt.Println(KeyPosMap['m'])
	fmt.Println(KeyPosMap['#'])
}
