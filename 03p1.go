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

type Claim struct {
	id   int
	i, j int
	w, h int
}

var claims = []Claim{}
var M [][]int

func main() {
	buf, err := ioutil.ReadFile("03.txt")
	must(err)
	maxi, maxj := 0, 0
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		id := atoi(fields[0][1:])
		c := splitandclean(nolast(fields[2]), ",", 2)
		j := atoi(c[0])
		i := atoi(c[1])
		d := splitandclean(fields[3], "x", 2)
		w := atoi(d[0])
		h := atoi(d[1])
		claims = append(claims, Claim{
			id, i, j, w, h,
		})
		if i+h+1 > maxi {
			maxi = i + h + 1
		}
		if j+w+1 > maxj {
			maxj = j + w + 1
		}
	}

	//fmt.Printf("%v\n", claims)

	M = make([][]int, maxi*2)
	for i := range M {
		M[i] = make([]int, maxj*2)
	}

	for _, claim := range claims {
		lay(claim)
	}

	/*
		for i := range M {
			for j := range M[i] {
				if M[i][j] > 0 {
					fmt.Printf("%d", M[i][j])
				} else if M[i][j] < 0 {
					fmt.Printf("X")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
	*/

	contested := 0
	for i := range M {
		for j := range M[i] {
			if M[i][j] < 0 {
				contested++
			}
		}
	}
	fmt.Printf("%d\n", contested)
}

func lay(claim Claim) {
	for i := 0; i < claim.h; i++ {
		for j := 0; j < claim.w; j++ {
			if M[i+claim.i][j+claim.j] != 0 {
				M[i+claim.i][j+claim.j] = -1
			} else {
				M[i+claim.i][j+claim.j] = claim.id
			}
		}
	}
}
