package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type AsteroidMap [][]bool

type Pos struct {
	X int
	Y int
}

var nukes = [][]int{}
var nuked = 0

func load(filename string) AsteroidMap {
	m := [][]bool{}

	file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break // EOF
		}

		row := []bool{}
		nukeRow := []int{}

		for i := 0; i < len(line); i++ {

			if line[i] == '.' {
				row = append(row, false)
				nukeRow = append(nukeRow, -1)
			}

			if line[i] == '#' {
				row = append(row, true)
				nukeRow = append(nukeRow, -1)
			}
		}

		m = append(m, row)
		nukes = append(nukes, nukeRow)
	}

	return AsteroidMap(m)
}

func (m *AsteroidMap) Display(pos Pos) {
	for i, row := range *m {
		for j, c := range row {
			if i == pos.Y && j == pos.X {
				fmt.Print("@")
			}

			if c {
				fmt.Print(".")
			} else {
				fmt.Print(" ")
			}
		}

		fmt.Println()
	}
}

func (m *AsteroidMap) Nukes(pos Pos) {
	if len(*m) != len(nukes) {
		panic("AAAA")
	}

	if len((*m)[0]) != len(nukes[0]) {
		fmt.Println(len(nukes[0]))
		panic(len((*m)[0]))
	}
	for j, row := range nukes {
		for i, c := range row {
			if j == pos.Y && i == pos.X {
				fmt.Print(" % ")
				continue
			}

			if c > 0 {
				fmt.Printf("%3d", c)
				continue
			}

			if (*m)[j][i] {
				fmt.Printf("  x")
			} else {
				fmt.Printf("  .")
			}

			// fmt.Printf("   ")

		}

		fmt.Println()
	}
}

func (m *AsteroidMap) Nuke(pos Pos) {
	nukes[pos.Y][pos.X] = nuked
	nuked++

	(*m)[pos.Y][pos.X] = false
}

func (m *AsteroidMap) Each(f func(Pos)) {
	for i := 0; i < len(*m); i++ {
		for j := 0; j < len((*m)[i]); j++ {
			if (*m)[i][j] {
				f(Pos{X: j, Y: i})
			}
		}
	}
}

func (m *AsteroidMap) InBetween(p1, p, p2 Pos) bool {
	if p1.X == p2.X {
		if p1.X != p.X {
			return false
		}

		if p1.Y <= p.Y && p.Y <= p2.Y {
			return true
		}

		if p1.Y >= p.Y && p.Y >= p2.Y {
			return true
		}

		return false
	}

	if p1.Y == p2.Y {
		if p1.Y != p.Y {
			return false
		}

		if p1.X <= p.X && p.X <= p2.X {
			return true
		}

		if p1.X >= p.X && p.X >= p2.X {
			return true
		}

		return false
	}

	a1 := p1.X - p2.X
	a2 := p1.X - p.X

	b1 := p1.Y - p2.Y
	b2 := p1.Y - p.Y

	if a1*b2 == a2*b1 {
		return ((p1.X <= p.X && p.X <= p2.X) || (p1.X >= p.X && p.X >= p2.X)) && ((p1.Y <= p.Y && p.Y <= p2.Y) || (p1.Y >= p.Y && p.Y >= p2.Y))
	} else {
		return false
	}
}

func (m *AsteroidMap) CanSee(p1, p2 Pos) bool {
	hasSomethingInBetween := true

	m.Each(func(p Pos) {
		if p1 == p || p2 == p {
			return
		}

		ok := m.InBetween(p1, p, p2)

		// fmt.Println("InBetween", p1, p, p2, ok)

		if ok {
			hasSomethingInBetween = false
		}
	})

	return hasSomethingInBetween
}

func (m *AsteroidMap) ListVisible(p Pos) []Pos {
	result := []Pos{}

	m.Each(func(p2 Pos) {
		if m.CanSee(p, p2) {
			result = append(result, p2)
		}
	})

	return result
}

func CompareAngles(laserPos Pos, p1 Pos, p2 Pos) bool {
	x1 := laserPos.X - p1.X
	y1 := laserPos.Y - p1.Y

	x2 := laserPos.X - p2.X
	y2 := laserPos.Y - p2.Y

	if x1 == 0 && y1 > 0 {
		return true
	}

	if x2 == 0 {
		return !CompareAngles(laserPos, p2, p1)
	}

	if x1 > 0 && x2 > 0 {
		return math.Atan(float64(y1)/float64(x1)) < math.Atan(float64(y2)/float64(x2))
	}

	if x1 < 0 && x2 < 0 {
		return math.Atan(float64(y1)/float64(x1)) < math.Atan(float64(y2)/float64(x2))
	}

	return x1 < x2
}

func SortByAngle(laserPos Pos, positions []Pos) []Pos {
	sort.SliceStable(positions, func(i, j int) bool {
		return CompareAngles(laserPos, positions[i], positions[j])
	})

	return positions
}

func main() {
	m := load("input.txt")

	result := map[Pos]int{}

	m.Each(func(p1 Pos) {
		m.Each(func(p2 Pos) {
			if p1 == p2 {
				return
			}

			can := m.CanSee(p1, p2)

			if can {
				result[p1]++
			}
		})
	})

	max := 0
	p := Pos{}

	for k, v := range result {
		fmt.Print(k)
		fmt.Println(v)

		if v > max {
			max = v
			p = k
		}
	}

	fmt.Println(max)
	fmt.Println("-------------------------")

	// Part 2

	laserPos := p
	m.Nuke(laserPos)

	nukeOrder := []Pos{}

	for {
		visible := m.ListVisible(laserPos)
		fmt.Println(visible)

		if len(visible) == 0 {
			break
		}

		nukeOrder = append(nukeOrder, SortByAngle(laserPos, visible)...)

		for _, n := range nukeOrder {
			m.Nuke(n)

			// cmd := exec.Command("clear")
			// cmd.Stdout = os.Stdout
			// cmd.Run()

			// m.Nukes(laserPos)

			// time.Sleep(50 * time.Millisecond)

			if nuked > 200 {
				break
			}
		}

		if nuked > 200 {
			break
		}
	}

	fmt.Println(len(nukeOrder))
	fmt.Println(nukeOrder)
	fmt.Println(nukeOrder[199])
	fmt.Println(nukeOrder[199].X*100 + nukeOrder[199].Y)

	m.Nukes(laserPos)
}
