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
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol.0p;/'"
	for i := range keys {
		if KeyPosArr[keys[i]].Fin == 0 {
			continue
		}
		fmt.Printf("key: %s, IsLeft: %v, Finger: %d\n", string(keys[i]), KeyPosArr[keys[i]].IsLeft, KeyPosArr[keys[i]].Fin)
	}
}
