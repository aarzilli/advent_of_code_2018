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

var opcodesMap = []Opcode{
	{"banr", banr},
	{"eqrr", eqrr},
	{"setr", setr},
	{"eqir", eqir},
	{"bori", bori},
	{"muli", muli},
	{"bani", bani},
	{"borr", borr},
	{"gtir", gtir},
	{"gtrr", gtrr},
	{"addi", addi},
	{"gtri", gtri},
	{"eqri", eqri},
	{"addr", addr},
	{"mulr", mulr},
	{"seti", seti},
}

type Opcode struct {
	Name string
	F    func(instr Instr, regs *Regs) bool
}

type Regs struct {
	r [6]int
}

func (regs *Regs) Set(i int, val int) {
	regs.r[i] = val
}

func (regs *Regs) Val(i int) int {
	return regs.r[i]
}

func addr(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])+regs.Val(instr.arg[2]))
	return true
}

func addi(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])+instr.arg[2])
	return true
}

func mulr(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])*regs.Val(instr.arg[2]))
	return true
}

func muli(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])*instr.arg[2])
	return true
}

func banr(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])&regs.Val(instr.arg[2]))
	return true
}

func bani(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])&instr.arg[2])
	return true
}

func borr(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])|regs.Val(instr.arg[2]))
	return true
}

func bori(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1])|instr.arg[2])
	return true
}

func setr(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], regs.Val(instr.arg[1]))
	return true
}

func seti(instr Instr, regs *Regs) (ok bool) {
	regs.Set(instr.arg[3], instr.arg[1])
	return true
}

func gtir(instr Instr, regs *Regs) (ok bool) {
	if instr.arg[1] > regs.Val(instr.arg[2]) {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

func gtri(instr Instr, regs *Regs) (ok bool) {
	if regs.Val(instr.arg[1]) > instr.arg[2] {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

func gtrr(instr Instr, regs *Regs) (ok bool) {
	if regs.Val(instr.arg[1]) > regs.Val(instr.arg[2]) {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

func eqir(instr Instr, regs *Regs) (ok bool) {
	if instr.arg[1] == regs.Val(instr.arg[2]) {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

func eqri(instr Instr, regs *Regs) (ok bool) {
	if regs.Val(instr.arg[1]) == instr.arg[2] {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

func eqrr(instr Instr, regs *Regs) (ok bool) {
	if regs.Val(instr.arg[1]) == regs.Val(instr.arg[2]) {
		regs.Set(instr.arg[3], 1)
	} else {
		regs.Set(instr.arg[3], 0)
	}
	return true
}

type Instr struct {
	opcode string
	arg    []int
}

func parseInstr(line string) Instr {
	fields := splitandclean(line, " ", -1)
	opcode := fields[0]
	args := make([]int, len(fields))
	args[0] = 0
	for i, field := range fields[1:] {
		if field[len(field)-1] == ',' {
			field = field[:len(field)-1]
		}
		var err error
		args[i+1], err = strconv.Atoi(field)
		must(err)
	}
	return Instr{opcode, args}
}

var text []Instr
var ipreg int

func findOpcode(name string) Opcode {
	for i := range opcodesMap {
		if opcodesMap[i].Name == name {
			return opcodesMap[i]
		}
	}
	panic("not found???")
}

func run() {
	var regs Regs
	ip := 0
	for ip < len(text) {
		regs.Set(ipreg, ip)
		instr := text[regs.Val(ipreg)]
		fmt.Printf("%d %v %v\n", ip, instr, regs)
		ok := findOpcode(instr.opcode).F(instr, &regs)
		fmt.Printf("-> %v\n", regs)
		if !ok {
			panic("not ok???")
		}
		ip = regs.Val(ipreg)
		ip++
	}
	fmt.Printf("%v\n", regs)
}

func main() {
	buf, err := ioutil.ReadFile("19.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		instr := parseInstr(line)
		if instr.opcode == "#ip" {
			ipreg = instr.arg[1]
		} else {
			text = append(text, instr)
		}
	}
	fmt.Printf("ipreg %d\n", ipreg)
	fmt.Printf("%v\n", text)
	run()
}
