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

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("02.txt")
	must(err)
	boxes := []string{}
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		boxes = append(boxes, line)
	}
	
	for i := range boxes {
		for j := i+1; j < len(boxes); j++ {
			if diff1(boxes[i], boxes[j]) {
				fmt.Printf("%s\n%s\n", boxes[i], boxes[j])
			}
		}
	}
}

func diff1(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	
	diffs := 0
	for i := range a {
		if a[i] != b[i] {
			diffs++
		}
	}
	
	return diffs == 1
}

// rank 91

/*
srijafjzloguvlntqmphenbkd
srijafjzloguvlnctqmphenbkd
*/