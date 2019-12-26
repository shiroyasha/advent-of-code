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

var m [][]byte
var layer = 0

type Portal struct {
	innerPos  Vec
	warpInner Vec
	outerPos  Vec
	warpOuter Vec
}

func is(pos Vec, v byte) bool {
	if pos.Y < 0 || pos.X < 0 {
		return false
	}

	if pos.Y >= len(m) || pos.X >= len(m[0]) {
		return false
	}

	return m[pos.Y][pos.X] == v
}

func at(pos Vec) byte {
	return m[pos.Y][pos.X]
}

func isPortal(pos Vec) bool {
	if pos.Y < 0 || pos.X < 0 {
		return false
	}

	if pos.Y >= len(m) || pos.X >= len(m[0]) {
		return false
	}

	return m[pos.Y][pos.X] >= 'A' && m[pos.Y][pos.X] <= 'Z'
}

var portals = map[string]Portal{}

func load(filename string) {
	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		bytes := []byte{}

		for _, c := range line[0 : len(line)-1] {
			bytes = append(bytes, byte(c))
		}

		m = append(m, bytes)
	}

	for i, line := range m {
		for j, c := range line {
			if c >= 'A' && c <= 'Z' {
				pos := Vec{X: j, Y: i}

				if !isPortal(pos) {
					continue
				}

				name := ""
				warpPos := Vec{}

				pUp := up(pos)
				pDown := down(pos)
				pLeft := left(pos)
				pRight := right(pos)

				if !is(pDown, '.') && !is(pUp, '.') && !is(pLeft, '.') && !is(pRight, '.') {
					continue
				}

				if isPortal(pUp) {
					name = string(at(pUp)) + string(c)
					warpPos = pDown
				}

				if isPortal(pDown) {
					name = string(c) + string(at(pDown))
					warpPos = pUp
				}

				if isPortal(pLeft) {
					name = string(at(pLeft)) + string(c)
					warpPos = pRight
				}

				if isPortal(pRight) {
					name = string(c) + string(at(pRight))
					warpPos = pLeft
				}

				inner := pos.X > 3 && pos.X <= len(line) && pos.Y >= 3 && pos.Y <= len(m)

				portal := portals[name]
				if inner {
					portal.innerPos = pos
					portal.warpInner = warpPos
				} else {
					portal.outerPos = pos
					portal.warpOuter = warpPos
				}

				portals[name] = portal
			}
		}
	}

	fmt.Println(portals)
}

func show(pos Vec) {
	res := ""

	for i, line := range m {
		for j, c := range line {
			if j == pos.X && i == pos.Y {
				res += fmt.Sprint("\033[31m@\033[0m")
			} else {
				found := false
				for depth := len(visited) - 1; depth >= 0; depth-- {
					if visited[depth][Vec{X: j, Y: i}] {
						res += fmt.Sprintf("\033[4%dm%s\033[0m", depth+1, string(c))
						found = true
						break
					}
				}
				if !found {
					res += fmt.Sprint(string(c))
				}
			}
		}

		res += "\n"
	}

	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println(res)
	fmt.Println("layer", layer)
}

var visited = []map[Vec]bool{}

func solve(pos Vec, steps int) (int, bool) {
	if len(visited) <= layer {
		visited = append(visited, map[Vec]bool{})
	}

	show(pos)

	visited[layer][pos] = true

	minOk := false
	min := 1000000001

	for _, nextPos := range []Vec{up(pos), down(pos), left(pos), right(pos)} {
		if isPortal(nextPos) {
			n1 := nextPos
			n2 := Vec{
				X: nextPos.X + nextPos.X - pos.X,
				Y: nextPos.Y + nextPos.Y - pos.Y,
			}

			name := ""

			if n1.X == n2.X {
				if n1.Y < n2.Y {
					name = string(at(n1)) + string(at(n2))
				} else {
					name = string(at(n2)) + string(at(n1))
				}
			}

			if n1.Y == n2.Y {
				if n1.X < n2.X {
					name = string(at(n1)) + string(at(n2))
				} else {
					name = string(at(n2)) + string(at(n1))
				}
			}

			fmt.Println(n1, n2)
			fmt.Println(string(at(n1)), string(at(n2)))
			fmt.Println(name)
			fmt.Println(portals[name])

			if name == "AA" {
				continue
			}

			if name == "ZZ" {
				minOk = true
				min = steps
				continue
			}
			if portals[name].warpInner == pos {
				layer += 1
				if len(visited) <= layer {
					visited = append(visited, map[Vec]bool{})
				}

				nextPos = portals[name].warpOuter

				visited[layer-1][portals[name].outerPos] = true
				visited[layer][portals[name].outerPos] = true
			}

			if layer > 0 && !visited[layer][pos] && portals[name].warpOuter == pos {
				layer -= 1
				nextPos = portals[name].warpInner
			}
		}

		fmt.Println(layer)

		if visited[layer][nextPos] || at(nextPos) == ' ' || at(nextPos) == '#' {
			continue
		}

		steps, ok := solve(nextPos, steps+1)
		if ok {
			if steps < min {
				min = steps
				minOk = true
			}
		}
	}

	visited[layer][pos] = false

	return min, minOk
}

func main() {
	load("input4.txt")

	pos := portals["AA"].warpInner

	fmt.Println(pos)
	for k, p := range portals {
		fmt.Printf("%s %+v\n", k, p)
	}

	steps, ok := solve(pos, 0)

	fmt.Println(steps, ok)
}
