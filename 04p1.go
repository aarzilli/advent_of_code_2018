package main

import (
	"fmt"
	"io/ioutil"
	"sort"
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

type Shift struct {
	id    int
	key   string
	times []int
}

func (s *Shift) asleep() int {
	awake := true
	start := 0
	r := 0
	for _, t := range s.times[1:] {
		if awake {
			start = t
			awake = false
		} else {
			r += t - start
			awake = true
		}
	}
	return r
}

func (s *Shift) isAsleep(minute int) bool {
	awake := true
	for _, t := range s.times[1:] {
		if minute < t {
			return !awake
		}
		awake = !awake
	}
	return !awake
}

var monthlen = []int{
	0,  //nn
	31, //gen
	28, //feb
	31, //mar
	30, //apr
	31, //mag
	30, //giu
	31, //lug
	31, //ago
	30, //set
	31, //ott
	30, //nov
	31, //div
}

func getkey(line string) string {
	fields := splitandclean(line, " ", -1)
	key := fields[0][1:]

	hour := splitandclean(fields[1], ":", -1)[0]
	if hour == "00" {
		return key
	}

	if hour != "23" {
		panic("wtf")
	}

	v := vatoi(splitandclean(key, "-", -1))

	v[2]++

	if v[2] > monthlen[v[1]] {
		v[2] = 1
		v[1]++
	}

	if v[1] > 12 {
		v[1] = 1
		v[0]++
	}

	return fmt.Sprintf("%04d-%02d-%02d", v[0], v[1], v[2])
}

func main() {
	buf, err := ioutil.ReadFile("04.txt")
	must(err)

	lines := strings.Split(string(buf), "\n")
	sort.Strings(lines)

	shifts := []Shift{}
	var curshift *Shift

	awake := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fmt.Printf("current line: %q\n", line)

		shouldAdd := false

		if strings.Contains(line, "begins shift") {
			id := atoi(splitandclean(line, " ", -1)[3][1:])
			shifts = append(shifts, Shift{id: id})
			curshift = &shifts[len(shifts)-1]
			awake = true
			shouldAdd = true
			curshift.key = getkey(line)
			fmt.Printf("\tshift start %q\n", curshift.key)
		}

		if !shouldAdd {
			if strings.Contains(line, "wakes up") && !awake {
				shouldAdd = true
			} else if strings.Contains(line, "falls asleep") && awake {
				shouldAdd = true
			}
			if shouldAdd {
				fmt.Printf("\tawake: %v -> %v\n", awake, !awake)
				awake = !awake
			}
		}

		if !shouldAdd {
			continue
		}

		v := vatoi(splitandclean(nolast(splitandclean(line, " ", -1)[1]), ":", 2))
		minute := v[0]*60 + v[1]
		if minute > 60 {
			minute = minute - (24 * 60)
		}

		if curshift != nil && getkey(line) == curshift.key {
			fmt.Printf("\tappending minute %d\n", minute)
			curshift.times = append(curshift.times, minute)
		}
	}

	fmt.Printf("%#v\n", shifts)

	asleep := map[int]int{}
	for i := range shifts {
		n := shifts[i].asleep()
		asleep[shifts[i].id] += n
	}

	fmt.Printf("%#v\n", asleep)

	sleepiestguard := 0
	for i := range asleep {
		if asleep[i] > asleep[sleepiestguard] {
			sleepiestguard = i
		}
	}

	fmt.Printf("sleepiest guard %d\n", sleepiestguard)

	asleepcnt := make([]int, 60)

	for i := range shifts {
		s := &shifts[i]
		if s.id != sleepiestguard {
			continue
		}
		for m := 0; m < 60; m++ {
			if s.isAsleep(m) {
				asleepcnt[m]++
			}
		}
	}

	maxasleepmin := 0
	for i := range asleepcnt {
		if asleepcnt[i] > asleepcnt[maxasleepmin] {
			maxasleepmin = i
		}
	}

	fmt.Printf("max asleep min %d\n", maxasleepmin)

	fmt.Printf("part 1: %d\n", maxasleepmin*sleepiestguard)

	/*
		for i := 0; i < 60; i++ {
			if shifts[0].isAsleep(i) {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	*/
}
