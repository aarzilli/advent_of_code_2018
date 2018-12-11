package main

import (
	"fmt"
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

type Coord struct {
	x, y int
}


func powerlvl(p Coord, serialnum int) int {
	rackID := p.x + 10
	lvlstart := rackID * p.y
	lvl := lvlstart + serialnum
	lvl *= rackID
	lvl = (lvl/100) % 10
	return lvl-5
}

func makegrid(serialnum int) [][]int {
	M := make([][]int, 300)
	for i := range M {
		M[i] = make([]int, 300)
	}
	
	for i := range M {
		for j := range M[i] {
			M[i][j] = powerlvl(Coord{ j+1, i+1 }, serialnum)
		}
	}
	
	return M
}

func printmatrix(M [][]int, start Coord) {
	for i := start.y-2; i < start.y-1+4; i++ {
		for j := start.x-2; j < start.x-1+4; j++ {
			fmt.Printf("%+d ", M[i][j])
		}
		fmt.Printf("\n")
	}
}

func sum3(M [][]int, starti, startj int) int {
	r := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r += M[starti+i][startj+j]
		}
	}
	return r
}

func sumn(M [][]int, starti, startj int, n int) int {
	r := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			r += M[starti+i][startj+j]
		}
	}
	return r
}

func findmax(M [][]int) (Coord, int) {
	var max Coord
	var maxval int
	first := true
	for i := range M {
		if i+3 >= len(M) {
			continue
		}
		for j := range M[i] {
			if j+3 >= len(M[i]) {
				continue
			}
			val := sum3(M, i, j)
			if val > maxval || first {
				first = false
				maxval = val
				max = Coord{ j+1, i+1}
			}
		}
	}
	return max, maxval
}

func findmaxpt2(M [][]int) (Coord, int, int) {
	var max Coord
	var maxn int
	var maxval int
	first := true
	for n := 0; n < len(M); n++ { // letting this run for 2 minutes is still faster than writing the dynamic programming solution
		fmt.Printf("%d\n", n)
		for i := range M {
			if i+n >= len(M) {
				continue
			}
			for j := range M[i] {
				if j+n >= len(M[i]) {
					continue
				}
				val := sumn(M, i, j, n)
				if val > maxval || first {
					first = false
					maxval = val
					max = Coord{ j+1, i+1}
					maxn = n
				}
			}
		}
	}
	return max, maxn, maxval
}

func main() {
	/*
	fmt.Printf("%d\n", powerlvl(Coord{ 3, 5 }, 8))
	fmt.Printf("%d\n", powerlvl(Coord{ 122, 79 }, 57))
	fmt.Printf("%d\n", powerlvl(Coord{ 217, 196 }, 39))
	fmt.Printf("%d\n", powerlvl(Coord{ 101, 153 }, 71))
	*/
	
	//M := makegrid(18)
	//fmt.Printf("%v\n", findmax(M))
	//printmatrix(M, Coord{ 33, 45 })
	
	//M := makegrid(42)
	//p, pow := findmax(M)
	//fmt.Printf("%v %v\n", p, pow)
	//printmatrix(M, Coord{ 21, 61 })
	
// 	M := makegrid(3031)
// 	p, pow := findmax(M)
// 	fmt.Printf("%v %v\n", p, pow)

	M := makegrid(3031)
	p, n, pow := findmaxpt2(M)
	fmt.Printf("%d,%d,%d (%d)\n", p.x, p.y, n, pow)
}

/*
Fuel cell at  122,79, grid serial number 57: power level -5.
Fuel cell at 217,196, grid serial number 39: power level  0.
Fuel cell at 101,153, grid serial number 71: power level  4.
*/