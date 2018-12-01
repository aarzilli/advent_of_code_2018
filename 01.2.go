package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

func main() {
	buf, err := ioutil.ReadFile("01.txt")
	must(err)
	freqs := []int{}
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		freqs = append(freqs, atoi(line))
	}
	
	freq := 0
	seen := map[int]bool{}
	i := 0
	for {
		if seen[freq] {
			fmt.Printf("seen again %d\n", freq)
			break
		}
		seen[freq] = true
		freq += freqs[i % len(freqs)]
		i++
	}
}

// rank 103
