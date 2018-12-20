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

type Graph map[Coord]*Node

type Node struct {
	Pos   Coord
	Child []*Node
}

type Coord struct {
	i, j int
}

func (Graph Graph) lookup(p Coord) *Node {
	if _, ok := Graph[p]; !ok {
		Graph[p] = &Node{Pos: p}
	}
	return Graph[p]
}

const debugParse = false

func (Graph Graph) parseAlternatives(in []byte, vpos []Coord, toplevel bool, indent string) ([]byte, []Coord) {
	if debugParse {
		fmt.Printf("%sstart alternative:\n", indent)
	}
	outpos := []Coord{}
	for len(in) > 0 {
		var newpos []Coord
		in, newpos = Graph.parseSingleRun(in, clone(vpos), toplevel, indent+"\t")
		outpos = append(outpos, newpos...)
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
			return in, dedup(outpos)
		default:
			panic("wtf?")
		}

		if debugParse {
			fmt.Printf("%sor:\n", indent)
		}
	}

	return in, dedup(outpos)
}

func clone(vpos []Coord) []Coord {
	r := make([]Coord, len(vpos))
	copy(r, vpos)
	return r
}

func dedup(vpos []Coord) []Coord {
	m := make(map[Coord]bool, len(vpos))
	for i := range vpos {
		m[vpos[i]] = true
	}

	r := make([]Coord, 0, len(m))
	for pos := range m {
		r = append(r, pos)
	}
	return r
}

func (Graph Graph) addlink(a, b Coord) {
	na := Graph.lookup(a)
	nb := Graph.lookup(b)
	na.Child = append(na.Child, nb)
	nb.Child = append(nb.Child, na)
}

func (Graph Graph) haslink(a, b Coord) bool {
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

func (Graph Graph) parseSingleRun(in []byte, vpos []Coord, toplevel bool, indent string) ([]byte, []Coord) {
	for len(in) > 0 {
		ch := in[0]
		oldin := in
		in = in[1:]

		switch ch {
		case 'N':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			for i := range vpos {
				nextPos := Coord{vpos[i].i - 1, vpos[i].j}
				Graph.addlink(vpos[i], nextPos)
				vpos[i] = nextPos
			}
		case 'S':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			for i := range vpos {
				nextPos := Coord{vpos[i].i + 1, vpos[i].j}
				Graph.addlink(vpos[i], nextPos)
				vpos[i] = nextPos
			}
		case 'E':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			for i := range vpos {
				nextPos := Coord{vpos[i].i, vpos[i].j + 1}
				Graph.addlink(vpos[i], nextPos)
				vpos[i] = nextPos
			}
		case 'W':
			if debugParse {
				fmt.Printf("%schar %c\n", indent, ch)
			}
			for i := range vpos {
				nextPos := Coord{vpos[i].i, vpos[i].j - 1}
				Graph.addlink(vpos[i], nextPos)
				vpos[i] = nextPos
			}
		case '(':
			in, vpos = Graph.parseAlternatives(in, vpos, false, indent+"\t")
			//TODO: this is actually wrong, I should be collecting the end position from every alternative and using it instead of pos for every subsequent iteration (pos should be a list everywhere)
		case '|':
			return oldin, vpos
		case ')':
			if toplevel {
				panic("wtf")
			}
			return oldin, vpos
		default:
			panic("wtf?")
		}
	}

	if !toplevel {
		panic("wtf?")
	}

	return in, vpos
}

type path struct {
	at   Coord
	dist int
}

func (Graph Graph) shortestdist(start Coord) map[Coord]int {
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
			if !Graph.haslink(cur.at, Coord{i, j}) {
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

func problem(input string) (max, cnt1000 int) {
	Graph := Graph{}

	if input[0] != '^' {
		panic("wtf")
	}
	input = input[1:]
	if input[len(input)-1] != '$' {
		panic("wtf")
	}
	input = input[:len(input)-1]
	Graph.parseAlternatives([]byte(input), []Coord{{0, 0}}, true, "")
	distmap := Graph.shortestdist(Coord{0, 0})
	if debugDist {
		fmt.Printf("%v\n", distmap)
	}
	cnt1000 = 0
	max = 0
	for _, d := range distmap {
		if d > max {
			max = d
		}
		if d >= 1000 {
			cnt1000++
		}
	}
	return max, cnt1000
}

type Example struct {
	input string
	max   int
}

var examples = []Example{
	{"^WNE$", 3},
	{"^ENWWW(NEEE|SSE(EE|N))$", 10},
	{"^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$", 18},
	{"^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$", 23},
	{"^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$", 31},
}

func main() {
	for _, example := range examples {
		max, _ := problem(example.input)
		if max != example.max {
			fmt.Printf("mismatch %q expected %d got %d\n", example.input, example.max, max)
		}
	}

	buf, err := ioutil.ReadFile("20.txt")
	must(err)
	input := strings.TrimSpace(string(buf))
	max, cnt1000 := problem(input)
	fmt.Printf("PART1: %d (3151)\n", max)
	fmt.Printf("PART2: %d (8784)\n", cnt1000)
}
