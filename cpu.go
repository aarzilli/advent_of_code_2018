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

type instr struct {
	opcode string
	args []arg
}

type arg struct {
	reg string
	val int
}

func (a arg) value(regs map[string]int) int {
	if a.reg == "" {
		return a.val
	}
	return regs[a.reg]
}

var text []instr

func run() {
	pc := 0
	snd := 0
	regs := map[string]int{}
	interpLoop:
	for {
		instr := text[pc]
		switch instr.opcode {
		case "snd":
			snd = instr.args[0].value(regs)
			pc++
		case "set":
			if instr.args[0].reg == "" {
				panic("blah")
			}
			regs[instr.args[0].reg] = instr.args[1].value(regs)
			pc++
		case "add":
			if instr.args[0].reg == "" {
				panic("blah")
			}
			regs[instr.args[0].reg] = instr.args[0].value(regs) + instr.args[1].value(regs)
			pc++
		case "mul":
			if instr.args[0].reg == "" {
				panic("blah")
			}
			regs[instr.args[0].reg] = instr.args[0].value(regs) * instr.args[1].value(regs)
			pc++
		case "mod":
			if instr.args[0].reg == "" {
				panic("blah")
			}
			regs[instr.args[0].reg] = instr.args[0].value(regs) % instr.args[1].value(regs)
			pc++
		case "rcv":
			if instr.args[0].value(regs) != 0 {
				fmt.Printf("recovered %d\n", snd)
				break interpLoop
			}
			pc++
		case "jgz":
			if instr.args[0].value(regs) > 0 {
				pc += instr.args[1].value(regs)
			} else {
				pc++
			}
		default:
			panic("blah")
		}
	}
}

func main() {
	buf, err := ioutil.ReadFile("18.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		opcode := fields[0]
		args := make([]arg, len(fields)-1)
		for i, field := range fields[1:] {
			n, err := strconv.Atoi(field)
			if err == nil {
				args[i] = arg{ val: n }
			} else {
				args[i] = arg{ reg: field }
			}
		}
		text = append(text, instr{ opcode, args })
	}
	
	//fmt.Printf("%#v\n", text)
	run()
}

// 14:02 (out of leaderboard)
