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

var erosionLvl = map[Coord]int{}
var M = map[Coord]byte{}

func geoidx(p Coord) int {
	switch {
	case p == Coord{0, 0}:
		return 0
	case p == target:
		return 0
	case p.y == 0:
		return p.x * 16807
	case p.x == 0:
		return p.y * 48271
	default:
		return erosionLvl[Coord{p.x - 1, p.y}] * erosionLvl[Coord{p.x, p.y - 1}]
	}
}

func getErosionLvl(p Coord) int {
	if e, ok := erosionLvl[p]; ok {
		return e
	}
	erosionLvl[p] = (geoidx(p) + depth) % 20183
	switch erosionLvl[p] % 3 {
	case 0:
		M[p] = '.'
	case 1:
		M[p] = '='
	case 2:
		M[p] = '|'
	}
	if p == (Coord{0, 0}) {
		M[p] = 'M'
	}
	if p == target {
		M[p] = 'T'
	}
	return erosionLvl[p]
}

// example
/*
var target Coord = Coord{ 10, 10 }
var depth int = 510
*/

// real input
var target = Coord{12, 763}
var depth int = 7740

func example(x, y int) {
	p := Coord{x, y}
	fmt.Printf("%d %d\n", geoidx(p), getErosionLvl(p))
	fmt.Printf("%c\n", M[p])
}

func main() {
	/*
		for y := 0; y <= 15; y++ {
			for x := 0; x <= 15; x++ {
				getErosionLvl(Coord{ x, y })
				fmt.Printf("%c", M[Coord{ x, y }])
			}
			fmt.Printf("\n")
		}
	*/

	risk := 0
	for y := 0; y <= target.y; y++ {
		for x := 0; x <= target.x; x++ {
			getErosionLvl(Coord{x, y})
			switch M[Coord{x, y}] {
			case '.':
				risk += 0
			case '=':
				risk += 1
			case '|':
				risk += 2
			}
		}
	}

	fmt.Printf("PART1 %d\n", risk)
}
