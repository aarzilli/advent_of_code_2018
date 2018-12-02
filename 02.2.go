package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	buf, err := ioutil.ReadFile("02.txt")
	must(err)
	boxes := []string{}
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		boxes = append(boxes, line)
	}
	
	for i := range boxes {
		for j := i+1; j < len(boxes); j++ {
			if diff1(boxes[i], boxes[j]) {
				fmt.Printf("%s\n%s\n", boxes[i], boxes[j])
			}
		}
	}
}

func diff1(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	
	diffs := 0
	for i := range a {
		if a[i] != b[i] {
			diffs++
		}
	}
	
	return diffs == 1
}

// rank 91

/*
srijafjzloguvlntqmphenbkd
srijafjzloguvlnctqmphenbkd
*/