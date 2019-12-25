package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Vec struct {
	X, Y int
}

var m [][]byte
var doorPositions = map[byte]Vec{}

func load(filename string) {
	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break // EOF
		}

		m = append(m, line)
	}

	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if isDoor(m[i][j]) {
				doorPositions[m[i][j]] = Vec{X: j, Y: i}
			}
		}
	}

	keys = findKeys()
}

func show(pos Vec) {
	res := ""

	for i, line := range m {
		for j, c := range line {
			if pos.X == j && pos.Y == i {
				res += fmt.Sprintf("\033[31m%s\033[0m", string('@'))
			} else {
				res += fmt.Sprint(string(c))
			}
		}

		res += "\n"
	}

	time.Sleep(10 * time.Millisecond)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println(res)
}

func findPosition() Vec {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j] == '@' {
				return Vec{X: j, Y: i}
			}
		}
	}

	panic("No entrance")
}

func findKeys() map[byte]bool {
	res := map[byte]bool{}

	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j] >= 'a' && m[i][j] <= 'z' {
				res[m[i][j]] = false
			}
		}
	}

	return res
}

func showKeys(k map[byte]bool) {
	for k, v := range k {
		fmt.Printf("%s:%t ", string(k), v)
	}
	fmt.Println()
}

func get(pos Vec) byte {
	return m[pos.Y][pos.X]
}

func allKeysFound(k map[byte]bool) bool {
	for _, v := range k {
		if v == false {
			return false
		}
	}
	return true
}

func isDoor(v byte) bool {
	return v >= 'A' && v <= 'Z'
}

func isKey(v byte) bool {
	return v >= 'a' && v <= 'z'
}

func doorForKey(v byte) byte {
	return (v - 'a') + 'A'
}

func set(pos Vec, v byte) {
	m[pos.Y][pos.X] = v
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

var minRes = []Vec{}
var minSet = false
var steps = []Vec{}
var keys = map[byte]bool{}

func solve(pos Vec, visits map[Vec]bool) {
	if allKeysFound(keys) {
		if minSet == false || len(minRes) > len(steps)-1 {
			fmt.Println("found ", len(steps))
			minRes = append([]Vec{}, steps[0:len(steps)-1]...)
			minSet = true
		}
		return
	}

	// show(pos)
	// for k, v := range keys {
	// 	if v {
	// 		fmt.Print(string(k), " ")
	// 	}
	// }

	v := get(pos)

	if v == '#' || isDoor(v) {
		return
	}

	if isKey(v) {
		set(pos, '.')
		keys[v] = true

		door := doorForKey(v)

		unlockPos, ok := doorPositions[door]
		if ok {
			set(unlockPos, '.')
		}

		defer func() {
			if ok {
				set(unlockPos, door)
			}
		}()
		defer func() { keys[v] = false }()

		visits = map[Vec]bool{}
	}

	if get(pos) != '.' {
		panic("how did we end up here. PANIC!!! " + string(get(pos)))
	}

	visits[pos] = true

	for _, p := range []Vec{up(pos), down(pos), left(pos), right(pos)} {
		if visits[p] || get(p) == '#' || isDoor(get(p)) {
			continue
		}

		old := get(p)
		steps = append(steps, p)

		solve(p, visits)

		steps = steps[0 : len(steps)-1]
		set(p, old)
	}
}

func main() {
	load("input4.txt")

	pos := findPosition()

	fmt.Println(pos)

	// show(pos)
	// showKeys(keys)

	set(pos, '.')

	solve(pos, map[Vec]bool{})

	fmt.Println(minRes)
	fmt.Println(len(minRes))

	// m = [][]byte{}

	// load("input2.txt")

	// for _, s := range shortest {
	// 	set(s, '.')
	// 	show(s)

	// 	if isKey(get(s)) {
	// 		unlock(get(s))
	// 	}
	// }

	// fmt.Println(pos)
	// fmt.Println(shortest)
	// fmt.Println(len(shortest))
}
