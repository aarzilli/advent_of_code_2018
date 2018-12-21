package main

import "fmt"
import "os"

func main() {
	cnt := 0
	seen := make(map[int]bool)
	reg1 := 0
	reg4 := 0
	reg0 := 0
	//FUCKYLOOP
	for {
		reg5 := reg4 | 65536

		reg4 = 1765573

		// OUTERLOOP:
		for {
			reg1 = reg5 & 255
			reg4 = reg4 + reg1
			reg4 = reg4 & 16777215
			reg4 = reg4 * 65899
			reg4 = reg4 & 16777215
			if reg5 < 256 {
				break
			}

			// INNERLOOP
			reg1 = 0
			for {
				reg3 := (reg1 + 1) * 256
				if reg3 > reg5 {
					reg5 = reg1
					break
				}
				reg1++
			}

			reg5 = reg1
		}

		// ENDCHECK
		fmt.Printf("%d %d (%v)\n", cnt, reg4, seen[reg4])
		if seen[reg4] {
			return
		}
		cnt++
		seen[reg4] = true
		if reg4 == reg0 {
			os.Exit(0)
		}
	}
}
