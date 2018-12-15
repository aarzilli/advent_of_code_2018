package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
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

type Coord struct {
	i, j int
}

type Unit struct {
	typ    byte
	pos    Coord
	hp, ap int
	dead   bool
}

var M = [][]byte{}
var units []Unit

func printmatrix(matrix [][]byte) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%c", matrix[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func findreachable(pos Coord, reachset map[Coord]bool) {
	reachset[pos] = true
	for _, pos := range []Coord{{pos.i + 1, pos.j}, {pos.i - 1, pos.j}, {pos.i, pos.j - 1}, {pos.i, pos.j + 1}} {
		ti, tj := pos.i, pos.j
		if ti < 0 || ti >= len(M) {
			continue
		}
		if tj < 0 || tj >= len(M[ti]) {
			continue
		}
		if M[ti][tj] != '.' {
			continue
		}
		if reachset[pos] {
			continue
		}
		findreachable(pos, reachset)
	}
}

func posbefore(pos1, pos2 Coord) bool {
	if pos1.i < pos2.i {
		return true
	}
	if pos1.i > pos2.i {
		return false
	}
	return pos1.j < pos2.j
}

func dist(a, b Coord) int {
	return abs(a.i-b.i) + abs(a.j-b.j)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const debugMove = false
const debugAttack = false

type moveDest struct {
	pos   Coord
	steps []Coord
}

func (u *Unit) move() {
	reachset := map[Coord]bool{}
	findreachable(u.pos, reachset)

	moveTargets := map[Coord]bool{}

	addreachable := func(ti, tj int) {
		tp := Coord{ti, tj}
		if reachset[tp] {
			moveTargets[tp] = true
		}
	}

	for j := range units {
		if units[j].typ == u.typ {
			continue
		}
		p := units[j].pos
		addreachable(p.i-1, p.j)
		addreachable(p.i+1, p.j)
		addreachable(p.i, p.j-1)
		addreachable(p.i, p.j+1)
	}

	movev := []moveDest{}
	for p := range moveTargets {
		movev = append(movev, moveDest{
			pos:   p,
			steps: shortestpath(u.pos, p),
		})
	}

	sort.Slice(movev, func(i, j int) bool {
		if len(movev[i].steps) == len(movev[j].steps) {
			return posbefore(movev[i].pos, movev[j].pos)
		}
		return len(movev[i].steps) < len(movev[j].steps)
	})

	if debugMove {
		fmt.Printf("move targets for unit at %v: %v\n", u.pos, movev)
	}

	if len(movev) == 0 {
		return
	}

	if len(movev[0].steps) == 0 {
		return
	}

	if debugMove {
		fmt.Printf("unit at %v moves to %v\n", u.pos, movev[0].steps[0])
	}

	M[u.pos.i][u.pos.j] = '.'
	u.pos.i = movev[0].steps[0].i
	u.pos.j = movev[0].steps[0].j
	M[u.pos.i][u.pos.j] = u.typ
}

type path struct {
	at    Coord
	steps []Coord
}

func shortestpath(start, end Coord) []Coord {
	//fmt.Printf("shortest path %v %v\n", start, end)
	//defer fmt.Printf("done\n")
	fringe := []path{path{at: start, steps: nil}}
	seen := map[Coord]bool{}

	for len(fringe) > 0 {
		cur := fringe[0]
		seen[cur.at] = true
		fringe = fringe[1:]

		if cur.at == end {
			return cur.steps
		}

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
			if i < 0 || i >= len(M) {
				return
			}
			if j < 0 || j >= len(M[i]) {
				return
			}
			if M[i][j] != '.' {
				return
			}
			if infringe(i, j) || seen[Coord{i, j}] {
				return
			}

			var n path
			n.at = Coord{i, j}
			n.steps = make([]Coord, 0, len(cur.steps)+1)
			n.steps = append(n.steps, cur.steps...)
			n.steps = append(n.steps, n.at)

			fringe = append(fringe, n)
		}

		addstep(cur.at.i-1, cur.at.j)
		addstep(cur.at.i, cur.at.j-1)
		addstep(cur.at.i, cur.at.j+1)
		addstep(cur.at.i+1, cur.at.j)
	}

	panic("unreachable")
}

func findUnit(p Coord) int {
	for i := range units {
		if units[i].pos == p {
			return i
		}
	}
	return -1
}

func enemyof(typ byte) byte {
	switch typ {
	case 'E':
		return 'G'
	case 'G':
		return 'E'
	}
	panic("blah")
}

func (u *Unit) attack() bool {
	findenemy := func(si, sj int) int {
		if si < 0 || si >= len(M) {
			return -1
		}
		if sj < 0 || sj >= len(M[si]) {
			return -1
		}
		if M[si][sj] == enemyof(u.typ) {
			return findUnit(Coord{si, sj})
		}
		return -1
	}

	enemyv := []int{}

	addenemy := func(ei, ej int) {
		if id := findenemy(ei, ej); id >= 0 {
			enemyv = append(enemyv, id)
		}
	}

	addenemy(u.pos.i-1, u.pos.j)
	addenemy(u.pos.i+1, u.pos.j)
	addenemy(u.pos.i, u.pos.j-1)
	addenemy(u.pos.i, u.pos.j+1)

	if len(enemyv) == 0 {
		return false
	}

	sort.Slice(enemyv, func(i, j int) bool {
		eu1 := units[enemyv[i]]
		eu2 := units[enemyv[j]]
		if eu1.hp == eu2.hp {
			return posbefore(eu1.pos, eu2.pos)
		}
		return eu1.hp < eu2.hp
	})

	eu := &units[enemyv[0]]

	if debugAttack {
		fmt.Printf("unit at %v attacks unit at %v\n", u.pos, eu.pos)
	}

	eu.hp -= u.ap

	if eu.hp < 0 {
		if debugAttack {
			fmt.Printf("\tunit at %v is dead\n")
		}
		if eu.typ == 'E' && stoponelfdeath {
			fmt.Printf("elf died %d\n", round)
			os.Exit(1)
		}
		M[eu.pos.i][eu.pos.j] = '.'
		eu.pos.i = -100
		eu.pos.j = -100
		eu.dead = true
	}

	return true
}

func runturn() {
	order := make([]int, len(units))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		c1 := &units[order[i]]
		c2 := &units[order[j]]

		if c1.pos == c2.pos && c1.pos.i != -100 {
			panic("unfound collision?")
		}
		return posbefore(c1.pos, c2.pos)
	})

	for ii := range order {
		i := order[ii]
		u := &units[i]

		if u.dead {
			continue
		}

		attacked := u.attack()
		if !attacked {
			u.move()
			u.attack()
		}
	}
}

const debugFinished = false

func checkfinished(round int) bool {
	cnt := map[byte]int{}
	for i := range units {
		if units[i].dead {
			continue
		}
		cnt[units[i].typ] += units[i].hp
		if debugFinished {
			fmt.Printf("unit %c at %v is alive with %d hp\n", units[i].typ, units[i].pos, units[i].hp)
		}
	}
	if len(cnt) == 1 {
		for k, v := range cnt {
			fmt.Printf("%c won: %d (round: %d)\n", k, v, round)
			fmt.Printf("outcome %d\n", v*round)
			fmt.Printf("alt outcome %d\n", v*(round+1))
		}
		return true
	}
	return false
}

var stoponelfdeath = false

var round int

const debugMinimal = false

func main() {
	elfattack := 3
	if len(os.Args) == 2 {
		elfattack = atoi(os.Args[1])
		stoponelfdeath = true
	}
	buf, err := ioutil.ReadFile("15.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	for i := range M {
		for j := range M[i] {
			if M[i][j] == 'E' || M[i][j] == 'G' {
				attack := 3
				if M[i][j] == 'E' {
					attack = elfattack
				}
				units = append(units, Unit{M[i][j], Coord{i, j}, 200, attack, false})
			}
		}
	}

	for round = 0; ; round++ {
		if debugMinimal {
			fmt.Printf("at %d\n", round)
			printmatrix(M)
		}
		runturn()
		if checkfinished(round) {
			break
		}
	}
}

// run with
// for i in $(seq 4 100); do echo attack power $i; go run 15.go $i; done
// to find the answer for part2

/* PART1 OUTPUT:
G won: 2336 (round: 105)
outcome 245280
alt outcome 247616
*/

/* PART2 OUTPUT:
attack power 19
E won: 1442 (round: 52)
outcome 74984
alt outcome 76426
*/
