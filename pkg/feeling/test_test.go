package feeling

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	a := Comb['c']['f']
	fmt.Println(a.Equivalent)
	a = Comb['f']['r']
	fmt.Println(a.SingleSpan)
	a = Comb['b']['t']
	fmt.Println(a.MultiSpan)
	a = Comb['x']['e']
	fmt.Println(a.Staggered)
	a = Comb['a']['w']
	fmt.Println(a.Disturb)
	a = Comb['c']['c']
	fmt.Println(a.Staggered)
}

func TestKeyPos(t *testing.T) {
	fmt.Println(KeyPos('A'))
	fmt.Println(KeyPos('A'))
	keys := "1qaz2wsx3edc4rfv5tgb_+6yhn7ujm8ik,9ol.0p;/'[]-="
	for i := range keys {
		isLeft, finger := KeyPos(keys[i])
		fmt.Printf("key: %v, IsLeft: %v, Finger: %d\n",
			string(keys[i]), isLeft, finger)
	}
}
