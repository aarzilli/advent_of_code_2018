package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
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

var M = [][]byte{}

const SZ = 2000
const debugIn = false

type Coord struct {
	i, j int
}

func (c *Coord) At() byte {
	return M[c.i][c.j]
}

func (c *Coord) Flood(ch byte) {
	M[c.i][c.j] = ch
}

var dispcnt = 0

func dbgdisp(maxi, minj, maxj int) {
	for i := 0; i <= maxi+1; i++ {
		for j := minj; j <= maxj; j++ {
			fmt.Printf("%c", M[i][j])
		}
		fmt.Printf("\n")
	}

	fmt.Printf("disp %d\n", dispcnt)
}

func makeimage(maxi, minj, maxj int) {
	fmt.Printf("makeimage %d\n", dispcnt)
	im := image.NewRGBA(image.Rect(0, 0, maxj-minj, maxi))
	stride := (maxj - minj) * 4
	for i := 0; i < maxi; i++ {
		for j := 0; j < maxj-minj; j++ {
			switch M[i][j+minj] {
			case '#':
				im.Pix[i*stride+j*4] = 0x00
				im.Pix[i*stride+j*4+1] = 0x00
				im.Pix[i*stride+j*4+2] = 0x00
				im.Pix[i*stride+j*4+3] = 0xff
			case '~':
				im.Pix[i*stride+j*4] = 0x00
				im.Pix[i*stride+j*4+1] = 0x00
				im.Pix[i*stride+j*4+2] = 0xff
				im.Pix[i*stride+j*4+3] = 0xff
			case '|':
				im.Pix[i*stride+j*4] = 0xff
				im.Pix[i*stride+j*4+1] = 0x00
				im.Pix[i*stride+j*4+2] = 0x00
				im.Pix[i*stride+j*4+3] = 0xff
			case '.':
				im.Pix[i*stride+j*4] = 0xff
				im.Pix[i*stride+j*4+1] = 0xff
				im.Pix[i*stride+j*4+2] = 0xff
				im.Pix[i*stride+j*4+3] = 0xff
			}
		}
	}
	fh, err := os.Create(fmt.Sprintf("fuck/%05d.png", dispcnt))
	must(err)
	png.Encode(fh, im)
	fh.Close()
	dispcnt++
}

func findBottom(pos Coord, maxi int) Coord {
	for M[pos.i][pos.j] != '#' {
		M[pos.i][pos.j] = '|'
		pos.i++
		if pos.i > maxi {
			return pos
		}
	}
	pos.i--
	return pos
}

func fill(pos Coord) []Coord {
	overflows := []Coord{}
	didoverflow := false

	for {
		tostill := []Coord{}
		for cur := pos; cur.At() != '#'; cur.j-- {
			cur.Flood('|')
			tostill = append(tostill, cur)
			if M[cur.i+1][cur.j] != '#' && M[cur.i+1][cur.j] != '~' {
				if M[cur.i+1][cur.j] != '|' {
					overflows = append(overflows, Coord{cur.i + 1, cur.j})
				}
				didoverflow = true
				break
			}
		}

		for cur := pos; cur.At() != '#'; cur.j++ {
			cur.Flood('|')
			tostill = append(tostill, cur)
			if M[cur.i+1][cur.j] != '#' && M[cur.i+1][cur.j] != '~' {
				if M[cur.i+1][cur.j] != '|' {
					overflows = append(overflows, Coord{cur.i + 1, cur.j})
				}
				didoverflow = true
				break
			}
		}

		if didoverflow {
			return overflows
		}

		for _, cur := range tostill {
			cur.Flood('~')
		}

		pos.i--
	}
}

const debugExample = false
const debugPix = false

func main() {
	buf, err := ioutil.ReadFile("17.txt")
	must(err)

	M = make([][]byte, SZ)
	for i := range M {
		M[i] = make([]byte, SZ)
		for j := range M[i] {
			M[i][j] = '.'
		}
	}

	mini := SZ
	maxi := 0
	maxj := 0
	minj := 0

	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := splitandclean(line, ",", -1)

		if strings.HasPrefix(v[0], "x=") {
			// vertical line
			j := getints(v[0], false)[0]
			rng := vatoi(splitandclean(v[1][2:], "..", -1))
			starti, endi := rng[0], rng[1]
			if debugIn {
				fmt.Printf("vertical line at j=%d from %d to %d\n", j, starti, endi)
			}

			if j > maxj {
				maxj = j
			}
			if j < minj {
				minj = j
			}

			if starti > maxi {
				maxi = starti
			}
			if starti < mini {
				mini = starti
			}
			if endi > maxi {
				maxi = endi
			}
			if endi < mini {
				mini = endi
			}

			for i := starti; i <= endi; i++ {
				M[i][j] = '#'
			}
		} else if strings.HasPrefix(v[0], "y=") {
			// horizontal line

			i := getints(v[0], false)[0]

			rng := vatoi(splitandclean(v[1][2:], "..", -1))
			startj, endj := rng[0], rng[1]
			if debugIn {
				fmt.Printf("horizontal line at i=%d from %d to %d\n", i, startj, endj)
			}

			if i > maxi {
				maxi = i
			}
			if i < mini {
				mini = i
			}

			if startj > maxj {
				maxj = startj
			}
			if startj < minj {
				minj = startj
			}
			if endj > maxj {
				maxj = endj
			}
			if endj < minj {
				minj = endj
			}

			for j := startj; j <= endj; j++ {
				M[i][j] = '#'
			}
		}
	}

	fmt.Printf("y range %d..%d\n", mini, maxi)

	rivers := []Coord{{0, 500}}

	curmaxi := 0
	minj = 300

	for len(rivers) > 0 {
		pos := rivers[len(rivers)-1]
		rivers = rivers[:len(rivers)-1]
		pos = findBottom(pos, maxi)
		if pos.i >= maxi {
			continue
		}
		if pos.i > curmaxi {
			curmaxi = pos.i
		}
		rivers = append(rivers, fill(pos)...)
		if debugExample {
			dbgdisp(curmaxi, minj, maxj)
			fmt.Println()
			fmt.Println()
		}
		if debugPix {
			makeimage(curmaxi, minj, maxj)
		}
	}

	if debugExample {
		dbgdisp(maxi, minj, maxj)
	}

	makeimage(SZ, 0, SZ)

	cnt := 0
	cnt2 := 0
	for i := range M {
		if i > maxi {
			continue
		}
		if i < mini {
			continue
		}
		for j := range M[i] {
			if i == 0 && j == 500 {
				continue
			}
			if M[i][j] == '~' || M[i][j] == '|' {
				cnt++
			}
			if M[i][j] == '~' {
				cnt2++
			}
		}
	}

	fmt.Printf("PART1 %d\n", cnt)
	fmt.Printf("PART2 %d\n", cnt2)
}
