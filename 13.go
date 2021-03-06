package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var M [][]byte

type Coord struct {
	i, j int
}

type Cart struct {
	pos     Coord
	dir     string
	turn    int
	enabled bool
}

var carts = []Cart{}

func simstep() {
	order := make([]int, len(carts))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		c1 := &carts[order[i]]
		c2 := &carts[order[j]]

		if c1.pos.i < c2.pos.i {
			return true
		}
		if c1.pos.i > c2.pos.i {
			return false
		}
		if c1.pos.j < c2.pos.j {
			return true
		}
		if c1.pos.j > c2.pos.j {
			return false
		}
		if part1 {
			panic("unfound collision?")
		} else {
			return true
		}
	})

	for ii := range order {
		i := order[ii]
		cart := &carts[i]

		if !cart.enabled {
			continue
		}

		switch M[cart.pos.i][cart.pos.j] {
		case '/':
			switch cart.dir {
			case "up":
				cart.dir = "right"
			case "down":
				cart.dir = "left"
			case "left":
				cart.dir = "down"
			case "right":
				cart.dir = "up"
			}
		case '\\':
			switch cart.dir {
			case "up":
				cart.dir = "left"
			case "down":
				cart.dir = "right"
			case "left":
				cart.dir = "up"
			case "right":
				cart.dir = "down"
			}

		case '-':
			switch cart.dir {
			case "up":
				panic("bad dir")
			case "down":
				panic("bad dir")
			}

		case '|':
			switch cart.dir {
			case "left":
				panic("bad dir")
			case "right":
				panic("bad dir")
			}

		case '+':
			switch cart.dir {
			case "up":
				switch cart.turn {
				case 0:
					cart.dir = "left"
				case 1:
					cart.dir = "up"
				case 2:
					cart.dir = "right"
				default:
					panic("wtf")
				}
			case "down":
				switch cart.turn {
				case 0:
					cart.dir = "right"
				case 1:
					cart.dir = "down"
				case 2:
					cart.dir = "left"
				default:
					panic("wtf")
				}
			case "left":
				switch cart.turn {
				case 0:
					cart.dir = "down"
				case 1:
					cart.dir = "left"
				case 2:
					cart.dir = "up"
				default:
					panic("wtf")
				}
			case "right":
				switch cart.turn {
				case 0:
					cart.dir = "up"
				case 1:
					cart.dir = "right"
				case 2:
					cart.dir = "down"
				default:
					panic("wtf")
				}
			}

			cart.turn++
			if cart.turn > 2 {
				cart.turn = 0
			}

		default:
			fmt.Printf("cart %d innawoods %c going %s\n", i, M[cart.pos.i][cart.pos.j], cart.dir)
			panic("innawoods")
		}

		switch cart.dir {
		case "up":
			cart.pos.i--
		case "down":
			cart.pos.i++
		case "left":
			cart.pos.j--
		case "right":
			cart.pos.j++
		default:
			panic("unknown dir")
		}

		if debugsim {
			fmt.Printf("cart %d at i=%d,j=%d going %s\n", i, cart.pos.i, cart.pos.j, cart.dir)
		}
		checkcollision()
	}
}

const part1 = false
const debugsim = false

func checkcollision() {
	collided := false
	for i := range carts {
		if !carts[i].enabled {
			continue
		}
		for j := i + 1; j < len(carts); j++ {
			if !carts[j].enabled {
				continue
			}
			if carts[i].pos == carts[j].pos {
				fmt.Printf("collision at %d,%d (for %d %d)\n", carts[i].pos.j, carts[i].pos.i, i, j)
				if part1 {
					collided = true
				} else {
					carts[i].enabled = false
					carts[j].enabled = false
				}
			}
		}
	}
	if collided {
		os.Exit(1)
	}
}

func main() {
	buf, err := ioutil.ReadFile("13.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}
		M = append(M, []byte(line))
	}

	for i := range M {
		for j := range M[i] {
			switch M[i][j] {
			case '<':
				carts = append(carts, Cart{Coord{i, j}, "left", 0, true})
				M[i][j] = '-'
			case '>':
				carts = append(carts, Cart{Coord{i, j}, "right", 0, true})
				M[i][j] = '-'
			case '^':
				carts = append(carts, Cart{Coord{i, j}, "up", 0, true})
				M[i][j] = '|'
			case 'v':
				carts = append(carts, Cart{Coord{i, j}, "down", 0, true})
				M[i][j] = '|'
			}
		}
	}

	fmt.Printf("%v\n", carts)

	for i := 0; i < 100000; i++ {
		simstep()
		checkcollision()
		if debugsim {
			fmt.Println()
			fmt.Println()
		}

		cnt := 0
		idx := 0
		for i := range carts {
			if carts[i].enabled {
				idx = i
				cnt++
			}
		}
		if cnt == 1 {
			fmt.Printf("last cart is %d at %d,%d\n", idx, carts[idx].pos.j, carts[idx].pos.i)
			return
		}
	}
}
