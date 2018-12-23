package main

import (
	"fmt"
	"io/ioutil"
	_ "os"
	"sort"
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

var maxgroup []Ptc
var L = 416

func enumerate(start int, used map[int]bool, possible []bool, possiblecnt int, group []Ptc) {
	if len(group) > len(maxgroup) {
		maxgroup = make([]Ptc, len(group))
		copy(maxgroup, group)
		fmt.Printf("new max overlapping group found %d (possibles %d)\n", len(maxgroup), possiblecnt)
		if len(maxgroup) > L {
			L = len(maxgroup)
		}
	}
	if possiblecnt == 0 {
		return
	}
	//fmt.Printf("%d starting at %d\n", len(group), start)
	for i := start; i < len(ptcs); i++ {
		if used[i] {
			continue
		}
		if !possible[i] {
			continue
		}
		if len(group)+overlapcnt[ptcs[i]] < L {
			continue
		}

		if groupOverlap(group, ptcs[i]) {
			used[i] = true

			newpossible := make([]bool, len(ptcs))
			newpossiblecnt := 0

			for j := range ptcs {
				newpossible[j] = !used[j] && possible[j] && overlaps[Overlap{ptcs[i], ptcs[j]}]
				if newpossible[j] {
					newpossiblecnt++
				}
			}

			enumerate(i+1, used, newpossible, newpossiblecnt, append(group, ptcs[i]))

			used[i] = false
		}
	}
}

type Overlap struct {
	a, b Ptc
}

var overlaps = make(map[Overlap]bool)
var overlapcnt = make(map[Ptc]int)

func groupOverlap(group []Ptc, p Ptc) bool {
	for i := range group {
		if !overlaps[Overlap{group[i], p}] {
			return false
		}
	}
	for _, pp := range p.extremes() {
		ok := true
		for i := range group {
			if !group[i].contains(pp) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
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

	sort.Slice(ptcs, func(i, j int) bool { return ptcs[i].r > ptcs[j].r })

	for i := range ptcs {
		ocnt := 0
		for j := i + 1; j < len(ptcs); j++ {
			overlap := false
			for _, p := range ptcs[j].extremes() {
				if ptcs[i].contains(p) {
					overlap = true
					break
				}
			}
			if overlap {
				overlaps[Overlap{ptcs[i], ptcs[j]}] = true
				ocnt++
			}
		}
		overlapcnt[ptcs[i]] = ocnt
	}

	possible := make([]bool, len(ptcs))
	for i := range possible {
		possible[i] = true
	}

	enumerate(0, make(map[int]bool), possible, len(possible), nil)

	fmt.Printf("%v\n", maxgroup)
}
