package main

import (
	"bufio"
	"fmt"
	"os"
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

func main() {
	m := load("input.txt")
	m.Display()

	result := map[Pos]int{}

	m.Each(func(p1 Pos) {
		m.Each(func(p2 Pos) {
			if p1 == p2 {
				return
			}

			can := m.CanSee(p1, p2)
			// fmt.Println("Testing", p1, p2, can)

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
}
