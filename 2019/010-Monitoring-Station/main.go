package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type AsteroidMap [][]bool

type Pos struct {
	X int
	Y int
}

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

		for i := 0; i < len(line); i++ {
			if line[i] == '.' {
				row = append(row, false)
			}

			if line[i] == '#' {
				row = append(row, true)
			}
		}

		m = append(m, row)
	}

	return AsteroidMap(m)
}

func (m *AsteroidMap) Display() {
	for _, row := range *m {
		for _, c := range row {
			if c {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println()
	}
}

func (m *AsteroidMap) Nuke(pos Pos) {
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

func Quadrant(p Pos) int {
	if p.X < 0 && p.Y < 0 {
		return 1
	}

	if p.X > 0 && p.Y < 0 {
		return 2
	}

	if p.X > 0 && p.Y > 0 {
		return 3
	}

	if p.X < 0 && p.Y > 0 {
		return 4
	}

	return 0
}

func CompareAngles(laserPos Pos, p1 Pos, p2 Pos) bool {
	x1 := p1.X - laserPos.X
	y1 := p1.Y - laserPos.Y

	x2 := p2.X - laserPos.X
	y2 := p2.Y - laserPos.Y

	// fmt.Printf("%v %v =>   (%5d %5d)   (%5d %5d)    =>   ", p1, p2, x1, y1, x2, y2)
	fmt.Printf("%v %v\n", p1, p2)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == p1.Y && j == p1.X {
				fmt.Printf("1")
				continue
			}

			if i == p2.Y && j == p2.X {
				fmt.Printf("2")
				continue
			}

			if i == 1 && j == 1 {
				fmt.Printf("@")
				continue
			}

			fmt.Print(".")
		}
		fmt.Println()
	}
	fmt.Println()

	defer func() {
		fmt.Println()
	}()

	res := false

	if x1 == 0 && x2 == 0 {
		res = y1 < y2
		fmt.Print("A-1 ", res)
		return res
	}

	if y1 == 0 && y2 == 0 {
		res = x1 < x2
		fmt.Print("A0 ", res)
		return res
	}

	if y1 == 0 {
		if x1 < 0 {
			res = true
			fmt.Print("A1 ", res)
			return res
		}

		if x1 > 0 {
			res = y2 > 0
			fmt.Print("A2 ", res)
			return res
		}
	}

	if y2 == 0 {
		if x2 < 0 {
			res = false
			fmt.Print("A3 ", res)
			return res
		}

		if x2 > 0 {
			res = y1 < 0
			fmt.Print("A32 ", res)
			return res
		}
	}

	if x1 == 0 {
		if y1 < 0 {
			res = x2 > 0
			fmt.Print("A4 ", res)
			return res
		}

		if y1 > 0 {
			res = x2 < 0
			fmt.Print("A5 ", res)
			return res
		}
	}

	if x2 == 0 {
		if y2 < 0 {
			res = x1 < 0
			fmt.Print("A6 ", res)
			return res
		}

		if y2 > 0 {
			res = x1 > 0
			fmt.Print("A6 ", res)
			return res
		}
	}

	// nothing can be zero anymore

	// quadrants

	q1 := Quadrant(Pos{X: x1, Y: y1})
	q2 := Quadrant(Pos{X: x2, Y: y2})

	if q1 != q2 {
		res = q1 < q2

		fmt.Printf("Q0 %d %d %t", q1, q2, res)
		return res
	}

	// same quadrant

	if q1 == 1 {
		res = float64(x1)/float64(y1) > float64(x2)/float64(y2)
		fmt.Print("Q1 ", res)
		return res
	}

	if q1 == 2 {
		res = float64(x1)/float64(y1) > float64(x2)/float64(y2)
		fmt.Print("Q2 ", res)
		return res
	}

	if q1 == 3 {
		res = float64(x1)/float64(y1) > float64(x2)/float64(y2)
		fmt.Print("Q3 ", res)
		return res
	}

	if q1 == 4 {
		res = float64(x1)/float64(y1) < float64(x2)/float64(y2)
		fmt.Print("Q4 ", res)
		return res
	}

	fmt.Print("NO")

	return true
}

func SortByAngle(laserPos Pos, positions []Pos) []Pos {
	sort.SliceStable(positions, func(i, j int) bool {
		return CompareAngles(laserPos, positions[i], positions[j])
	})

	return positions
}

func main() {
	m := load("input.txt")
	laserPos := Pos{X: 11, Y: 1}

	m.Nuke(laserPos)
	m.Display()

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

	for k, v := range result {
		fmt.Print(k)
		fmt.Println(v)

		if v > max {
			max = v
		}
	}

	fmt.Println(max)
	fmt.Println("-------------------------")

	// Part 2

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
		}
	}

	fmt.Println(len(nukeOrder))
	fmt.Println(nukeOrder)
	fmt.Println(nukeOrder[199])
	fmt.Println(nukeOrder[199].X*100 + nukeOrder[199].Y)
}
