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

func getints(in string, hasneg bool) []int {
	v := getnums(in, hasneg, false)
	return vatoi(v)
}

func getfloats(in string, hasneg bool) []float64 {
	v := getnums(in, hasneg, false)
	r := make([]float64, len(v))
	for i := range v {
		var err error
		r[i], err = strconv.ParseFloat(v[i], 64)
		must(err)
	}
	return r
}

func getnums(in string, hasneg, hasdot bool) []string {
	r := []string{}
	start := -1

	flush := func(end int) {
		if start < 0 {
			return
		}
		hasdigit := false
		for i := start; i < end; i++ {
			if in[i] >= '0' && in[i] <= '9' {
				hasdigit = true
				break
			}
		}
		if hasdigit {
			r = append(r, in[start:end])
		}
		start = -1
	}

	for i, ch := range in {
		isnumch := false

		switch {
		case hasneg && (ch == '-'):
			isnumch = true
		case hasdot && (ch == '.'):
			isnumch = true
		case ch >= '0' && ch <= '9':
			isnumch = true
		}

		if start >= 0 {
			if !isnumch {
				flush(i)
			}
		} else {
			if isnumch {
				start = i
			}
		}
	}
	flush(len(in))
	return r
}

var state []byte
var rules = map[string]byte{}
var startid = 0

type Rule struct {
	pat []byte
	out byte
}

func runrules(in []byte) byte {
	out, ok := rules[string(in)]
	if !ok {
		return '.'
	}
	return out
}

func simstep() {
	n0 := make([]byte, len(state)+8)
	n0[0], n0[1], n0[2], n0[3] = '.', '.', '.', '.'
	copy(n0[4:], state)
	n0[len(n0)-1], n0[len(n0)-2], n0[len(n0)-3], n0[len(n0)-4] = '.', '.', '.', '.'
	n1 := make([]byte, len(n0))
	startid = startid - 4

	for i := range n1 {
		if i-1 < 0 || i-2 < 0 {
			n1[i] = runrules([]byte{'.', '.', n0[i], n0[i+1], n0[i+2]})
		} else if i+1 >= len(n0) || i+2 >= len(n0) {
			n1[i] = runrules([]byte{n0[i-2], n0[i-1], n0[i], '.', '.'})
		} else {
			n1[i] = runrules([]byte{n0[i-2], n0[i-1], n0[i], n0[i+1], n0[i+2]})
		}
	}

	start := 0
	for i := range n1 {
		if n1[i] == '#' {
			start = i
			break
		}
	}
	startid = startid + start

	end := len(n1)
	for i := len(n1) - 1; i >= 0; i-- {
		if n1[i] == '#' {
			end = i + 1
			break
		}
	}

	state = n1[start:end]
}

func count() int {
	r := 0
	for i := range state {
		if state[i] == '#' {
			r += startid + i
		}
	}
	return r
}

func main() {
	buf, err := ioutil.ReadFile("12.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "initial state:") {
			v := splitandclean(line, ":", -1)
			state = []byte(v[1])
		} else {
			v := splitandclean(line, "=>", -1)
			rules[v[0]] = v[1][0]
		}
	}

	fmt.Printf("initial state: <%s>\n", string(state))
	fmt.Printf("rules %v\n", rules)

	for gen := 0; gen < 20; gen++ {
		fmt.Printf("%d %s\n", gen, string(state))
		simstep()
	}

	fmt.Printf("final state: <%s> %d\n", string(state), startid)

	fmt.Printf("%d\n", count())

}
