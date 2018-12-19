package main

import (
	"fmt"
)

func main() {
	reg0 := 0
	reg3 := 1
	//reg3_loop
	for {
		reg5 := 1
		for {
			if reg5*reg3 == 915 {
				fmt.Printf("pair %d %d\n", reg5, reg3)
				reg0 += reg3
			}
			reg5++
			if reg5 > 915 {
				break
			}
		}
		reg3++
		if reg3 > 915 {
			break
		}
	}
	fmt.Printf("%d\n", reg0)
}
