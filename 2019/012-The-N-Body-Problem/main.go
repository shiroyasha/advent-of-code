package main

import (
	"fmt"
)

type Body struct {
	pos, vel Vec
}

type Vec struct {
	X, Y, Z int
}

func Add(v1, v2 Vec) Vec {
	return Vec{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

func Energy(v Vec) int {
	return Abs(v.X) + Abs(v.Y) + Abs(v.Z)
}

func gravityOnAxis(a, b int) int {
	if a < b {
		return 1
	}

	if a > b {
		return -1
	}

	return 0
}

func gravity(p1, p2 Vec) Vec {
	return Vec{
		X: gravityOnAxis(p1.X, p2.X),
		Y: gravityOnAxis(p1.Y, p2.Y),
		Z: gravityOnAxis(p1.Z, p2.Z),
	}
}

func part1() {
	steps := 1000
	bodies := []Body{
		Body{
			pos: Vec{X: 5, Y: 4, Z: 4},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -11, Y: -11, Z: -3},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: 0, Y: 7, Z: 0},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -13, Y: 2, Z: 10},
			vel: Vec{},
		},
	}

	for step := 0; step < steps; step++ {
		// apply gravity
		for i, b1 := range bodies {
			acc := Vec{}

			for _, b2 := range bodies {
				acc = Add(acc, gravity(b1.pos, b2.pos))
			}

			bodies[i].vel = Add(bodies[i].vel, acc)
		}

		// apply velocity

		for i, b := range bodies {
			bodies[i].pos = Add(b.pos, b.vel)
		}

		fmt.Println(bodies[0].pos.X)
	}

	total := 0

	for _, b := range bodies {
		total += Energy(b.pos) * Energy(b.vel)
	}

	fmt.Println(total)
}

func part2() {
	fmt.Println("---- Part 2 ----")

	start := []Body{
		Body{
			pos: Vec{X: 5, Y: 4, Z: 4},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -11, Y: -11, Z: -3},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: 0, Y: 7, Z: 0},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -13, Y: 2, Z: 10},
			vel: Vec{},
		},
	}

	bodies := []Body{
		Body{
			pos: Vec{X: 5, Y: 4, Z: 4},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -11, Y: -11, Z: -3},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: 0, Y: 7, Z: 0},
			vel: Vec{},
		},
		Body{
			pos: Vec{X: -13, Y: 2, Z: 10},
			vel: Vec{},
		},
	}

	xDone := false
	yDone := false
	zDone := false

	xCycle := 0
	yCycle := 0
	zCycle := 0

	for {
		if !xDone {
			xCycle++
		}

		if !yDone {
			yCycle++
		}

		if !zDone {
			zCycle++
		}

		// apply gravity
		for i, b1 := range bodies {
			acc := Vec{}

			for _, b2 := range bodies {
				acc = Add(acc, gravity(b1.pos, b2.pos))
			}

			bodies[i].vel = Add(bodies[i].vel, acc)
		}

		// apply velocity

		for i, b := range bodies {
			bodies[i].pos = Add(b.pos, b.vel)
		}

		if !xDone {
			ok := true

			for i, _ := range bodies {
				if ok && (bodies[i].vel.X != start[i].vel.X || bodies[i].pos.X != start[i].pos.X) {
					ok = false
				}
			}

			if ok {
				xDone = true
			}
		}

		if !yDone {
			ok := true

			for i, _ := range bodies {
				if ok && (bodies[i].vel.Y != start[i].vel.Y || bodies[i].pos.Y != start[i].pos.Y) {
					ok = false
				}
			}

			if ok {
				yDone = true
			}
		}

		if !zDone {
			ok := true

			for i, _ := range bodies {
				if ok && (bodies[i].vel.Z != start[i].vel.Z || bodies[i].pos.Z != start[i].pos.Z) {
					ok = false
				}
			}

			if ok {
				zDone = true
			}
		}

		if xDone && yDone && zDone {
			break
		}
	}

	result := lcm(xCycle, lcm(yCycle, zCycle))

	fmt.Println(result)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	if a == 0 {
		return b
	}

	return gcd(b%a, a)
}

func main() {
	part1()
	part2()
}
