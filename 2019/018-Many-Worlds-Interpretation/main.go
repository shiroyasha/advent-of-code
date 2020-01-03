package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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

	robotPos Vec
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

			if result.field[i][j] >= 'a' && result.field[i][j] <= 'z' {
				result.keys[result.field[i][j]] = Vec{X: j, Y: i}
			}

			if result.field[i][j] == '@' {
				result.robotPos = Vec{X: j, Y: i}
			}
		}
	}

	return result
}

func (m *Map) show(s State) {
	res := ""

	for i, line := range m.field {
		for j, c := range line {
			pos := Vec{X: j, Y: i}

			if s.pos == pos {
				res += fmt.Sprintf("\033[41m%s\033[0m", string(c))
			} else {
				res += fmt.Sprint(string(c))
			}
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

func Add(s uint, i uint) uint {
	if i > 64 {
		panic("more then 64")
	}

	s = s | (uint(1) << i)

	return s
}

func Remove(s uint, i uint) uint {
	if i > 64 {
		panic("more then 64")
	}

	s = s & (^(uint(1) << i))

	return s
}

func Has(s uint, i uint) bool {
	if i > 64 {
		panic("more then 64")
	}

	mask := uint(1) << i

	return s&mask == mask
}

const MAX = 10000000000000

type State struct {
	pos      Vec
	keys     uint
	distance int
}

type StateWithoutD struct {
	pos  Vec
	keys uint
}

func (m *Map) hasAllKeys(set uint) bool {
	all := uint(0)

	for k, _ := range m.keys {
		all = Add(all, uint(k-'a'))
	}

	return set == all
}

func (m *Map) steps() int {
	seen := map[StateWithoutD]bool{}
	q := []State{}

	q = append(q, State{
		pos:      m.robotPos,
		keys:     0,
		distance: 0,
	})

	for len(q) > 0 {
		// fmt.Println("Queue len ", len(q))
		// fmt.Println("Seen len ", seen)
		// for _, k := range q {
		// 	fmt.Print(k.pos)
		// }
		// fmt.Println()
		// fmt.Println()

		// pop from queue

		current := q[0]
		q = q[1:]

		m.show(current)
		time.Sleep(1 * time.Second)

		// fmt.Printf("%v   -> %030b\n", current.pos, current.keys)

		if seen[StateWithoutD{pos: current.pos, keys: current.keys}] {
			continue
		}
		seen[StateWithoutD{pos: current.pos, keys: current.keys}] = true

		if m.at(current.pos) >= 'A' && m.at(current.pos) <= 'Z' && !Has(current.keys, uint(m.at(current.pos))-uint('A')) {
			// fmt.Println("DOOR")
			continue
		}

		newKeys := current.keys

		if m.at(current.pos) >= 'a' && m.at(current.pos) <= 'z' {
			newKeys = Add(newKeys, uint(m.at(current.pos)-'a'))
		}

		if m.hasAllKeys(newKeys) {
			fmt.Println(current.distance)
			fmt.Println(newKeys)
			os.Exit(1)
		}

		for _, next := range []Vec{up(current.pos), down(current.pos), left(current.pos), right(current.pos)} {
			if m.at(next) == '#' {
				continue
			}

			q = append(q, State{
				pos:      next,
				keys:     newKeys,
				distance: current.distance + 1,
			})
		}
	}

	return 0
}

func main() {
	m := load("input2.txt")
	m.field[m.robotPos.Y][m.robotPos.X] = '.'

	result := m.steps()

	fmt.Println(result)
}
