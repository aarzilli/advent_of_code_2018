package main

import (
	"fmt"
)

//const N = 915
const N = 10551315

func main() {
	reg0 := 0
	for i := 1; i <= N; i++ {
		if N%i != 0 {
			continue
		}
		fmt.Printf("exact divisor %d\n", i)
		reg0 += i
	}
	fmt.Printf("%d\n", reg0)
}
