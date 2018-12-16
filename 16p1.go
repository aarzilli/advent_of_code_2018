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

type Instr [4]int
type Regs struct {
	r [4]int
}

type Opcode func(instr Instr, regs *Regs) bool

var opcodes = []Opcode{
	addr,
	addi,
	mulr,
	muli,
	banr,
	bani,
	borr,
	bori,
	setr,
	seti,
	gtir,
	gtri,
	gtrr,
	eqir,
	eqri,
	eqrr,
}

func (regs *Regs) Val(i int) int {
	if i > 3 {
		panic("out of registers")
	}
	return regs.r[i]
}

func (regs *Regs) Set(i int, val int) {
	if i > 3 {
		panic("out of registers")
	}
	regs.r[i] = val
}

func addr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])+regs.Val(instr[2]))
	return true
}

func addi(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])+instr[2])
	return true
}

func mulr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])*regs.Val(instr[2]))
	return true
}

func muli(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])*instr[2])
	return true
}

func banr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])&regs.Val(instr[2]))
	return true
}

func bani(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])&instr[2])
	return true
}

func borr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])|regs.Val(instr[2]))
	return true
}

func bori(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1])|instr[2])
	return true
}

func setr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], regs.Val(instr[1]))
	return true
}

func seti(instr Instr, regs *Regs) (ok bool) {
	if instr[3] > 3 {
		return false
	}
	regs.Set(instr[3], instr[1])
	return true
}

func gtir(instr Instr, regs *Regs) (ok bool) {
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if instr[1] > regs.Val(instr[2]) {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func gtri(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if regs.Val(instr[1]) > instr[2] {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func gtrr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if regs.Val(instr[1]) > regs.Val(instr[2]) {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func eqir(instr Instr, regs *Regs) (ok bool) {
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if instr[1] == regs.Val(instr[2]) {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func eqri(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if regs.Val(instr[1]) == instr[2] {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func eqrr(instr Instr, regs *Regs) (ok bool) {
	if instr[1] > 3 {
		return false
	}
	if instr[2] > 3 {
		return false
	}
	if instr[3] > 3 {
		return false
	}
	if regs.Val(instr[1]) == regs.Val(instr[2]) {
		regs.Set(instr[3], 1)
	} else {
		regs.Set(instr[3], 0)
	}
	return true
}

func check(instr Instr, before, after Regs) int {
	cnt := 0
	for i := range opcodes {
		copy := before
		if !opcodes[i](instr, &copy) {
			continue
		}
		if copy == after {
			cnt++
		}
	}
	return cnt
}

const debug = false

func main() {
	buf, err := ioutil.ReadFile("16.txt")
	must(err)

	var before Regs
	var after Regs
	var instr Instr
	r := 0
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Before: ") {
			for i, v := range getints(line, true) {
				before.r[i] = v
			}
		} else if strings.HasPrefix(line, "After: ") {
			for i, v := range getints(line, true) {
				after.r[i] = v
			}
			cnt := check(instr, before, after)
			if debug {
				fmt.Printf("check %v %v %v: %d\n", before, instr, after, cnt)
			}
			if cnt >= 3 {
				r++
			}
		} else {
			for i, v := range vatoi(splitandclean(line, " ", -1)) {
				instr[i] = v
			}
		}
	}

	fmt.Printf("PART1 %d\n", r)

	//check(Instr{ 9, 2, 1, 2 }, Regs{ [4]int{ 3, 2, 1, 1 } }, Regs{ [4]int{ 3, 2, 2, 1 } })
}
