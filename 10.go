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

func parsepoint(in string) (int, int) {
	v := splitandclean(in, "<", -1)
	vv := splitandclean(v[1], ",", -1)
	return atoi(vv[0]), atoi(nolast(vv[1]))
}

type Coord struct {
	x, y int
}

type Ptc struct {
	pos, v Coord
}

var ptcs = []Ptc{}
var clock = 0

func display() {
	var min, max Coord
	first := true
	M := map[Coord]bool{}
	for _, ptc := range ptcs {
		if first {
			first = false
			min = ptc.pos
			max = ptc.pos
		}
		if ptc.pos.x < min.x {
			min.x = ptc.pos.x
		}
		if ptc.pos.x > max.x {
			max.x = ptc.pos.x
		}
		if ptc.pos.y < min.y {
			min.y = ptc.pos.y
		}
		if ptc.pos.y > max.y {
			max.y = ptc.pos.y
		}
		M[ptc.pos] = true
	}
	
	pixels := (max.y - min.y) * (max.x - min.x)
	
	fmt.Printf("PIXELS %d\n", pixels)
	
	
	if pixels < 5000 {
		for y := min.y; y <= max.y; y++ {
			for x := min.x; x <= max.x; x++ {
				if M[Coord{ x, y }] {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Println()
		fmt.Println()
	}
}

func simstep() {
	for i := range ptcs {
		ptc := &ptcs[i]
		ptc.pos.x += ptc.v.x
		ptc.pos.y += ptc.v.y
	}
	clock++
}

func main() {
	buf, err := ioutil.ReadFile("10.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fuck := splitandclean(line, " velocity=", 2)
		x, y := parsepoint(fuck[0])
		vx, vy := parsepoint(fuck[1])
		
		ptcs = append(ptcs, Ptc{ pos: Coord{ x, y }, v: Coord{ vx, vy }  })
	}
	
	skip := 10
	
	for {
		if skip <= 0 {
			fmt.Printf("NEXT %d\n", clock)
			display()
		} else {
			skip--
		}
		simstep()
	}
}

// EHAZPZHP