package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// returns x without the last character
func nolast(x string) string {
	return x[:len(x)-1]
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

// convert vector of strings to integer
func vatoi(in []string) []int {
	r := make([]int, len(in))
	for i := range in {
		var err error
		r[i], err = strconv.Atoi(in[i])
		must(err)
	}
	return r
}

func printmatrix(matrix [][]byte) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%c ", matrix[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func countpixels(matrix [][]byte) (cnt int) {
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == '#' {
				cnt++
			}
		}
	}
	return cnt
}

func reduce(in []byte) []byte {
	for {
		found := false
		for i := 0; i < len(in)-1; i++ {
			if in[i] == in[i+1]+32 || in[i]+32 == in[i+1] {
				//fmt.Printf("at %d %c %c\n", i, in[i], in[i+1])
				copy(in[i:], in[i+2:])
				in = in[:len(in)-2]
				found = true
			}
		}
		if !found {
			break
		}
	}
	return in
}

func removeAll(in []byte, ch byte) []byte {
	for i := 0; i < len(in); i++ {
		if in[i] == ch || in[i] == ch+32 {
			copy(in[i:], in[i+1:])
			in = in[:len(in)-1]
			i--
		}
	}
	return in
}

func part2(in []byte) {
	shortest := len(in)
	for ch := byte('A'); ch <= 'Z'; ch++ {
		incopy := make([]byte, len(in))
		copy(incopy, in)
		incopy = removeAll(incopy, ch)
		incopy = reduce(incopy)
		if len(incopy) < shortest {
			shortest = len(incopy)
		}
	}
	fmt.Printf("%d\n", shortest)
}

const part1 = false

func main() {
	buf, err := ioutil.ReadFile("05.txt")
	must(err)
	trueInput := buf
	_ = trueInput
	if trueInput[len(trueInput)-1] == '\n' {
		trueInput = trueInput[:len(trueInput)-1]
		fmt.Printf("fuck\n")
	}
	example := []byte("dabAcCaCBAcCcaDA")
	if part1 {
		fmt.Printf("%s\n", string(example))
		example = reduce(example)
		fmt.Printf("%s (%d)\n", example, len(example))
		trueInput = reduce(trueInput)
		fmt.Printf("%d\n", len(trueInput))
	}
	if !part1 {
		part2(example)
		part2(trueInput)
	}
}
