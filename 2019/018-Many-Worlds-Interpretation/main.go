package main

import (
	"bufio"
	"fmt"
	"os"
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
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		r := []byte{}
		for _, k := range line[0 : len(line)-1] {
			r = append(r, byte(k))
		}

		m = append(m, r)
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

func show(path []Vec) {
	res := ""

	for i, line := range m {
		for j, c := range line {
			found := false

			for _, pos := range path {
				if pos.X == j && pos.Y == i {
					found = true
					break
				}
			}

			if found {
				res += fmt.Sprintf("\033[41m%s\033[0m", string(c))
			} else {
				res += fmt.Sprint(string(c))
			}
		}

		res += "\n"
	}

	// time.Sleep(1000 * time.Millisecond)

	// cmd := exec.Command("clear")
	// cmd.Stdout = os.Stdout
	// cmd.Run()

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

func at(p Vec) byte {
	return m[p.Y][p.X]
}

func findKeys() map[byte]Vec {
	res := map[byte]Vec{}

	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j] >= 'a' && m[i][j] <= 'z' {
				res[m[i][j]] = Vec{X: j, Y: i}
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

var keys = map[byte]Vec{}

func isVisited(path []Vec, pos Vec) bool {
	for _, p := range path {
		if p == pos {
			return true
		}
	}

	return false
}

func dfs(from Vec, goal Vec, path []Vec) ([]Vec, bool) {
	if from == goal {
		return path, true
	}

	minRes := []Vec{}
	minFound := false

	for _, next := range []Vec{up(from), left(from), right(from), down(from)} {
		if at(next) == '#' {
			continue
		}

		if isVisited(path, next) {
			continue
		}

		newPath, ok := dfs(next, goal, append(path, next))

		if ok {
			if !minFound || len(minRes) > len(newPath) {
				minRes = newPath
				minFound = true
			}
		}
	}

	return minRes, minFound
}

func filterDoors(path []Vec) []byte {
	res := []byte{}

	for _, p := range path {
		if isDoor(at(p)) {
			res = append(res, at(p))
		}
	}

	return res
}

type Path struct {
	steps  []Vec
	length int
	doors  []byte
}

var distances = [27][27]Path{}

func showDistances() {
	fmt.Print("     ")

	for i := 0; i < 27; i++ {
		if i == 0 {
			fmt.Print("@   ")
		} else {
			fmt.Print(string(i+'a'-1) + "   ")
		}
	}

	fmt.Println()

	for i := 0; i < 27; i++ {
		if i == 0 {
			fmt.Print("@ ")
		} else {
			fmt.Print(string(i+'a'-1) + " ")
		}

		for j := 0; j < 27; j++ {
			if i == j || distances[i][j].length == 0 {
				fmt.Print("   -")
				continue
			}

			fmt.Printf("%4d", distances[i][j].length)
		}

		fmt.Println()
	}
}

func showDoors() {
	fmt.Print("     ")

	for i := 0; i < 27; i++ {
		if i == 0 {
			fmt.Print("@   ")
		} else {
			fmt.Print(string(i+'a'-1) + "   ")
		}
	}

	fmt.Println()

	for i := 0; i < 27; i++ {
		if i == 0 {
			fmt.Print("@ ")
		} else {
			fmt.Print(string(i+'a'-1) + " ")
		}

		for j := 0; j < 27; j++ {
			if i == j || len(distances[i][j].doors) == 0 {
				fmt.Print("   -")
				continue
			}

			fmt.Printf("%4d", len(distances[i][j].doors))
		}

		fmt.Println()
	}
}

func hasKeysForDoors(doors []byte, keys []byte) bool {
	for _, d := range doors {
		found := false

		for _, k := range keys {

			// fmt.Println(string(d), string(k), string(byte(k)-'a'+'A'))

			if d == byte(k)-'a'+'A' {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

var queue = []byte{0}

func solve(visited []byte, steps int) ([]byte, int, bool) {
	// fmt.Println(keys)
	if len(visited) == len(keys) {
		return visited, steps, true
	}

	nextQueue := []byte{}

	for _, k := range keys {
		kk := k
		if k > 0 {
			kk = at(k) - byte('a') + 1
		}

		for _, k := range keys {
			if kk == from {
				continue
			}

			found := false
			for _, v := range visited {
				if v == at(k) {
					found = true
					break
				}
			}

			if found {
				continue
			}

			if hasKeysForDoors(path.doors, visited) {
				nextQueue = append(nextQueue, kk)
				res, steps, ok := solve(kk, append(visited, at(k)), steps+path.length)

				if ok {
					if !minFound || minSteps >= steps {
						minSteps = steps
						minRes = res
						minFound = true
					}
				}
			}

		}

	}

	return minRes, minSteps, minFound
}

func main() {
	load("input5.txt")

	pos := findPosition()
	fmt.Println(pos)

	for k1, p1 := range keys {
		path1, _ := dfs(pos, keys[k1], []Vec{pos})

		distances[0][k1-'a'+1] = Path{
			steps:  path1,
			length: len(path1) - 1,
			doors:  filterDoors(path1),
		}

		for k2, p2 := range keys {
			if k1 == k2 {
				continue
			}

			path2, ok := dfs(p1, p2, []Vec{p1})

			if !ok {
				fmt.Println(p1, p2)
				fmt.Println("Path not found")
				os.Exit(1)
			}

			distances[k1-'a'+1][k2-'a'+1] = Path{
				steps:  path2,
				length: len(path2) - 1,
				doors:  filterDoors(path2),
			}

			// if !(path2[0] == p1 && path2[len(path2)-1] == p2) {
			// 	fmt.Println(string(k1), string(k2))
			// 	fmt.Println(path2)

			// 	show(distances[k1-'a'+1][k2-'a'+1].steps)
			// 	os.Exit(1)
			// }
		}
	}

	showDistances()
	showDoors()

	steps, l, _ := solve(0, []byte{}, 0)

	for _, s := range steps {
		fmt.Println(string(s))
	}

	fmt.Println(l)

	show([]Vec{})

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
