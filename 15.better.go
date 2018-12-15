package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"os"
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

func posbefore(pos1, pos2 Coord) bool {
	if pos1.i == pos2.i {
		return pos1.j < pos2.j
	}
	return pos1.i < pos2.i
}

type moveDest struct {
	pos  Coord
	path path
}

func (u *Unit) move() {
	pathmap := shortestpaths(u.pos)

	moveTargets := map[Coord]*path{}

	addifempty := func(ti, tj int) {
		p := Coord{ti, tj}
		if M[ti][tj] != '.' || moveTargets[p] != nil {
			return
		}
		moveTargets[p] = pathmap[p]
	}

	for j := range units {
		if units[j].dead || units[j].typ == u.typ {
			continue
		}
		p := units[j].pos
		addifempty(p.i-1, p.j)
		addifempty(p.i+1, p.j)
		addifempty(p.i, p.j-1)
		addifempty(p.i, p.j+1)
	}

	movev := []moveDest{}
	for p, path := range moveTargets {
		if path == nil {
			continue
		}
		movev = append(movev, moveDest{
			pos:  p,
			path: *path,
		})
	}

	sort.Slice(movev, func(i, j int) bool {
		if movev[i].path.dist == movev[j].path.dist {
			return posbefore(movev[i].pos, movev[j].pos)
		}
		return movev[i].path.dist < movev[j].path.dist
	})

	if len(movev) == 0 || movev[0].path.dist == 0 {
		return
	}

	M[u.pos.i][u.pos.j] = '.'
	u.pos = movev[0].path.firstStep
	M[u.pos.i][u.pos.j] = u.typ
}

type path struct {
	at        Coord
	firstStep Coord
	dist      int
}

func shortestpaths(start Coord) map[Coord]*path {
	fringe := []path{path{at: start, dist: 0}}
	pathmap := map[Coord]*path{}

	for len(fringe) > 0 {
		cur := fringe[0]
		pathmap[cur.at] = &cur
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
			if M[i][j] != '.' || infringe(i, j) || (pathmap[Coord{i, j}] != nil) {
				return
			}

			var n path
			n.at = Coord{i, j}
			n.dist = cur.dist + 1
			if cur.dist == 0 {
				n.firstStep = n.at
			} else {
				n.firstStep = cur.firstStep
			}
			fringe = append(fringe, n)
		}

		addstep(cur.at.i-1, cur.at.j)
		addstep(cur.at.i, cur.at.j-1)
		addstep(cur.at.i, cur.at.j+1)
		addstep(cur.at.i+1, cur.at.j)
	}

	return pathmap
}

func findUnit(p Coord) *Unit {
	for i := range units {
		if units[i].pos == p {
			return &units[i]
		}
	}
	return nil
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

const elfDied = "elf died"

func (u *Unit) attack() bool {
	findenemy := func(si, sj int) *Unit {
		if M[si][sj] == enemyof(u.typ) {
			return findUnit(Coord{si, sj})
		}
		return nil
	}

	enemyv := []*Unit{}

	addenemy := func(ei, ej int) {
		if u := findenemy(ei, ej); u != nil {
			enemyv = append(enemyv, u)
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
		if enemyv[i].hp == enemyv[j].hp {
			return posbefore(enemyv[i].pos, enemyv[j].pos)
		}
		return enemyv[i].hp < enemyv[j].hp
	})

	enemyv[0].hp -= u.ap

	if enemyv[0].hp < 0 {
		if enemyv[0].typ == 'E' && stoponelfdeath {
			panic(elfDied)
		}
		M[enemyv[0].pos.i][enemyv[0].pos.j] = '.'
		enemyv[0].pos.i = -100
		enemyv[0].pos.j = -100
		enemyv[0].dead = true
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

func checkfinished(round int) bool {
	cnt := map[byte]int{}
	for i := range units {
		if units[i].dead {
			continue
		}
		cnt[units[i].typ] += units[i].hp
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

func initialize(elfattack int) {
	path := "15.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	buf, err := ioutil.ReadFile(path)
	must(err)
	M = nil
	units = nil
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
}

func alive(typ byte) int {
	hp := 0
	for i := range units {
		if units[i].typ == typ {
			hp += units[i].hp
		}
	}
	return hp
}

func battle() bool {
	var round int
	defer func() {
		ierr := recover()
		if ierr == nil {
			return
		}
		errstr, _ := ierr.(string)
		if errstr != elfDied {
			panic(ierr)
		}
		fmt.Printf("%v at %d\n", ierr, round)
	}()
	for round = 0; ; round++ {
		if round % 1000 == 0 && round != 0 {
			fmt.Printf("round %d: elves: %dHP goblins %dHP\n", round, alive('E'), alive('G'))
		}
		runturn()
		if checkfinished(round) {
			return true
		}
	}
	return false
}

func main() {
	// PART 1

	fmt.Printf("PART 1:\n")
	initialize(3)
	battle()

	fmt.Printf("PART 2:\n")
	stoponelfdeath = true
	for elfattack := 4; ; elfattack++ {
		fmt.Printf("trying with elf attack power %d\n", elfattack)
		initialize(elfattack)
		if battle() {
			break
		}
	}
}
