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

type Coord struct {
	x, y, z, t int
}

var coords []Coord
var clique = map[Coord]int{}

var adj = map[Coord]map[Coord]bool{}

func dist(a, b Coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z) + abs(a.t-b.t)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func color(cur Coord, curcolor int) {
	clique[cur] = curcolor
	for next := range adj[cur] {
		if _, ok := clique[next]; ok {
			continue
		}
		color(next, curcolor)
	}
}

func main() {
	buf, err := ioutil.ReadFile("25.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		f := vatoi(splitandclean(line, ",", -1))
		coords = append(coords, Coord{
			x: f[0],
			y: f[1],
			z: f[2],
			t: f[3],
		})
	}
	//fmt.Printf("%v\n", coords)

	for i := range coords {
		for j := i + 1; j < len(coords); j++ {
			if dist(coords[i], coords[j]) <= 3 {
				if adj[coords[i]] == nil {
					adj[coords[i]] = make(map[Coord]bool)
				}
				adj[coords[i]][coords[j]] = true
				if adj[coords[j]] == nil {
					adj[coords[j]] = make(map[Coord]bool)
				}
				adj[coords[j]][coords[i]] = true
			}
		}
	}

	curcolor := 1

	for {
		didsomething := false
		for i := range coords {
			if clique[coords[i]] == 0 {
				color(coords[i], curcolor)
				curcolor++
				didsomething = true
			}
		}
		if !didsomething {
			break
		}
	}

	fmt.Printf("PART1: %d\n", curcolor-1)
}
