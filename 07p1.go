package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
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

var depends = map[string][]string{}
var ready = map[string]bool{}

func findready() {
	for node := range ready {
		if len(depends[node]) == 0 {
			ready[node] = true
			//fmt.Printf("%s is ready\n", node)
		}
	}
}

func minready() string {
	areready := []string{}
	for node := range ready {
		if ready[node] {
			areready = append(areready, node)
		}
	}
	sort.Strings(areready)
	return areready[0]
}

func main() {
	fmt.Printf("hello\n")
	buf, err := ioutil.ReadFile("07.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		depends[fields[7]] = append(depends[fields[7]], fields[1])
		ready[fields[7]] = false
		ready[fields[1]] = false
	}
	
	fmt.Printf("%v\n", depends)
	
	findready()
	
	r := []string{}
	
	for len(ready) > 0 {
		node := minready()
		//fmt.Printf("processing %s\n", node)
		r = append(r, node)
		delete(ready, node)
		for node2 := range depends {
			if r, ok := ready[node2]; !ok || r {
				continue
			}
			newdep := []string{}
			for _, dep := range depends[node2] {
				if dep != node {
					newdep = append(newdep, dep)
				}
			}
			depends[node2] = newdep
		}
		//fmt.Printf("\n")
		findready()
	}
	
	fmt.Printf("%s\n", strings.Join(r, ""))
}
