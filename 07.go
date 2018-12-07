package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
	"os"
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
			fmt.Printf("%s is ready\n", node)
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
	if len(areready) == 0 {
		return ""
	}
	sort.Strings(areready)
	return areready[0]
}

func finishnode(node string) {
	fmt.Printf("%d %s finished\n", clock, node)
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
	findready()
}

const part1 = false

var workers []int
var workerjob []string
var clock int

func freeworker() int {
	for i := range workerjob {
		if workerjob[i] == "" {
			return i
		}
	}
	return -1
}

func runwork() {
	fmt.Printf("runwork\n")
	mini := -1
	for i := range workers {
		if workerjob[i] == "" {
			continue
		}
		if mini < 0 {
			mini = i
		}
		if workers[i] < workers[mini] {
			mini = i
		}
	}
	n := workers[mini]
	fmt.Printf("working for %d units of time\n", n)
	clock += n
	for i := range workers {
		if workerjob[i] == "" {
			continue
		}
		workers[i] -= n
		if workers[i] == 0 {
			finishnode(workerjob[i])
			workerjob[i] = ""
		}
	}
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
	
	r := []string{}
	
	if part1 {
		fmt.Printf("%v\n", depends)
		findready()
		
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
	
	if part1 {
		os.Exit(0)
	}
	
	
	var numworkers = 2
	var extratime = 0
	
	if len(depends) > 7 {
		numworkers = 5
		extratime = 60
	}
	
	workers = make([]int, numworkers)
	workerjob = make([]string, numworkers)
	
	fmt.Printf("starting\n")
	fmt.Printf("%v\n", depends)
	findready()
	
	for len(ready) > 0 {
		var i int
		for {
			fmt.Printf("workerloop\n")
			i = freeworker()
			if i >= 0 {
				break
			}
			runwork()
		}
		var node string
		for {
			fmt.Printf("nodeloop\n")
			node = minready()
			if node != "" {
				break
			}
			runwork()
		}
		fmt.Printf("%d starting %s on worker %d\n", clock, node, i)
		delete(ready, node)
		workerjob[i] = node
		workers[i] = (int(node[0]) - 'A') + 1 + extratime
	}
	
	for {
		busy := false
		for i := range workers {
			if workerjob[i] != "" {
				busy = true
				break
			}
		}
		if !busy {
			break
		}
		runwork()
	}
	fmt.Printf("clock %d\n", clock)
}
