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

type Node struct {
	Child    []*Node
	Metadata []int
}

func parseTree(in []int) (*Node, []int) {
	r := &Node{}
	nchild := in[0]
	in = in[1:]
	nmeta := in[0]
	in = in[1:]
	for i := 0; i < nchild; i++ {
		var child *Node
		child, in = parseTree(in)
		r.Child = append(r.Child, child)
	}
	for i := 0; i < nmeta; i++ {
		r.Metadata = append(r.Metadata, in[0])
		in = in[1:]
	}
	return r, in
}

func (node *Node) SumMeta() int {
	r := 0
	for _, child := range node.Child {
		r += child.SumMeta()
	}
	for _, md := range node.Metadata {
		r += md
	}
	return r
}

func main() {
	buf, err := ioutil.ReadFile("08.txt")
	must(err)
	in := vatoi(splitandclean(strings.TrimSpace(string(buf)), " ", -1))
	root, _ := parseTree(in)
	fmt.Printf("%d\n", root.SumMeta())
}
