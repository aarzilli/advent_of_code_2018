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

type Claim struct {
	id        int
	i, j      int
	w, h      int
	contested bool
}

var claims = []Claim{}
var M [][]int

func main() {
	buf, err := ioutil.ReadFile("03.txt")
	must(err)
	maxi, maxj := 0, 0
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := splitandclean(line, " ", -1)
		id := atoi(fields[0][1:])
		c := splitandclean(nolast(fields[2]), ",", 2)
		j := atoi(c[0])
		i := atoi(c[1])
		d := splitandclean(fields[3], "x", 2)
		w := atoi(d[0])
		h := atoi(d[1])
		claims = append(claims, Claim{
			id, i, j, w, h, false,
		})
		if i+h+1 > maxi {
			maxi = i + h + 1
		}
		if j+w+1 > maxj {
			maxj = j + w + 1
		}
	}

	//fmt.Printf("%v\n", claims)

	M = make([][]int, maxi*2)
	for i := range M {
		M[i] = make([]int, maxj*2)
	}

	for _, claim := range claims {
		lay(claim)
	}

	/*
		for i := range M {
			for j := range M[i] {
				if M[i][j] > 0 {
					fmt.Printf("%d", M[i][j])
				} else if M[i][j] < 0 {
					fmt.Printf("X")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
	*/

	contested := 0
	for i := range M {
		for j := range M[i] {
			if M[i][j] < 0 {
				contested++
			}
		}
	}
	fmt.Printf("%d\n", contested)

	for i := range claims {
		if !claims[i].contested {
			fmt.Printf("uncontested %d\n", claims[i].id)
		}
	}
}

func lay(claim Claim) {
	for i := 0; i < claim.h; i++ {
		for j := 0; j < claim.w; j++ {
			if M[i+claim.i][j+claim.j] != 0 {
				if M[i+claim.i][j+claim.j] > 0 {
					getClaimForID(M[i+claim.i][j+claim.j]).contested = true
				}
				getClaimForID(claim.id).contested = true
				M[i+claim.i][j+claim.j] = -1
			} else {
				M[i+claim.i][j+claim.j] = claim.id
			}
		}
	}
}

func getClaimForID(id int) *Claim {
	for i := range claims {
		if claims[i].id == id {
			return &claims[i]
		}
	}
	return nil
}
