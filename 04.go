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

type Shift struct {
	id    int
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

type guardmin struct {
	id int
	m  int
}

const part1 = true

func main() {
	buf, err := ioutil.ReadFile("04.txt")
	must(err)

	lines := strings.Split(string(buf), "\n")
	sort.Strings(lines) // fuck you eric

	shifts := []Shift{}
	var curshift *Shift

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "begins shift") {
			id := atoi(splitandclean(line, " ", -1)[3][1:])
			shifts = append(shifts, Shift{id: id})
			curshift = &shifts[len(shifts)-1]
		}

		v := vatoi(splitandclean(nolast(splitandclean(line, " ", -1)[1]), ":", 2))
		minute := v[0]*60 + v[1]
		if minute > 60 {
			minute = minute - (24 * 60)
		}

		if curshift != nil {
			curshift.times = append(curshift.times, minute)
		}
	}

	if part1 {
		//fmt.Printf("%#v\n", shifts)

		asleep := map[int]int{}
		for i := range shifts {
			n := shifts[i].asleep()
			asleep[shifts[i].id] += n
		}

		//fmt.Printf("%#v\n", asleep)

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
	}
	asleepguardmin := map[guardmin]int{}

	for i := range shifts {
		s := &shifts[i]

		for m := 0; m < 60; m++ {
			if s.isAsleep(m) {
				asleepguardmin[guardmin{s.id, m}]++
			}
		}
	}

	maxgm := guardmin{}
	for gm := range asleepguardmin {
		if asleepguardmin[gm] > asleepguardmin[maxgm] {
			maxgm = gm
		}
	}
	fmt.Printf("PART2: %v %d\n", maxgm, maxgm.id*maxgm.m)
}
