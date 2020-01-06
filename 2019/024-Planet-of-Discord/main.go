package main

import (
	"fmt"
)

var input1 = []string{
	"...#.",
	"#.##.",
	"#..##",
	"#.###",
	"##...",
}

var input2 = []string{
	"....#",
	"#..#.",
	"#..##",
	"..#..",
	"#....",
}

func Set(s uint, i, j int, v bool) uint {
	if v {
		return Add(s, uint(i*5+j))
	} else {
		return Remove(s, uint(i*5+j))
	}
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

func load(input []string) uint {
	res := uint(0)

	for i, line := range input {
		for j, c := range line {
			if c == '#' {
				res = Add(res, uint(i*len(line)+j))
			}
		}
	}

	return res
}

func show(pattern uint) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if alive(pattern, i, j) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println()
	}
}

func alive(pattern uint, i, j int) bool {
	if i < 0 || j < 0 {
		return false
	}

	if i > 4 || j > 4 {
		return false
	}

	return Has(pattern, uint(i*5+j))
}

func bugCount(pattern uint, i, j int) int {
	count := 0

	if alive(pattern, i-1, j) {
		count += 1
	}

	if alive(pattern, i, j-1) {
		count += 1
	}

	if alive(pattern, i, j+1) {
		count += 1
	}

	if alive(pattern, i+1, j) {
		count += 1
	}

	return count
}

func shouldLive(pattern uint, i, j int) bool {
	return bugCount(pattern, i, j) == 1
}

func shouldSpawn(pattern uint, i, j int) bool {
	c := bugCount(pattern, i, j)
	return c == 1 || c == 2
}

func part1() {
	pattern := load(input1)

	layouts := map[uint]bool{}

	for {
		nextPattern := uint(0)

		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if alive(pattern, i, j) && shouldLive(pattern, i, j) {
					nextPattern = Add(nextPattern, uint(i*5+j))
				}

				if !alive(pattern, i, j) && shouldSpawn(pattern, i, j) {
					nextPattern = Add(nextPattern, uint(i*5+j))
				}
			}
		}

		pattern = nextPattern

		if layouts[pattern] {
			fmt.Println(pattern)
			return
		}

		layouts[pattern] = true
	}
}

type Field struct {
	levels []uint
}

func NewField(pattern uint) *Field {
	f := &Field{
		levels: make([]uint, 1000),
	}

	f.levels[0] = pattern

	return f
}

type Vec2 struct {
	x, y int
}

type Vec3 struct {
	x, y, levelInc int
}

// "...#.",
// "#.##.",
// "#..##",
// "#.###",
// "##...",

func (f *Field) get(level int, i, j int) int {
	level = level * 2
	if level < 0 {
		level = -level
		level -= 1
	}

	// fmt.Println(level, i, j)
	// fmt.Println()

	if Has(f.levels[level], uint(i*5+j)) {
		return 1
	} else {
		return 0
	}
}

func (f *Field) up(level, i, j int) int {
	if i == 0 {
		return f.get(level+1, 1, 2)
	} else if i == 3 {
		return 0 +
			f.get(level-1, 4, 0) +
			f.get(level-1, 4, 1) +
			f.get(level-1, 4, 2) +
			f.get(level-1, 4, 3) +
			f.get(level-1, 4, 4)
	} else {
		return f.get(level, i-1, j)
	}
}

func (f *Field) down(level, i, j int) int {
	if i == 4 {
		return f.get(level+1, 3, 2)
	} else if i == 1 {
		return 0 +
			f.get(level-1, 0, 0) +
			f.get(level-1, 0, 1) +
			f.get(level-1, 0, 2) +
			f.get(level-1, 0, 3) +
			f.get(level-1, 0, 4)
	} else {
		return f.get(level, i+1, j)
	}
}

func (f *Field) left(level, i, j int) int {
	if j == 0 {
		return f.get(level+1, 2, 1)
	} else if j == 3 {
		return 0 +
			f.get(level-1, 0, 4) +
			f.get(level-1, 1, 4) +
			f.get(level-1, 2, 4) +
			f.get(level-1, 3, 4) +
			f.get(level-1, 4, 4)
	} else {
		return f.get(level, i, j-1)
	}
}

func (f *Field) right(level, i, j int) int {
	if j == 4 {
		return f.get(level+1, 2, 3)
	} else if j == 1 {
		return 0 +
			f.get(level-1, 0, 0) +
			f.get(level-1, 1, 0) +
			f.get(level-1, 2, 0) +
			f.get(level-1, 3, 0) +
			f.get(level-1, 4, 0)
	} else {
		return f.get(level, i, j+1)
	}
}

func (f *Field) count(level, i, j int) int {
	return f.up(level, i, j) + f.down(level, i, j) + f.left(level, i, j) + f.right(level, i, j)
}

func (f *Field) Process() {
	newLevels := make([]uint, len(f.levels))

	for l := 0; l < 200; l++ {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				if i == 2 && j == 2 {
					continue
				}

				count := f.count(l, i, j)

				if l == 0 {
					fmt.Println(i, j, count)
				}

				if f.get(l, i, j) == 1 {
					newLevels[l] = Set(newLevels[l], i, j, count == 1)
				} else {
					newLevels[l] = Set(newLevels[l], i, j, count == 1 || count == 2)
				}
			}
		}
	}

	f.levels = newLevels
}

func part2() {
	pattern := load(input2)
	field := NewField(pattern)

	show(field.levels[0])
	fmt.Println()

	// for i := 0; i < 200; i++ {
	field.Process()
	// }

	show(field.levels[0])
}

func main() {
	// part1()
	part2()
}
