package main

import "fmt"

var (
	start = 171309
	end   = 643603
)

func meets_criteria(i int) bool {
	d1 := (i / 1) % 10
	d2 := (i / 10) % 10
	d3 := (i / 100) % 10
	d4 := (i / 1000) % 10
	d5 := (i / 10000) % 10
	d6 := (i / 100000) % 10

	adj := (d1 == d2 && d2 != d3) ||
		(d2 == d3 && d2 != d1 && d3 != d4) ||
		(d3 == d4 && d3 != d2 && d4 != d5) ||
		(d4 == d5 && d4 != d3 && d5 != d6) ||
		(d5 == d6 && d5 != d4)

	inc := (d6 <= d5) && (d5 <= d4) && (d4 <= d3) && (d3 <= d2) && (d2 <= d1)

	return adj && inc
}

func main() {
	count := 0

	for i := start; i < end; i++ {
		if meets_criteria(i) {
			fmt.Println(i)
			count++
		}
	}

	fmt.Println(count)
}
