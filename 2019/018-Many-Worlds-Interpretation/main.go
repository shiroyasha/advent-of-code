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

	robot Vec

	robots Robots
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
				result.robot = Vec{X: j, Y: i}
				// result.field[i+1][j] = '#'
				// result.field[i-1][j] = '#'

				// result.field[i][j+1] = '#'
				// result.field[i][j-1] = '#'

				// result.robots = Robots{
				// 	p1: Vec{X: j - 1, Y: i - 1},
				// 	p2: Vec{X: j - 1, Y: i + 1},
				// 	p3: Vec{X: j + 1, Y: i - 1},
				// 	p4: Vec{X: j + 1, Y: i + 1},
				// }
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

			if pos == m.robots.p1 {
				res += fmt.Sprintf("\033[41m%s\033[0m", string(c))
			} else if pos == m.robots.p2 {
				res += fmt.Sprintf("\033[42m%s\033[0m", string(c))
			} else if pos == m.robots.p3 {
				res += fmt.Sprintf("\033[43m%s\033[0m", string(c))
			} else if pos == m.robots.p4 {
				res += fmt.Sprintf("\033[44m%s\033[0m", string(c))
			} else {
				res += fmt.Sprint(string(c))
			}
		}
		res += "\n"
	}

	fmt.Println(res)
}

type Graph map[byte]map[byte]Path

type Path struct {
	distance int
	doorsSet uint
}

type FindItem struct {
	pos      Vec
	distance int
	doorsSet uint
}

func (m *Map) findPath(p1, p2 Vec) (int, uint, bool) {
	seen := map[Vec]bool{}

	q := []FindItem{}
	q = append(q, FindItem{pos: p1, distance: 0, doorsSet: 0})

	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		if seen[current.pos] {
			continue
		}
		seen[current.pos] = true

		if m.at(current.pos) == '#' {
			continue
		}

		if current.pos == p2 {
			return current.distance, current.doorsSet, true
		}

		nextDoors := current.doorsSet
		if m.at(current.pos) >= 'A' && m.at(current.pos) <= 'Z' {
			nextDoors = Add(nextDoors, uint(m.at(current.pos)-'A'))
		}

		for _, next := range []Vec{up(current.pos), down(current.pos), left(current.pos), right(current.pos)} {
			q = append(q, FindItem{
				distance: current.distance + 1,
				doorsSet: nextDoors,
				pos:      next,
			})
		}
	}

	return MAX, uint(0), false
}

func (m *Map) DistanceGraph() Graph {
	g := Graph{}

	g['@'] = map[byte]Path{}

	for k, p := range m.keys {
		distance, doors, _ := m.findPath(m.robot, p)

		g['@'][k] = Path{
			distance: distance,
			doorsSet: doors,
		}
	}

	for k1, p1 := range m.keys {
		g[k1] = map[byte]Path{}

		for k2, p2 := range m.keys {
			if p1 != p2 {
				distance, doors, _ := m.findPath(p1, p2)

				g[k1][k2] = Path{
					distance: distance,
					doorsSet: doors,
				}
			}
		}
	}

	return g
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

type Robots struct {
	p1 Vec
	p2 Vec
	p3 Vec
	p4 Vec
}

type State struct {
	robots   Robots
	keys     uint
	distance int
}

type StateWithoutD struct {
	robots uint
}

func (m *Map) hasAllKeys(set uint) bool {
	all := uint(0)

	for k, _ := range m.keys {
		all = Add(all, uint(k-'a'))
	}

	return set == all
}

func zip(r Robots, keys uint) uint {
	res := keys

	res += uint(r.p1.X)
	res += uint(r.p1.Y * 100)

	res += uint(r.p2.X * 10000)
	res += uint(r.p2.Y * 1000000)

	res += uint(r.p3.X * 100000000)
	res += uint(r.p3.Y * 10000000000)

	res += uint(r.p4.X * 1000000000000)
	res += uint(r.p4.Y * 100000000000000)

	return res
}

func (m *Map) steps() int {
	seen := map[uint]struct{}{}
	q := []State{}

	q = append(q, State{
		robots:   m.robots,
		keys:     0,
		distance: 0,
	})

	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		// if len(seen)%100000 == 0 {
		// 	fmt.Println(len(seen))
		// }

		// m.show(current)
		// time.Sleep(1 * time.Second)

		// fmt.Printf("%v   -> %030b\n", current.pos, current.keys)

		z := zip(current.robots, current.keys)
		if _, ok := seen[z]; ok {
			continue
		}
		seen[z] = struct{}{}

		if m.at(current.robots.p1) >= 'A' && m.at(current.robots.p1) <= 'Z' && !Has(current.keys, uint(m.at(current.robots.p1))-uint('A')) {
			// fmt.Println("DOOR")
			continue
		}

		if m.at(current.robots.p2) >= 'A' && m.at(current.robots.p2) <= 'Z' && !Has(current.keys, uint(m.at(current.robots.p2))-uint('A')) {
			// fmt.Println("DOOR")
			continue
		}

		if m.at(current.robots.p3) >= 'A' && m.at(current.robots.p3) <= 'Z' && !Has(current.keys, uint(m.at(current.robots.p3))-uint('A')) {
			// fmt.Println("DOOR")
			continue
		}

		if m.at(current.robots.p4) >= 'A' && m.at(current.robots.p4) <= 'Z' && !Has(current.keys, uint(m.at(current.robots.p4))-uint('A')) {
			// fmt.Println("DOOR")
			continue
		}

		newKeys := current.keys

		if m.at(current.robots.p1) >= 'a' && m.at(current.robots.p1) <= 'z' {
			newKeys = Add(newKeys, uint(m.at(current.robots.p1)-'a'))
		}

		if m.at(current.robots.p2) >= 'a' && m.at(current.robots.p2) <= 'z' {
			newKeys = Add(newKeys, uint(m.at(current.robots.p2)-'a'))
		}

		if m.at(current.robots.p3) >= 'a' && m.at(current.robots.p3) <= 'z' {
			newKeys = Add(newKeys, uint(m.at(current.robots.p3)-'a'))
		}

		if m.at(current.robots.p4) >= 'a' && m.at(current.robots.p4) <= 'z' {
			newKeys = Add(newKeys, uint(m.at(current.robots.p4)-'a'))
		}

		if m.hasAllKeys(newKeys) {
			fmt.Println(current.distance)
			fmt.Println(newKeys)
			os.Exit(1)
		}

		for _, next := range []Vec{up(current.robots.p1), down(current.robots.p1), left(current.robots.p1), right(current.robots.p1)} {
			if m.at(next) == '#' {
				continue
			}

			q = append(q, State{
				robots: Robots{
					p1: next,
					p2: current.robots.p2,
					p3: current.robots.p3,
					p4: current.robots.p4,
				},
				keys:     newKeys,
				distance: current.distance + 1,
			})
		}

		for _, next := range []Vec{up(current.robots.p2), down(current.robots.p2), left(current.robots.p2), right(current.robots.p2)} {
			if m.at(next) == '#' {
				continue
			}

			q = append(q, State{
				robots: Robots{
					p1: current.robots.p1,
					p2: next,
					p3: current.robots.p3,
					p4: current.robots.p4,
				},
				keys:     newKeys,
				distance: current.distance + 1,
			})
		}

		for _, next := range []Vec{up(current.robots.p3), down(current.robots.p3), left(current.robots.p3), right(current.robots.p3)} {
			if m.at(next) == '#' {
				continue
			}

			q = append(q, State{
				robots: Robots{
					p1: current.robots.p1,
					p2: current.robots.p2,
					p3: next,
					p4: current.robots.p4,
				},
				keys:     newKeys,
				distance: current.distance + 1,
			})
		}

		for _, next := range []Vec{up(current.robots.p4), down(current.robots.p4), left(current.robots.p4), right(current.robots.p4)} {
			if m.at(next) == '#' {
				continue
			}

			q = append(q, State{
				robots: Robots{
					p1: current.robots.p1,
					p2: current.robots.p2,
					p3: current.robots.p3,
					p4: next,
				},
				keys:     newKeys,
				distance: current.distance + 1,
			})
		}
	}

	return 0
}

type Part1Item struct {
	pos      byte
	distance int
	keys     uint
}

type Part1Seen struct {
	pos  byte
	keys uint
}

type Heap []Part1Item

func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Part1Item))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func part1(m Map, g Graph) int {
	seen := map[Part1Seen]bool{}
	q := Heap{}

	q = append(q, Part1Item{pos: '@', keys: 0, distance: 0})

	possibleNext := func(current Part1Item) []byte {
		res := []byte{}

		for k, _ := range g[current.pos] {
			if Has(current.keys, uint(k-'a')) {
				continue
			}

			if g[current.pos][k].doorsSet|current.keys == current.keys {
				res = append(res, k)
			}
		}

		return res
	}

	for len(q) > 0 {
		current := q.Pop().(Part1Item)

		s := Part1Seen{pos: current.pos, keys: current.keys}
		if seen[s] {
			continue
		}
		seen[s] = true

		if m.hasAllKeys(current.keys) {
			return current.distance
		}

		for _, next := range possibleNext(current) {
			nextKeys := current.keys
			if next >= 'a' && next <= 'z' {
				nextKeys = Add(nextKeys, uint(next-'a'))
			}

			item := Part1Item{
				distance: current.distance + g[current.pos][next].distance,
				pos:      next,
				keys:     nextKeys,
			}

			q.Push(item)
		}
	}

	panic("never go here")
}

func main() {
	m := load("input3.txt")
	g := m.DistanceGraph()

	fmt.Println("Graph ready")

	fmt.Println(part1(*m, g))

	// result := m.steps()

	// fmt.Println(result)
}
