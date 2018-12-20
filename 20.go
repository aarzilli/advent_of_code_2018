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

var Graph = map[Coord]*Node{}

type Node struct {
	Pos   Coord
	Child []*Node
}

type Coord struct {
	i, j int
}

var max, min Coord

func lookup(p Coord) *Node {
	if p.i > max.i {
		max.i = p.i
	}
	if p.j > max.j {
		max.j = p.j
	}
	if p.i < min.i {
		min.i = p.i
	}
	if p.j < min.j {
		min.j = p.j
	}
	if _, ok := Graph[p]; !ok {
		Graph[p] = &Node{Pos: p}
	}
	return Graph[p]
}

const debugParse = false

func parseAlternatives(in []byte, pos Coord, toplevel bool, indent string) []byte {
	if debugParse {
		fmt.Printf("%sstart alternative:\n", indent)
	}
	for len(in) > 0 {
		in = parseSingleRun(in, pos, toplevel, indent+"\t")
		if len(in) == 0 {
			break
		}

		ch := in[0]
		in = in[1:]

		switch ch {
		case '|':
			// continue iterating
		case ')':
			if toplevel {
				panic("wtf")
			}
			return in
		default:
			panic("wtf?")
		}

		if debugParse {
			fmt.Printf("%sor:\n", indent)
		}
	}

	return in
}

func addlink(a, b Coord) {
	na := lookup(a)
	nb := lookup(b)
	na.Child = append(na.Child, nb)
	nb.Child = append(nb.Child, na)
}

func haslink(a, b Coord) bool {
	na := Graph[a]
	nb := Graph[b]
	if na == nil || nb == nil {
		return false
	}
	for _, child := range na.Child {
		if child.Pos == b {
			return true
		}
	}
	return false
}

func parseSingleRun(in []byte, pos Coord, toplevel bool, indent string) []byte {
	for len(in) > 0 {
		ch := in[0]
		oldin := in
		in = in[1:]

		switch ch {
		case 'N':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			nextPos := Coord{pos.i - 1, pos.j}
			addlink(pos, nextPos)
			pos = nextPos
		case 'S':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			nextPos := Coord{pos.i + 1, pos.j}
			addlink(pos, nextPos)
			pos = nextPos
		case 'E':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			nextPos := Coord{pos.i, pos.j + 1}
			addlink(pos, nextPos)
			pos = nextPos
		case 'W':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			nextPos := Coord{pos.i, pos.j - 1}
			addlink(pos, nextPos)
			pos = nextPos
		case '(':
			in = parseAlternatives(in, pos, false, indent+"\t")
		case '|':
			return oldin
		case ')':
			if toplevel {
				panic("wtf")
			}
			return oldin
		default:
			panic("wtf?")
		}
	}

	if !toplevel {
		panic("wtf?")
	}

	return in
}

type path struct {
	at   Coord
	dist int
}

func shortestdist(start Coord) map[Coord]int {
	fringe := []path{path{at: start, dist: 0}}
	distmap := map[Coord]int{}

	for len(fringe) > 0 {
		cur := fringe[0]
		distmap[cur.at] = cur.dist
		fringe = fringe[1:]

		infringe := func(i, j int) bool {
			p := Coord{i, j}
			for k := range fringe {
				if fringe[k].at == p {
					return true
				}
			}
			return false
		}
		addstep := func(i, j int) {
			if infringe(i, j) {
				return
			}
			if _, ok := distmap[Coord{i, j}]; ok {
				return
			}
			if !haslink(cur.at, Coord{i, j}) {
				return
			}
			fringe = append(fringe, path{Coord{i, j}, cur.dist + 1})
		}

		addstep(cur.at.i-1, cur.at.j)
		addstep(cur.at.i, cur.at.j-1)
		addstep(cur.at.i, cur.at.j+1)
		addstep(cur.at.i+1, cur.at.j)
	}

	return distmap
}

const debugDist = false

func main() {
	buf, err := ioutil.ReadFile("20.txt")
	must(err)
	input := strings.TrimSpace(string(buf))
	if input[0] != '^' {
		panic("wtf")
	}
	input = input[1:]
	if input[len(input)-1] != '$' {
		panic("wtf")
	}
	input = input[:len(input)-1]
	parseAlternatives([]byte(input), Coord{0, 0}, true, "")
	distmap := shortestdist(Coord{0, 0})
	if debugDist {
		fmt.Printf("%v\n", distmap)
	}
	cnt1000 := 0
	max := 0
	for _, d := range distmap {
		if d > max {
			max = d
		}
		if d >= 1000 {
			cnt1000++
		}
	}
	fmt.Printf("PART1: %d\n", max)
	fmt.Printf("PART2: %d\n", cnt1000)
}

// with greather than: 8782
// with greather than or equal: 8784
