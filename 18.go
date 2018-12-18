package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
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

// convert vector of strings to integer, discard non-ints
func vatoiSkip(in []string) []int {
	r := make([]int, 0, len(in))
	for i := range in {
		n, err := strconv.Atoi(in[i])
		if err == nil {
			r = append(r, n)
		}
	}
	return r
}

func getints(in string, hasneg bool) []int {
	v := getnums(in, hasneg, false)
	return vatoi(v)
}

func getfloats(in string, hasneg bool) []float64 {
	v := getnums(in, hasneg, false)
	r := make([]float64, len(v))
	for i := range v {
		var err error
		r[i], err = strconv.ParseFloat(v[i], 64)
		must(err)
	}
	return r
}

func getnums(in string, hasneg, hasdot bool) []string {
	r := []string{}
	start := -1

	flush := func(end int) {
		if start < 0 {
			return
		}
		hasdigit := false
		for i := start; i < end; i++ {
			if in[i] >= '0' && in[i] <= '9' {
				hasdigit = true
				break
			}
		}
		if hasdigit {
			r = append(r, in[start:end])
		}
		start = -1
	}

	for i, ch := range in {
		isnumch := false

		switch {
		case hasneg && (ch == '-'):
			isnumch = true
		case hasdot && (ch == '.'):
			isnumch = true
		case ch >= '0' && ch <= '9':
			isnumch = true
		}

		if start >= 0 {
			if !isnumch {
				flush(i)
			}
		} else {
			if isnumch {
				start = i
			}
		}
	}
	flush(len(in))
	return r
}

func printmatrix(matrix [][]byte) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%c", matrix[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

var M [][]byte
var AllM [][][]byte

func countAdjacent(i, j int, ch byte) int {
	r := 0
	c := func(i, j int) {
		if i < 0 || i >= len(M) {
			return
		}
		if j < 0 || j >= len(M[i]) {
			return
		}
		if M[i][j] == ch {
			r++
		}
	}
	c(i-1, j-1)
	c(i-1, j)
	c(i-1, j+1)
	c(i, j-1)
	c(i, j+1)
	c(i+1, j-1)
	c(i+1, j)
	c(i+1, j+1)
	return r
}

func step() {
	AllM = append(AllM, M)
	NM := make([][]byte, len(M))

	for i := range NM {
		NM[i] = make([]byte, len(M[i]))

		for j := range NM[i] {
			switch M[i][j] {
			case '.':
				if countAdjacent(i, j, '|') >= 3 {
					NM[i][j] = '|'
				} else {
					NM[i][j] = '.'
				}
			case '|':
				if countAdjacent(i, j, '#') >= 3 {
					NM[i][j] = '#'
				} else {
					NM[i][j] = '|'
				}
			case '#':
				if countAdjacent(i, j, '#') >= 1 && countAdjacent(i, j, '|') >= 1 {
					NM[i][j] = '#'
				} else {
					NM[i][j] = '.'
				}
			default:
				panic("blah")
			}
		}
	}
	M = NM
}

func count() int {
	trees := 0
	lumber := 0
	for i := range M {
		for j := range M[i] {
			switch M[i][j] {
			case '|':
				trees++
			case '#':
				lumber++
			}
		}
	}
	fmt.Printf("%d*%d\n", trees, lumber)
	return trees * lumber
}

func equals(M2 [][]byte) bool {
	for i := range M {
		for j := range M[i] {
			if M[i][j] != M2[i][j] {
				return false
			}
		}
	}
	return true
}

var cycleStart int

func iscycle() bool {
	for i := range AllM {
		if equals(AllM[i]) {
			fmt.Printf("cycle detected to %d\n", i)
			cycleStart = i
			return true
		}
	}
	return false
}

func main() {
	buf, err := ioutil.ReadFile("18.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	//const N = 10 // part 1
	const N = 1000000000 // part 2
	//const N = 600

	var firstCycleIdx int

	for i := 0; i < N; i++ {
		fmt.Printf("cycle %d\n", i)
		printmatrix(M)
		step()
		if iscycle() {
			AllM = AllM[:len(AllM)-1]
			firstCycleIdx = i + 1
			break
		}
	}

	for i := cycleStart; i < len(AllM); i++ {
		M = AllM[i]
		fmt.Printf("%d: %d\n", i, count())
	}

	p := cycleStart
	for i := 0; i+firstCycleIdx < N; i++ {
		p++
		if p > len(AllM) {
			p = cycleStart
		}
	}

	fmt.Printf("p is %d (N - firstCycleIdx = %d)\n", p, N-firstCycleIdx)

	//M = AllM[((N - firstCycleIdx) % (len(AllM) - cycleStart)) + cycleStart]
	M = AllM[p]
	fmt.Printf("PART2 %d\n", count())

	/*
		printmatrix(M)

		fmt.Printf("PART1 %d\n", count())
	*/
}

// N=600 -> 574*340 = 195160
