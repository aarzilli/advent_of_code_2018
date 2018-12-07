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

type Coord struct {
	x, y int
}

var M = map[Coord]int{}
var coords = []Coord{}

func findClosest(p Coord) int {
	mini := -1
	mind := 100000000
	for i, coord := range coords {
		d := dist(coord, p)
		if d < mind {
			mini = i
			mind = d
		} else if d == mind {
			mini = -1
		}
	}
	return mini
}

func dist(a, b Coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func largestArea(excl map[int]bool) int {
	area := map[int]int{}

	for _, idx := range M {
		if excl[idx] {
			continue
		}
		area[idx]++
	}

	maxi := 0
	for i := range area {
		if area[i] > area[maxi] {
			maxi = i
		}
	}

	return area[maxi]
}

const debugPart1 = false

func main() {
	buf, err := ioutil.ReadFile("06.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := splitandclean(line, ",", -1)
		coord := Coord{atoi(v[0]), atoi(v[1])}
		coords = append(coords, coord)
	}

	min, max := coords[0], coords[0]

	for _, coord := range coords {
		if coord.x < min.x {
			min.x = coord.x
		}
		if coord.y < min.y {
			min.y = coord.y
		}
		if coord.x > max.x {
			max.x = coord.x
		}
		if coord.y > max.y {
			max.y = coord.y
		}
	}

	fmt.Printf("%d\n", findClosest(Coord{5, 3}))

	for x := min.x - 2; x <= max.x+2; x++ {
		for y := min.y - 2; y <= max.y+2; y++ {
			M[Coord{x, y}] = findClosest(Coord{x, y})
		}
	}

	if debugPart1 {
		for y := min.y - 2; y <= max.y+2; y++ {
			for x := min.x - 2; x <= max.x+2; x++ {
				if M[Coord{x, y}] < 0 {
					fmt.Printf(".")
				} else {
					if dist(Coord{x, y}, coords[M[Coord{x, y}]]) == 0 {
						fmt.Printf("%c", M[Coord{x, y}]+'A')
					} else {
						fmt.Printf("%c", M[Coord{x, y}]+'a')
					}
				}
			}
			fmt.Printf("\n")
		}
	}

	excl := map[int]bool{}

	for y := min.y - 2; y <= max.y+2; y++ {
		excl[M[Coord{min.x - 2, y}]] = true
		excl[M[Coord{max.x + 2, y}]] = true
	}

	for x := min.x - 2; x <= max.x+2; x++ {
		excl[M[Coord{x, min.y - 2}]] = true
		excl[M[Coord{x, max.x + 2}]] = true
	}

	fmt.Printf("%v\n", excl)

	fmt.Printf("PART1: %d\n", largestArea(excl))
}
