package main

import (
	"fmt"
	"io/ioutil"
	_ "os"
	"strconv"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
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

func getints(in string, hasneg bool) []int {
	v := getnums(in, hasneg, false)
	return vatoi(v)
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
	pos Coord
	r   int
}

type Coord struct {
	x, y, z int
}

var ptcs = []Ptc{}

func dist(a, b Coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p Ptc) contains(q Coord) bool {
	return dist(p.pos, q) <= p.r
}

func (p Ptc) extremes() []Coord {
	return []Coord{
		Coord{p.pos.x - p.r, p.pos.y, p.pos.z},
		Coord{p.pos.x + p.r, p.pos.y, p.pos.z},
		Coord{p.pos.x, p.pos.y - p.r, p.pos.z},
		Coord{p.pos.x, p.pos.y + p.r, p.pos.z},
		Coord{p.pos.x, p.pos.y, p.pos.z - p.r},
		Coord{p.pos.x, p.pos.y, p.pos.z + p.r},
	}
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
		ptcs = append(ptcs, Ptc{pos: Coord{x: posv[0], y: posv[1], z: posv[2]}, r: radiumv[0]})
	}
	maxi := 0
	for i := range ptcs {
		if ptcs[i].r > ptcs[maxi].r {
			maxi = i
		}
	}
	cnt := 0
	for i := range ptcs {
		if ptcs[maxi].contains(ptcs[i].pos) {
			cnt++
		}
	}
	fmt.Printf("PART1 %d\n", cnt)

	// PART2

	ps := []Coord{}

	for i := range ptcs {
		ps = append(ps, ptcs[i].extremes()...)
		ps = append(ps, ptcs[i].pos)
	}

	maxov := 0
	ovps := []Coord{}

	for _, p := range ps {
		cnt := 0
		for i := range ptcs {
			if ptcs[i].contains(p) {
				cnt++
			}
		}
		switch {
		case cnt > maxov:
			fmt.Printf("for %d %v\n", maxov, ovps)
			maxov = cnt
			ovps = []Coord{p}
		case cnt == maxov:
			ovps = append(ovps, p)
		}
	}

	fmt.Printf("%d %v\n", maxov, ovps)

	fmt.Printf("%d\n", dist(Coord{0, 0, 0}, ovps[0]))
}

// 110690993 too low
