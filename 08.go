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

func (node *Node) Value() int {
	if len(node.Child) == 0 {
		r := 0
		for _, md := range node.Metadata {
			r += md
		}
		return r
	}

	r := 0
	for _, md := range node.Metadata {
		if md-1 >= len(node.Child) {
			continue
		}
		r += node.Child[md-1].Value()
	}
	return r
}

func main() {
	buf, err := ioutil.ReadFile("08.txt")
	must(err)
	in := vatoi(splitandclean(strings.TrimSpace(string(buf)), " ", -1))
	root, _ := parseTree(in)
	fmt.Printf("PART1: %d\n", root.SumMeta())
	fmt.Printf("PART2: %d\n", root.Value())
}
