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

func getints(in string, hasneg bool) []int {
	v := getnums(in, hasneg, false)
	return vatoi(v)
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

type Opcode struct {
	Name string
	F    func(instr Instr, regs *Regs) bool
}

var opcodes = []Opcode{
	{"addr", addr},
	{"addi", addi},
	{"mulr", mulr},
	{"muli", muli},
	{"banr", banr},
	{"bani", bani},
	{"borr", borr},
	{"bori", bori},
	{"setr", setr},
	{"seti", seti},
	{"gtir", gtir},
	{"gtri", gtri},
	{"gtrr", gtrr}, //ok
	{"eqir", eqir},
	{"eqri", eqri},
	{"eqrr", eqrr},
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

func check(instr Instr, before, after Regs) map[string]bool {
	opnames := map[string]bool{}
	for i := range opcodes {
		copy := before
		if !opcodes[i].F(instr, &copy) {
			continue
		}
		if copy == after {
			opnames[opcodes[i].Name] = true
		}
	}
	return opnames
}

const debug = false

func main() {
	buf, err := ioutil.ReadFile("16.txt")
	must(err)

	var before Regs
	var after Regs
	var instr Instr
	r := 0

	possibleNames := make([]map[string]bool, len(opcodes))

	for i := range opcodes {
		names := map[string]bool{}
		for _, opcode2 := range opcodes {
			names[opcode2.Name] = true
		}
		possibleNames[i] = names
	}

	realRegs := Regs{[4]int{0, 0, 0, 0}}
	_ = realRegs

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
			opnames := check(instr, before, after)
			if debug {
				fmt.Printf("check %v %v %v: %v\n", before, instr, after, opnames)
			}
			if len(opnames) >= 3 {
				r++
			}
			for name := range possibleNames[instr[0]] {
				if !opnames[name] {
					delete(possibleNames[instr[0]], name)
				}
			}
		} else {
			for i, v := range vatoi(splitandclean(line, " ", -1)) {
				instr[i] = v
			}
			opcodesMap[instr[0]].F(instr, &realRegs)
		}
	}

	fmt.Printf("PART1 %d\n", r)

	fmt.Printf("before:\n")
	for i := range possibleNames {
		fmt.Printf("%v\n", possibleNames[i])
	}

	finalMap := make([]string, len(opcodes))

	cleanupPossibleNames := func() bool {
		for i := range possibleNames {
			if len(possibleNames[i]) == 1 {
				for v := range possibleNames[i] {
					fmt.Printf("found %d is %s\n", i, v)
					finalMap[i] = v
				}
				for j := range possibleNames {
					if possibleNames[j][finalMap[i]] {
						delete(possibleNames[j], finalMap[i])
					}
				}
				return true
			}
		}
		return false
	}

	for {
		ok := cleanupPossibleNames()
		if !ok {
			break
		}
	}

	fmt.Printf("var opcodesMap = []Opcode{\n")
	for i := range finalMap {
		fmt.Printf("{ %q, %s },\n", finalMap[i], finalMap[i])
	}
	fmt.Printf("}\n")

	fmt.Printf("PART 2: %d\n", realRegs.r[0])
	//check(Instr{ 9, 2, 1, 2 }, Regs{ [4]int{ 3, 2, 1, 1 } }, Regs{ [4]int{ 3, 2, 2, 1 } })
}
