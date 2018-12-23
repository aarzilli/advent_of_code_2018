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

type Ptc struct {
	x, y, z int
	r       int
}

var ptcs = []Ptc{}

func dist(a, b Ptc) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	buf, err := ioutil.ReadFile("23.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := splitandclean(line, ", r=", -1)
		posv := getints(v[0], true)
		radiumv := getints(v[1], true)
		ptcs = append(ptcs, Ptc{x: posv[0], y: posv[1], z: posv[2], r: radiumv[0]})
	}
	maxi := 0
	for i := range ptcs {
		if ptcs[i].r > ptcs[maxi].r {
			maxi = i
		}
	}
	cnt := 0
	for i := range ptcs {
		if dist(ptcs[i], ptcs[maxi]) <= ptcs[maxi].r {
			cnt++
		}
	}
	fmt.Printf("PART1 %d\n", cnt)
}
