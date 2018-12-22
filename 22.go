package main

import (
	"fmt"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Coord struct {
	x, y int
}

var erosionLvl = map[Coord]int{}
var M = map[Coord]byte{}

func geoidx(p Coord) int {
	switch {
	case p == Coord{0, 0}:
		return 0
	case p == target:
		return 0
	case p.y == 0:
		return p.x * 16807
	case p.x == 0:
		return p.y * 48271
	default:
		return getErosionLvl(Coord{p.x - 1, p.y}) * getErosionLvl(Coord{p.x, p.y - 1})
		//return erosionLvl[] * erosionLvl[]
	}
}

func getErosionLvl(p Coord) int {
	if e, ok := erosionLvl[p]; ok {
		return e
	}
	erosionLvl[p] = (geoidx(p) + depth) % 20183
	switch erosionLvl[p] % 3 {
	case 0:
		M[p] = '.'
	case 1:
		M[p] = '='
	case 2:
		M[p] = '|'
	}
	/*
		if p == (Coord{0,0}) {
			M[p] = 'M'
		}
		if p == target {
			M[p] = 'T'
		}*/
	if p == target {
		if M[p] != '.' {
			panic("target isn't rocky")
		}
	}
	return erosionLvl[p]
}

// example
/*
var target Coord = Coord{ 10, 10 }
var depth int = 510
*/

// real input
var target = Coord{12, 763}
var depth int = 7740

func example(x, y int) {
	p := Coord{x, y}
	fmt.Printf("%d %d\n", geoidx(p), getErosionLvl(p))
	fmt.Printf("%c\n", M[p])
}

const (
	equipNeither = 0
	equipTorch   = 1
	equipClimb   = 2
)

type state struct {
	at        Coord
	equipment int
}

func minimum(fringe map[state]int, lastmin int) state {
	var mink state
	first := true
	for k, d := range fringe {
		if first {
			mink = k
			first = false
		}
		if d == lastmin {
			return k
		}
		if d < fringe[mink] {
			mink = k
		}
	}
	return mink
}

func search() {
	getErosionLvl(Coord{0, 0})
	fringe := map[state]int{state{at: Coord{0, 0}, equipment: equipTorch}: 0}   // nodes discovered but not visited (start at node 0 with distance 0)
	seen := map[state]bool{state{at: Coord{0, 0}, equipment: equipTorch}: true} // nodes already visited (we know the minimum distance of those)
	mindist := map[state]int{state{at: Coord{0, 0}, equipment: equipTorch}: 0}

	lastmin := 0

	cnt := 0

	for len(fringe) > 0 {
		cur := minimum(fringe, lastmin)

		if cnt%1000 == 0 {
			fmt.Printf("fringe %d (min dist %d)\n", len(fringe), fringe[cur])
		}
		cnt++

		if cur.at == target && cur.equipment == equipTorch {
			fmt.Printf("%v %d\n", cur, fringe[cur])
			return
		}

		distcur := fringe[cur]
		lastmin = distcur
		delete(fringe, cur)
		seen[cur] = true
		mindist[cur] = distcur

		switch M[cur.at] {
		case '.':
			if cur.equipment == 0 {
				panic("bad equipment on rock")
			}
		case '=':
			if cur.equipment == equipTorch {
				panic("bad equipment on water")
			}
		case '|':
			if cur.equipment == equipClimb {
				panic("bad equipment on narrow")
			}
		default:
			panic("blah! 3")
		}

		if cur.equipment != 0 && cur.equipment != equipTorch && cur.equipment != equipClimb {
			panic("blah! 4")
		}

		addswitch := func(newequip int) {
			getErosionLvl(cur.at)
			switch M[cur.at] {
			case '.':
				if newequip == 0 {
					return
				}
			case '=':
				if newequip == equipTorch {
					return
				}
			case '|':
				if newequip == equipClimb {
					return
				}
			default:
				panic("blah! 1")
			}

			if riskfrom(cur.at.x, cur.at.y) == newequip {
				panic("shouldn't have swapped")
			}

			newdist := distcur + 7
			newstate := state{at: cur.at, equipment: newequip}

			if seen[newstate] {
				return
			}

			if d, ok := fringe[newstate]; !ok || newdist < d {
				fringe[newstate] = newdist
			}
		}

		addmove := func(x, y int) {
			if x < 0 || y < 0 || y > depth {
				return
			}
			p := Coord{x, y}
			getErosionLvl(p)
			switch M[p] {
			case '.':
				if cur.equipment == 0 {
					return
				}
			case '=':
				if cur.equipment == equipTorch {
					return
				}
			case '|':
				if cur.equipment == equipClimb {
					return
				}
			default:
				panic("blah! 2")
			}

			if riskfrom(p.x, p.y) == cur.equipment {
				panic("shouldn't have added?")
			}

			newdist := distcur + 1
			newstate := state{at: p, equipment: cur.equipment}

			if seen[newstate] {
				return
			}

			if d, ok := fringe[newstate]; !ok || newdist < d {
				fringe[newstate] = newdist
			}
		}

		// try to add all possible neighbors
		switch cur.equipment {
		case 0:
			addswitch(equipTorch)
			addswitch(equipClimb)
		case equipTorch:
			addswitch(0)
			addswitch(equipClimb)
		case equipClimb:
			addswitch(0)
			addswitch(equipTorch)
		}
		addmove(cur.at.x-1, cur.at.y)
		addmove(cur.at.x+1, cur.at.y)
		addmove(cur.at.x, cur.at.y-1)
		addmove(cur.at.x, cur.at.y+1)
	}
}

func riskfrom(x, y int) int {
	getErosionLvl(Coord{x, y})
	switch M[Coord{x, y}] {
	case '.':
		return 0
	case '=':
		return 1
	case '|':
		return 2
	default:
		panic("blah")
	}
}

func main() {
	/*
		for y := 0; y < target.y+1000; y++ {
			for x := 0; x < target.x+1000; x++ {
				getErosionLvl(Coord{ x, y })
				if x == target.x+1000-1 {
					fmt.Printf("%c", M[Coord{ x, y }])
				} else {
					fmt.Printf("%c ", M[Coord{ x, y }])
				}
			}
			fmt.Printf("\n")
		}

		return*/

	risk := 0
	for y := 0; y <= target.y; y++ {
		for x := 0; x <= target.x; x++ {
			risk += riskfrom(x, y)
		}
	}

	fmt.Printf("PART1 %d\n", risk)

	getErosionLvl(Coord{0, 0})
	fmt.Printf("%c\n", M[Coord{0, 0}])

	search()
}

// 779 too low
