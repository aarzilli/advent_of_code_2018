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

func main() {
	buf, err := ioutil.ReadFile("09.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		part1(atoi(fields[0]), atoi(fields[6]))
	}
}

const debug = false
const debug2 = false

func part1(nplayers, lastpoints int) {
	fmt.Printf("# %d %d\n", nplayers, lastpoints)
	circle := []int{ 0 }
	i := 0
	players := make([]int, nplayers)
	_ = players
	for marble := 1; marble <= lastpoints; marble++ {
		if marble % 23 == 0 {
			if debug2 {
				fmt.Printf("before:\n")
				printcircle(circle, i)
			}
			players[(marble-1)%len(players)] += marble
			newcircle := make([]int, 0, len(circle)-1)
			ri1 := (i-7)
			if ri1 < 0 {
				ri1 = len(circle)+ri1
			} else {
				ri1 = ri1 % len(circle)
			}
			ri2 := (i-6)
			if ri2 < 0 {
				ri2 = len(circle)+ri2
			} else {
				ri2 = ri2 % len(circle)
			}
			newcircle = append(newcircle, circle[:ri1]...)
			if debug2 {
				fmt.Printf("removing %d\n", circle[ri1])
			}
			players[(marble-1)%len(players)] += circle[ri1]
			newi := len(newcircle)
			newcircle = append(newcircle, circle[ri2:]...)
			circle = newcircle
			i = newi
			if debug2 {
				fmt.Printf("after:\n")
				printcircle(circle, i)
			}
			continue
		}
		if debug {
			printcircle(circle, i)
		}
		newcircle := make([]int, 0, len(circle)+1)
		newcircle = append(newcircle, circle[:(i+1)%len(circle)+1]...)
		newi := len(newcircle)
		newcircle = append(newcircle, marble)
		newcircle = append(newcircle, circle[(i+1)%len(circle)+1:]...)
		circle = newcircle
		i = newi
	}
	if debug {
		printcircle(circle, i)
	}
	max := 0
	for i := range players {
		if players[i] > max {
			max = players[i]
		}
	}
	fmt.Printf("highest score %d\n", max)
}

func printcircle(circle []int, cur int) {
	for i := range circle {
		if i == cur {
			fmt.Printf("(%d)", circle[i])
		} else {
			fmt.Printf(" %d ", circle[i])
		}
	}
	fmt.Printf("\n")
}