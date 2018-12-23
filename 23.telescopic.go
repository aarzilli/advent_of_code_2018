package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func maxminof(max, min *int, cur int) {
	if cur > *max {
		*max = cur
	}
	if cur < *min {
		*min = cur
	}
}

func ptcsCopy(ptcs []Ptc) []Ptc {
	r := make([]Ptc, len(ptcs))
	copy(r, ptcs)
	return r
}

func findoverlappt(ptcs []Ptc, startMin, startMax *Coord, factor int) (int, []Coord) {
	ptcs = ptcsCopy(ptcs)

	if factor != 1 {
		for i := range ptcs {
			ptcs[i].pos.x = ptcs[i].pos.x / factor
			ptcs[i].pos.y = ptcs[i].pos.y / factor
			ptcs[i].pos.z = ptcs[i].pos.z / factor
			ptcs[i].r = (ptcs[i].r / factor) + 1
		}
	}

	var min, max Coord

	if startMin == nil {
		min = ptcs[0].pos
		max = ptcs[0].pos

		for i := range ptcs {
			maxminof(&max.x, &min.x, ptcs[i].pos.x)
			maxminof(&max.y, &min.y, ptcs[i].pos.y)
			maxminof(&max.z, &min.z, ptcs[i].pos.z)
		}

		fmt.Printf("%v %v\n", min, max)
	} else {
		min = *startMin
		max = *startMax
	}

	var maxoverlaps = 0
	var maxoverlapPos []Coord

	for x := min.x; x <= max.x; x++ {
		for y := min.y; y <= max.y; y++ {
			for z := min.z; z <= max.z; z++ {
				overlaps := 0
				overlapBitmap := make([]bool, len(ptcs))
				for i := range ptcs {
					if ptcs[i].contains(Coord{x, y, z}) {
						overlaps++
						overlapBitmap[i] = true
					}
				}
				if overlaps > maxoverlaps {
					maxoverlaps = overlaps
					maxoverlapPos = []Coord{Coord{x, y, z}}
				} else if overlaps == maxoverlaps {
					maxoverlapPos = append(maxoverlapPos, Coord{x, y, z})
				}
			}
		}
	}

	//fmt.Printf("maximum is %d particles at %v\n", maxoverlaps, maxoverlapPos)

	return maxoverlaps, maxoverlapPos
}

func closest2zero(candidates map[Coord]bool) Coord {
	first := true
	var closestp Coord
	for p := range candidates {
		if first {
			first = false
			closestp = p
		} else {
			if dist(Coord{0, 0, 0}, p) < dist(Coord{0, 0, 0}, closestp) {
				closestp = p
			}
		}
	}
	return closestp
}

func fusecandidates(candidates map[Coord]bool) (ps, pe Coord) {
	var min, max Coord
	first := true

	for p := range candidates {
		p.x *= 10
		p.y *= 10
		p.z *= 10
		pe := p
		pe.x += 9
		pe.y += 9
		pe.z += 9

		if first {
			min = p
			max = pe
			first = false
		} else {
			maxminof(&max.x, &min.x, p.x)
			maxminof(&max.y, &min.y, p.y)
			maxminof(&max.z, &min.z, p.z)
			maxminof(&max.x, &min.x, pe.x)
			maxminof(&max.y, &min.y, pe.y)
			maxminof(&max.z, &min.z, pe.z)
		}
	}

	return min, max
}

func main() {
	var ptcs = []Ptc{}
	
	path := "23.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	buf, err := ioutil.ReadFile(path)
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
	// maximum is 956 particles at [{2 4 4}]

	candidates := make(map[Coord]bool)
	factor := 10000000
	_, firstcndv := findoverlappt(ptcs, nil, nil, factor)
	for _, pos := range firstcndv {
		candidates[pos] = true
	}

	for factor != 1 {
		overlapsz := 0
		newcandidates := make(map[Coord]bool)
		fmt.Printf("considering %d macropoints %d\n", len(candidates), factor/10)
		//for p := range candidates {
		p := closest2zero(candidates)
		//p, pe := fusecandidates(candidates)
		p.x *= 10
		p.y *= 10
		p.z *= 10
		pe := p
		pe.x += 10
		pe.y += 10
		pe.z += 10
		p.x -= 10
		p.y -= 10
		p.z -= 10
		fmt.Printf("testing %v %v (%d) (number of points: %d)\n", p, pe, factor, (pe.x-p.x)*(pe.y-p.y)*(pe.z-p.z))

		sz, cndv := findoverlappt(ptcs, &p, &pe, factor/10)

		if sz > overlapsz {
			overlapsz = sz
			newcandidates = make(map[Coord]bool)
		}

		for _, np := range cndv {
			newcandidates[np] = true
		}
		//}
		candidates = newcandidates
		cp := closest2zero(candidates)
		fmt.Printf("%v %d (overlap size: %d)\n", cp, dist(Coord{0, 0, 0}, cp), overlapsz)
		factor = factor / 10
	}

	fmt.Printf("%d macropoints remain\n", len(candidates))
	fmt.Printf("%v\n", closest2zero(candidates))

	/*
		findoverlappt(ptcs, &Coord{ 20, 40, 40 }, &Coord{ 29, 49, 49}, 1000000)
		// maximum is 901 particles at [{20 48 43}]

		findoverlappt(ptcs, &Coord{ 200, 480, 430 }, &Coord{ 209, 489, 439}, 100000)
		// maximum is 871 particles at [{200 485 430} {200 486 431} {200 487 432} {200 488 433} {200 489 434} {201 486 430} {201 487 431} {201 488 432} {201 489 433} {202 487 430} {202 488 431} {202 489 432}]

		for _, p := range []Coord{ {200, 485, 430}, {200, 486, 431}, {200, 487, 432}, {200, 488, 433}, {200, 489, 434}, {201, 486, 430}, {201, 487, 431}, {201, 488, 432}, {201, 489, 433}, {202, 487, 430}, {202, 488, 431}, {202, 489, 432} } {
			p.x *= 10
			p.y *= 10
			p.z *= 10
			pe := p
			pe.x += 9
			pe.y += 9
			pe.z += 9
			//fmt.Printf("testing %v %v\n", p, pe)
			//findoverlappt(ptcs, &p, &pe, 10000)
		}

	*/
	/*
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
				ovps = []Coord{ p }
			case cnt == maxov:
				ovps = append(ovps, p)
			}
		}

		fmt.Printf("%d %v\n", maxov, ovps)

		fmt.Printf("%d\n", dist(Coord{0,0,0}, ovps[0]))*/
}

// 110690993 too low
// 111513514 too low
// 112997634 fuck me
