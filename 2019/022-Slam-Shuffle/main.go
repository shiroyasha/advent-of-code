package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DealWithIncrement = "deal-with-increment"
	DealIntoNewStack  = "deal-into-new-stack"
	Cut               = "cut"
)

type Action struct {
	name  string
	value int
}

func load(filename string) []Action {
	result := []Action{}

	file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		action := Action{}

		parts := strings.Split(line[:len(line)-1], " ")

		if parts[0] == "cut" {
			v, _ := strconv.Atoi(parts[1])

			action.name = Cut
			action.value = v
		}

		if parts[0] == "deal" && parts[1] == "into" {
			action.name = DealIntoNewStack
		}

		if parts[0] == "deal" && parts[1] == "with" {
			v, _ := strconv.Atoi(parts[3])

			action.name = DealWithIncrement
			action.value = v
		}

		result = append(result, action)
	}

	return result
}

func reverse(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func dealWithIncr(a []int, inc int) []int {
	res := make([]int, len(a))

	for i, el := range a {
		res[i*inc%len(a)] = el
	}

	return res
}

func part1() int {
	actions := load("input1.txt")
	cardNum := 10007

	cards := []int{}
	for i := 0; i < cardNum; i++ {
		cards = append(cards, i)
	}

	for _, a := range actions {
		for i, v := range cards {
			if v == 2019 {
				fmt.Println(i)
				break
			}
		}

		if a.name == DealIntoNewStack {
			reverse(cards)
		}

		if a.name == Cut {
			if a.value >= 0 {
				cards = append(cards[a.value:], cards[:a.value]...)
			} else {
				cards = append(cards[len(cards)+a.value:], cards[:len(cards)+a.value]...)
			}
		}

		if a.name == DealWithIncrement {
			cards = dealWithIncr(cards, a.value)
		}
	}

	res := 0

	for i, v := range cards {
		if v == 2019 {
			res = i
		}
	}

	return res
}

func part2() int {
	actions := load("input1.txt")
	cards := 10007

	res := 2019

	for _, a := range actions {
		fmt.Println(res)

		if a.name == DealIntoNewStack {
			res = cards - res - 1
		}

		if a.name == Cut {
			val := 0
			if a.value >= 0 {
				val = a.value
			} else {
				val = cards + a.value
			}

			if val > res {
				res += cards - val
			} else {
				res -= val
			}
		}

		if a.name == DealWithIncrement {
			res = res * a.value % cards
		}
	}

	return res
}

func main() {
	res1 := part1()
	fmt.Println(res1)

	fmt.Println("------------")

	res2 := part2()
	fmt.Println(res2)
}
