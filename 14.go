package main

import (
	"fmt"
	"os"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func todig(b []byte) []int {
	r := make([]int, len(b))
	for i := range b {
		r[i] = int(b[i] - '0')
	}
	return r
}

const inputNum = 306281
//const inputNum = 2018
var pattern = todig([]byte("306281"))
//var pattern = todig([]byte("01245"))
const debug = false
const part1 = false

var recipes = []int{ 3, 7 }
var elf1 = 0
var elf2 = 1

func checkpat(in []int) bool {
	if len(in) != len(pattern) {
		panic("wtf")
	}
	for i := range in {
		if in[i] != pattern[i] {
			return false
		}
	}
	return true
}

func runstep() {
	for _, ch := range fmt.Sprintf("%d", recipes[elf1] + recipes[elf2]) {
		recipes = append(recipes, int(ch-'0'))
		if !part1 {
			if len(recipes) >= len(pattern) && checkpat(recipes[len(recipes)-len(pattern):]) {
				fmt.Printf("found after %d recipes\n",len(recipes)-len(pattern))
				os.Exit(1)
			}
		}
	}
	elf1 = (elf1 + recipes[elf1] + 1) % len(recipes)
	elf2 = (elf2 + recipes[elf2] + 1) % len(recipes)
}

func main() {
	fmt.Printf("hello\n")
	for step := 0; ; step++ {
		if step % 10000 == 0 {
			fmt.Printf("step %d\n", step)
		}
		if debug {
			fmt.Printf("%v elf1=%d elf2=%d\n", recipes, elf1, elf2)
		}
		runstep()
		if part1 {
			if len(recipes) > inputNum + 10 {
				break
			}
		} else {
			if len(recipes) >= len(pattern) && checkpat(recipes[len(recipes)-len(pattern):]) {
				fmt.Printf("found after %d recipes\n",len(recipes)-len(pattern))
				return
			}
		}
	}
	
	if debug {
		fmt.Printf("%v elf1=%d elf2=%d\n", recipes, elf1, elf2)
	}
	fmt.Printf("OUT: ")
	for _, d := range recipes[inputNum:inputNum+10] {
		fmt.Printf("%d", d)
	}
	fmt.Printf("\n")
}
