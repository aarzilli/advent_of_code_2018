package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	_ "os"
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

const part2 = true

func main() {
	buf, err := ioutil.ReadFile("09.txt")
	must(err)
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		fields := splitandclean(line, " ", -1)
		nplayers, lastpoints := atoi(fields[0]), atoi(fields[6])
		if part2 {
			lastpoints = lastpoints * 100
		}
		part1alt(nplayers, lastpoints)
		//part1(nplayers, lastpoints)
		if lastpoints < 100 {
			break
		}
	}
}

func part1alt(nplayers, lastpoints int) {
	fmt.Printf("# %d %d\n", nplayers, lastpoints)
	circle := []int{ 0, 2,  1,  3 }
	newcircle := []int{ 0, 4 }
	var rewrite []int
	rewritescratch := [8]int{}
	players := make([]int, nplayers)
	
	
	src := 1
	marble := 5
	
	incplayer := func(points int) {
		players[(marble-1)%len(players)] += points
		if debugInc {
			fmt.Printf("player %d gets %d points\n", (marble-1)%len(players), points)
		}
	}
	
	
	for {
		if len(rewrite) > 0 {
			newcircle = append(newcircle, rewrite[0])
			rewrite = rewrite[1:]
		} else {
			newcircle = append(newcircle, circle[src])
			src++
		}
		
		if marble % 23 == 0 {
			incplayer(marble)
			if len(newcircle) > 9 {
				if debug2 {
					fmt.Printf("%v\n", newcircle)
				}
				incplayer(newcircle[len(newcircle)-9])
				if debug2 {
					fmt.Printf("adding: %d\n", newcircle[len(newcircle)-9])
				}
				rewrite = rewritescratch[:]
				copy(rewrite, newcircle[len(newcircle)-8:])
				newcircle = newcircle[:len(newcircle)-9]
				newcircle = append(newcircle, rewrite[:2]...)
				rewrite = rewrite[2:]
				if debug2 {
					fmt.Printf("restarting with: %v\n", newcircle)
					fmt.Printf("rewriting: %v\n", rewrite)
				}
			} else {
				//slow path
				if debug3 {
					fmt.Printf("%v\n", circle)
					fmt.Printf("%v\n", newcircle)
				}
				
				rem := len(circle)+len(newcircle)-9
				
				if debug3 {
					fmt.Printf("index %d: %d\n", rem, circle[rem])
				}
				
				incplayer(circle[rem])
				
				newcircle = append(newcircle, circle[src:rem]...)
				src = rem+1
				
				if src < len(circle) {
					newcircle = append(newcircle, circle[src])
					src++
					if src < len(circle) {
						newcircle = append(newcircle, circle[src])
						src++
					} else {
						circle = newcircle
						newcircle = make([]int, 0, len(circle)*3)
						newcircle = append(newcircle, circle[0])
						src = 1
					}
				} else {
					circle = newcircle
					newcircle = make([]int, 0, len(circle)*3)
					newcircle = append(newcircle, circle[:2]...)
					src = 2
				}
			}
			marble++
		}
		
		newcircle = append(newcircle, marble)
		marble++
		if marble > lastpoints {
			break
		}
		if src >= len(circle) && len(rewrite) <= 0 {
			circle = newcircle
			newcircle = make([]int, 0, len(circle)*3)
			src = 0
			if debug {
				printcircle(circle, src)
			}
		}
	}
	
	if debug {
		printcircle(circle, src)
	}
	max := 0
	for i := range players {
		if players[i] > max {
			max = players[i]
		}
	}
	fmt.Printf("highest score %d (exp: 393229)\n", max)
}




const debug = false
const debug2 = false
const debug3 = false
const debugInc = false

func part1(nplayers, lastpoints int) {
	fmt.Printf("# %d %d\n", nplayers, lastpoints)
	//circle := []int{ 0 }
	circle := make([]int, 0, lastpoints)
	newcircle := make([]int, 0, lastpoints)
	_ = newcircle
	circle = append(circle, 0)
	i := 0
	players := make([]int, nplayers)
	_ = players
	for marble := 1; marble <= lastpoints; marble++ {
		incplayer := func(points int) {
			players[(marble-1)%len(players)] += points
			if debugInc {
				fmt.Printf("player %d gets %d points\n", (marble-1)%len(players), points)
			}
		}
		if marble % 1000 == 0 {
			fmt.Printf("progress %d %0.02g%%\n", marble, float64(marble)/float64(lastpoints) * 100)
		}
		if marble % 23 == 0 {
			if debug2 {
				fmt.Printf("before:\n")
				printcircle(circle, i)
			}
			incplayer(marble)
			//newcircle := make([]int, 0, len(circle)-1)
			ri1 := (i-7)
			if ri1 < 0 {
				ri1 = len(circle)+ri1
			} else {
				ri1 = ri1 % len(circle)
			}
			ri2 := (i-6)
			if ri2 < 0 {
				ri2 = len(circle)+ri2
			} else {
				ri2 = ri2 % len(circle)
			}
			
			incplayer(circle[ri1])
			
			copy(circle[ri1:], circle[ri2:])
			newi := ri1
			circle = circle[:len(circle)-1]
			
			/*
			newcircle = newcircle[:0]
			newcircle = append(newcircle, circle[:ri1]...)
			if debug2 {
				fmt.Printf("removing %d\n", circle[ri1])
			}
			players[(marble-1)%len(players)] += circle[ri1]
			newi := len(newcircle)
			newcircle = append(newcircle, circle[ri2:]...)
			temp := circle
			circle = newcircle
			newcircle = temp
			*/
			
			i = newi
			if debug2 {
				fmt.Printf("after:\n")
				printcircle(circle, i)
			}
			continue
		}
		if debug {
			printcircle(circle, i)
		}
		//newcircle = make([]int, 0, len(circle)+1)
		
		r1 := (i+1)%len(circle)+1
		circle = circle[:len(circle)+1]
		copy(circle[r1+1:], circle[r1:])
		circle[r1] = marble
		newi := r1
		
		/*newcircle = newcircle[:0]
		newcircle = append(newcircle, circle[:(i+1)%len(circle)+1]...)
		newi := len(newcircle)
		newcircle = append(newcircle, marble)
		newcircle = append(newcircle, circle[(i+1)%len(circle)+1:]...)
		//circle = newcircle
		temp := circle
		circle = newcircle
		newcircle = temp*/
		
		i = newi
	}
	if debug {
		printcircle(circle, i)
	}
	max := 0
	for i := range players {
		if players[i] > max {
			max = players[i]
		}
	}
	fmt.Printf("highest score %d (exp: 393229)\n", max)
}

func printcircle(circle []int, cur int) {
	for i := range circle {
		if i == cur {
			fmt.Printf("(%d)", circle[i])
		} else {
			fmt.Printf(" %d ", circle[i])
		}
	}
	fmt.Printf("\n")
}
