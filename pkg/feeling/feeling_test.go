package feeling

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	a := combination['c']['f']
	fmt.Println(a.Equivalent)
	a = combination['f']['r']
	fmt.Println(a.SingleSpan)
	a = combination['b']['t']
	fmt.Println(a.MultiSpan)
	a = combination['x']['e']
	fmt.Println(a.Staggered)
	a = combination['a']['w']
	fmt.Println(a.Disturb)
	a = combination['c']['c']
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

type Result struct{}

func TestRace(t *testing.T) {
	res := make([][]*Result, 3)
	fmt.Println(res)
	for i := range 3 {
		fmt.Println(res[i])
		res[i] = append(res[i], &Result{}, &Result{}, &Result{})
		fmt.Println(res[i])
		res[i] = make([]*Result, 4)
		fmt.Println(res[i])
	}
}
