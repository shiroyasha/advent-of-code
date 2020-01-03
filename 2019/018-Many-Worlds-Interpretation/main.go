package main

import (
	"bufio"
	"fmt"
	"os"
)

type Vec struct {
	X, Y int
}

type Map struct {
	field [][]byte

	doors map[byte]Vec
	keys  map[byte]Vec

	width  int
	height int
}

func load(filename string) *Map {
	result := &Map{}
	result.doors = map[byte]Vec{}
	result.keys = map[byte]Vec{}

	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		r := []byte{}
		for _, k := range line[0 : len(line)-1] {
			r = append(r, byte(k))
		}

		result.field = append(result.field, r)
	}

	result.height = len(result.field)
	result.width = len(result.field[0])

	for i := 0; i < result.height; i++ {
		for j := 0; j < result.width; j++ {
			if result.field[i][j] >= 'A' && result.field[i][j] <= 'Z' {
				result.doors[result.field[i][j]] = Vec{X: j, Y: i}
			}

			if (result.field[i][j] >= 'a' && result.field[i][j] <= 'z') || result.field[i][j] == '@' {
				result.keys[result.field[i][j]] = Vec{X: j, Y: i}
			}
		}
	}

	return result
}

func (m *Map) show() {
	res := ""

	for _, line := range m.field {
		for _, c := range line {
			res += fmt.Sprint(string(c))
		}
		res += "\n"
	}

	fmt.Println(res)
}

func (m *Map) at(p Vec) byte {
	return m.field[p.Y][p.X]
}

func up(pos Vec) Vec {
	return Vec{X: pos.X, Y: pos.Y - 1}
}

func down(pos Vec) Vec {
	return Vec{X: pos.X, Y: pos.Y + 1}
}

func left(pos Vec) Vec {
	return Vec{X: pos.X - 1, Y: pos.Y}
}

func right(pos Vec) Vec {
	return Vec{X: pos.X + 1, Y: pos.Y}
}

const MAX = 10000000000000

func (m *Map) unlock(door byte) {
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			vec := Vec{X: j, Y: i}

			if m.at(vec) == door {
				fmt.Println("unlocked", door)
				m.field[i][j] = '.'
				m.show()
				return
			}
		}
	}
}

func (m *Map) steps() int {
	q := map[Vec]bool{}
	dist := map[Vec]int{}
	prev := map[Vec]Vec{}

	start := m.keys['@']

	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			vec := Vec{X: j, Y: i}
			val := m.at(vec)
			if val != '#' {
				q[vec] = true
			}

			dist[vec] = MAX
		}
	}

	dist[start] = 0
	q[start] = true

	fmt.Println(q)

	for len(q) > 0 {
		current := Vec{}
		currentFound := false

		for k, _ := range q {
			if m.at(k) >= 'A' && m.at(k) <= 'Z' {
				continue
			}

			if !currentFound || dist[k] < dist[current] {
				current = k
				currentFound = true
			}
		}

		delete(q, current)

		fmt.Println(current, q)

		for _, next := range []Vec{up(current), down(current), left(current), right(current)} {
			if q[next] == false {
				continue
			}

			if m.at(next) >= 'a' && m.at(next) <= 'z' {
				m.unlock(m.at(next) - 'a' + 'A')
			}

			alt := dist[current] + 1
			if alt < dist[next] {
				dist[next] = alt
				prev[next] = current
			}
		}
	}

	for k, v := range dist {
		fmt.Println(k, v)
	}

	return dist[m.keys['b']]
}

func main() {
	m := load("input1.txt")

	m.show()

	result := m.steps()

	fmt.Println(result)
}
