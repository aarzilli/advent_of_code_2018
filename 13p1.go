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

var M [][]byte

type Coord struct {
	i, j int
}

type Cart struct {
	pos  Coord
	dir  string
	turn int
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
		panic("unfound collision?")
	})

	for ii := range order {
		i := order[ii]
		cart := &carts[i]

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

const debugsim = true

func checkcollision() {
	collided := false
	for i := range carts {
		for j := i + 1; j < len(carts); j++ {
			if carts[i].pos == carts[j].pos {
				fmt.Printf("collision at %d,%d (for %d %d)\n", carts[i].pos.j, carts[i].pos.i, i, j)
				collided = true
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
				carts = append(carts, Cart{Coord{i, j}, "left", 0})
				M[i][j] = '-'
			case '>':
				carts = append(carts, Cart{Coord{i, j}, "right", 0})
				M[i][j] = '-'
			case '^':
				carts = append(carts, Cart{Coord{i, j}, "up", 0})
				M[i][j] = '|'
			case 'v':
				carts = append(carts, Cart{Coord{i, j}, "down", 0})
				M[i][j] = '|'
			}
		}
	}

	fmt.Printf("%v\n", carts)

	for i := 0; i < 1000; i++ {
		simstep()
		checkcollision()
		if debugsim {
			fmt.Println()
			fmt.Println()
		}
	}
}
