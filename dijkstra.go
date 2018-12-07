package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type state struct {
	i, j int
	path string
}

// finds the closest node in the fringe, lastmin is an optimization, if we find a node that is at that distance we return it immediately (there can be nothing that's closer)
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
	fringe := map[state]int{state{0, 0, ""}: 0}   // nodes discovered but not visited (start at node 0 with distance 0)
	seen := map[state]bool{state{0, 0, ""}: true} // nodes already visited (we know the minimum distance of those)

	lastmin := 0

	cnt := 0

	for len(fringe) > 0 {
		cur := minimum(fringe, lastmin)

		if cnt%1000 == 0 {
			fmt.Printf("fringe %d (min dist %d)\n", len(fringe), fringe[cur])
		}
		cnt++

		if cur.i == 3 && cur.j == 3 {
			fmt.Printf("%v %d\n", cur, fringe[cur])
			return
		}

		distcur := fringe[cur]
		lastmin = distcur
		delete(fringe, cur)
		seen[cur] = true

		maybeadd := func(nb state, door byte) {
			// check if we can add the node
			if nb.i < 0 || nb.j < 0 || nb.i >= 4 || nb.j >= 4 || seen[nb] {
				return
			}
			switch door {
			case 'b', 'c', 'd', 'e', 'f':
				//ok
			default:
				return
			}

			// if we can add the node add it to the fringe
			// but first check that it's either a new node or we improved its distance
			if d, ok := fringe[nb]; !ok || distcur+1 < d {
				fringe[nb] = distcur + 1
			}
		}

		h := fmt.Sprintf("%x", md5.Sum([]byte(PASSCODE+cur.path)))

		// try to add all possible neighbors
		maybeadd(state{cur.i - 1, cur.j, cur.path + "U"}, h[0])
		maybeadd(state{cur.i + 1, cur.j, cur.path + "D"}, h[1])
		maybeadd(state{cur.i, cur.j - 1, cur.path + "L"}, h[2])
		maybeadd(state{cur.i, cur.j + 1, cur.path + "R"}, h[3])
	}
}

const PASSCODE = "bwnlcvfs"

func main() {
	fmt.Printf("hello\n")
	strconv.Atoi("0")
	search()
}

// 20m8s
