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
}

func show(pos Vec) {
	time.Sleep(200 * time.Millisecond)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for i, line := range m {
		for j, c := range line {
			if pos.X == j && pos.Y == i {
				fmt.Printf("\033[31m%s\033[0m", string('@'))
			} else {
				fmt.Print(string(c))
			}
		}

		fmt.Println()
	}
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

func keyForDoor(v byte) byte {
	return v + ('a' - 'A')
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

func dup(k map[byte]bool) map[byte]bool {
	newMap := make(map[byte]bool)

	for key, value := range k {
		newMap[key] = value
	}

	return newMap
}

func solve(pos Vec, keys map[byte]bool, visits map[Vec]bool) ([]Vec, bool) {
	if allKeysFound(keys) {
		return []Vec{}, true
	}

	v := get(pos)

	if v == '#' {
		return []Vec{}, false
	}

	// show(pos)
	// fmt.Println(visits)

	if isDoor(v) {
		keyNeeded := keyForDoor(v)

		if v, ok := keys[keyNeeded]; ok && v == true {
			set(pos, '.')
			visits = map[Vec]bool{pos: true}
		} else {
			return []Vec{}, false
		}
	}

	if isKey(v) {
		keys = dup(keys)
		keys[v] = true
		set(pos, '.')

		visits = map[Vec]bool{pos: true}
	}

	if get(pos) != '.' {
		panic("how did we end up here. PANIC!!! " + string(get(pos)))
	}

	minOk := false
	minRest := []Vec{Vec{-1, -1}}
	visits[pos] = true

	for _, p := range []Vec{up(pos), down(pos), left(pos), right(pos)} {
		if visits[p] {
			continue
		}

		rest, ok := solve(p, keys, visits)

		if ok {
			if minOk == false {
				minRest = rest
				minOk = true
			} else {
				if len(minRest) > len(rest) {
					minRest = rest
				}
			}
		}
	}

	return append([]Vec{pos}, minRest...), minOk
}

func main() {
	load("input2.txt")

	pos := findPosition()
	keys := findKeys()

	fmt.Println(pos)

	// show(pos)
	// showKeys(keys)

	set(pos, '.')

	steps, _ := solve(pos, keys, map[Vec]bool{pos: true})

	fmt.Println(steps)
	fmt.Println(len(steps) - 1)

	m = [][]byte{}

	load("input2.txt")

	for _, s := range steps {
		set(s, '.')
		show(s)
	}

	fmt.Println(len(steps) - 1)
}
