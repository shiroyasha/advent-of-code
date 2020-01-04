package main

import (
	"bufio"
	"container/heap"
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

	robots map[byte]Vec
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
				result.field[i+1][j] = '#'
				result.field[i-1][j] = '#'

				result.field[i][j+1] = '#'
				result.field[i][j-1] = '#'

				result.robots = map[byte]Vec{
					'1': Vec{X: j - 1, Y: i - 1},
					'2': Vec{X: j - 1, Y: i + 1},
					'3': Vec{X: j + 1, Y: i - 1},
					'4': Vec{X: j + 1, Y: i + 1},
				}
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

	for _, r := range []byte{'1', '2', '3', '4'} {
		g[r] = map[byte]Path{}

		for k, p := range m.keys {
			distance, doors, ok := m.findPath(m.robots[r], p)

			if ok {
				g[r][k] = Path{
					distance: distance,
					doorsSet: doors,
				}
			}
		}
	}

	for k1, p1 := range m.keys {
		g[k1] = map[byte]Path{}

		for k2, p2 := range m.keys {
			distance, doors, ok := m.findPath(p1, p2)

			if ok {
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
	q := &Heap{}

	heap.Init(q)
	heap.Push(q, Part1Item{pos: '@', keys: 0, distance: 0})

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

	for q.Len() > 0 {
		current := heap.Pop(q).(Part1Item)

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

			heap.Push(q, item)
		}
	}

	panic("never go here")
}

type Part2Item struct {
	pos      string
	distance int
	keys     uint
}

type Part2Seen struct {
	pos  string
	keys uint
}

type Heap2 []Part2Item

func (h Heap2) Len() int           { return len(h) }
func (h Heap2) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h Heap2) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap2) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Part2Item))
}

func (h *Heap2) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func part2(m Map, g Graph) int {
	seen := map[Part2Seen]bool{}
	q := &Heap2{}

	heap.Init(q)
	heap.Push(q, Part2Item{pos: "1234", keys: 0, distance: 0})

	possibleNext := func(current Part2Item) []string {
		res := []string{}

		for i, r := range []byte(current.pos) {
			for k, _ := range g[r] {
				if Has(current.keys, uint(k-'a')) {
					continue
				}

				if g[r][k].doorsSet|current.keys == current.keys {
					next := append([]byte(current.pos), []byte{}...)

					next[i] = k

					res = append(res, string(next))
				}
			}
		}

		return res
	}

	for q.Len() > 0 {
		current := heap.Pop(q).(Part2Item)

		// fmt.Printf("Next %s, keys %010b, distance %d\n", current.pos, current.keys, current.distance)

		s := Part2Seen{pos: current.pos, keys: current.keys}
		if seen[s] {
			continue
		}
		seen[s] = true

		if m.hasAllKeys(current.keys) {
			return current.distance
		}

		for _, next := range possibleNext(current) {
			nextKeys := current.keys
			distance := current.distance

			for _, r := range []byte(next) {
				if r >= 'a' && r <= 'z' {
					nextKeys = Add(nextKeys, uint(r-'a'))
				}
			}

			for i, _ := range []byte(current.pos) {
				distance += g[byte(current.pos[i])][byte(next[i])].distance
			}

			item := Part2Item{
				distance: distance,
				pos:      next,
				keys:     nextKeys,
			}

			// fmt.Printf("Next %s, keys %010b, distance %d\n", item.pos, item.keys, item.distance)
			// time.Sleep(1 * time.Second)

			heap.Push(q, item)
		}
	}

	panic("never go here")
}

func main() {
	m := load("input5.txt")
	m.show()
	g := m.DistanceGraph()
	fmt.Println(g['1'])

	fmt.Println("Graph ready")

	fmt.Println(part2(*m, g))

	// result := m.steps()

	// fmt.Println(result)
}
