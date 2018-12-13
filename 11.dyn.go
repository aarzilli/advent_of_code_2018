package main

import (
	"fmt"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Coord struct {
	x, y int
}

func powerlvl(p Coord, serialnum int) int8 {
	rackID := p.x + 10
	lvlstart := rackID * p.y
	lvl := lvlstart + serialnum
	lvl *= rackID
	lvl = (lvl / 100) % 10
	return int8(lvl - 5)
}

func makegrid(serialnum int) [][]int8 {
	M := make([][]int8, 300)
	for i := range M {
		M[i] = make([]int8, 300)
	}

	for i := range M {
		for j := range M[i] {
			M[i][j] = powerlvl(Coord{j + 1, i + 1}, serialnum)
		}
	}

	return M
}

func sumn_old(M [][]int8, starti, startj int, n int) int {
	r := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			r += int(M[starti+i][startj+j])
		}
	}
	return r
}

func sum_rect(M [][]int8, SUMN [][][]int, starti, startj int, n, m int) int {
	if n > 1 && m > 1 {
		if n < m {
			r := SUMN[n][starti][startj]
			r += sum_rect(M, SUMN, starti, startj+n, n, m-n)
			return r
		} else {
			r := SUMN[m][starti][startj]
			r += sum_rect(M, SUMN, starti+m, startj, n-m, m)
			return r
		}
	}
	r := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			r += int(M[starti+i][startj+j])
		}
	}
	return r
}

func sumn_dyn(M [][]int8, SUMN [][][]int, starti, startj int, n int) int {
	if n == 1 {
		SUMN[n][starti][startj] = int(M[starti][startj])
		return SUMN[n][starti][startj]
	}

	if n%2 == 0 {
		SUMN[n][starti][startj] = SUMN[n/2][starti][startj] + SUMN[n/2][starti+n/2][startj] + SUMN[n/2][starti+n/2][startj+n/2] + SUMN[n/2][starti][startj+n/2]
		return SUMN[n][starti][startj]
	}

	m := n/2 + 1
	r := SUMN[m][starti][startj]
	r += SUMN[n-m][starti+m][startj+m]
	r += sum_rect(M, SUMN, starti+m, startj, n-m, m)
	r += sum_rect(M, SUMN, starti, startj+m, m, n-m)
	SUMN[n][starti][startj] = r

	return SUMN[n][starti][startj]
}

func sumn(M [][]int8, SUMN [][][]int, starti, startj int, n int) int {
	r := sumn_dyn(M, SUMN, starti, startj, n)
	/*r1 := sumn_old(M, starti, startj, n)
	if r != r1 {
		panic("bug found")
	}*/
	return r
}

func findmaxpt2(M [][]int8) (Coord, int, int) {
	SUMN := make([][][]int, len(M)+1)
	for n := range SUMN {
		SUMN[n] = make([][]int, len(M))
		for i := range SUMN[n] {
			SUMN[n][i] = make([]int, len(M))
		}
	}

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
				val := sumn(M, SUMN, i, j, n)
				if val > maxval || first {
					first = false
					maxval = val
					max = Coord{j + 1, i + 1}
					maxn = n
				}
			}
		}
	}
	return max, maxn, maxval
}

func main() {
	M := makegrid(3031)
	p, n, pow := findmaxpt2(M)
	fmt.Printf("%d,%d,%d (%d)\n", p.x, p.y, n, pow)
}

// 234,108,16 (160)

/*
O(n**5) solution:
	real	1m42,053s
	user	1m42,075s
	sys	0m0,131s

Dynamic programming (without rectangles optimization):
	real	0m27,171s
	user	0m27,159s
	sys	0m0,160s

Better dynamic programming:
	real	0m1,201s
	user	0m1,150s
	sys	0m0,177s
*/
