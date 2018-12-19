package main

import (
	"fmt"
)

const REG1 = 915

func main() {
	reg0 := 0
	reg3 := 1
	//reg3_loop
	for {
		reg5 := 1
		for {
			if reg5*reg3 == REG1 {
				fmt.Printf("pair %d %d\n", reg5, reg3)
				reg0 += reg3
			}
			reg5++
			if reg5 > REG1 {
				break
			}
		}
		reg3++
		if reg3 > REG1 {
			break
		}
	}
	fmt.Printf("%d\n", reg0)
}
