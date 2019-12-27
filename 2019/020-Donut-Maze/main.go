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

type Location struct {
	Pos   Vec
	Layer int
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

				inner := pos.X > 3 && pos.X <= len(line)-3 && pos.Y >= 3 && pos.Y <= len(m)-3

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

func show(loc Location) {
	res := ""

	for i, line := range m {
		for j, c := range line {
			if j == loc.Pos.X && i == loc.Pos.Y {
				res += fmt.Sprint("\033[31m@\033[0m")
			} else {
				found := false
				for depth := 50; depth >= 0; depth-- {
					if visited[Location{Pos: Vec{X: j, Y: i}, Layer: depth}] {
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

	time.Sleep(20 * time.Millisecond)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println(res)
	fmt.Println(loc)
}

var visited = map[Location]bool{}

func loadPortalName(currentPos, nextPos Vec) string {
	n1 := nextPos
	n2 := Vec{
		X: nextPos.X + nextPos.X - currentPos.X,
		Y: nextPos.Y + nextPos.Y - currentPos.Y,
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

	return name
}

func resolveNextLoc(current Location, nextPos Vec) (Location, bool, bool) {
	if at(nextPos) == '#' || at(nextPos) == ' ' {
		return current, false, false
	}

	if isPortal(nextPos) {
		name := loadPortalName(current.Pos, nextPos)

		fmt.Println(name)

		if name == "AA" {
			return current, false, false
		}

		if name == "ZZ" {
			return current, true, true
		}

		portal := portals[name]

		if portal.warpOuter == current.Pos && current.Layer == 0 {
			return current, false, false
		}

		if portal.warpInner == current.Pos {
			return Location{
				Layer: current.Layer + 1,
				Pos:   portal.warpOuter,
			}, true, false
		}

		if portal.warpOuter == current.Pos {
			return Location{
				Layer: current.Layer - 1,
				Pos:   portal.warpInner,
			}, true, false
		}
	}

	return Location{Pos: nextPos, Layer: current.Layer}, true, false
}

func solve(loc Location, steps int) (int, bool) {
	show(loc)

	visited[loc] = true

	minOk := false
	min := 1000000001

	for _, nextPos := range []Vec{up(loc.Pos), down(loc.Pos), left(loc.Pos), right(loc.Pos)} {
		nextLoc, isResolved, isDestination := resolveNextLoc(loc, nextPos)

		if visited[nextLoc] || !isResolved {
			continue
		}

		if isDestination {
			minOk = true
			min = steps
			continue
		}

		steps, ok := solve(nextLoc, steps+1)
		if ok {
			if steps < min {
				min = steps
				minOk = true
			}
		}
	}

	visited[loc] = false

	return min, minOk
}

func main() {
	load("input4.txt")

	pos := portals["AA"].warpOuter
	loc := Location{Pos: pos, Layer: 0}

	fmt.Println(pos)
	for k, p := range portals {
		fmt.Printf("%s %+v\n", k, p)
	}

	steps, ok := solve(loc, 0)

	fmt.Println(steps, ok)
}
