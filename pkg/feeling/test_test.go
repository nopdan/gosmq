package feeling

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	a := Comb['c']['f']
	fmt.Println(a & IsDKP)
	a = Comb['f']['r']
	fmt.Println(a & IsXKP)
	a = Comb['b']['t']
	fmt.Println(a & IsDKP)
	a = Comb['x']['e']
	fmt.Println(a & IsCS)
	a = Comb['a']['w']
	fmt.Println(a & IsXZGR)
	a = Comb['c']['c']
	fmt.Println(a & IsCS)
}

func TestKeyPos(t *testing.T) {
	fmt.Println(KeyPosArr['g'])
	fmt.Println(KeyPosArr['m'])
	fmt.Println(KeyPosArr['#'])
}
