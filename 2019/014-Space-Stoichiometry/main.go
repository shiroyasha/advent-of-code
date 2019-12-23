package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	name     string
	quantity int
}

type Reaction struct {
	output      Node
	requirments []Node
}

type NanoFactory struct {
	reactions map[string]Reaction
	stock     map[string]int
	ores      int
}

func parseNode(raw string) Node {
	parts := strings.Split(strings.TrimSpace(raw), " ")
	name := parts[1]

	quantity, _ := strconv.Atoi(parts[0])

	return Node{name: name, quantity: quantity}
}

func load(filename string) NanoFactory {
	result := NanoFactory{
		reactions: map[string]Reaction{},
		stock:     map[string]int{},
	}

	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break // EOF
		}

		parts := strings.Split(line, "=>")

		requirments := []Node{}
		for _, l := range strings.Split(parts[0], ",") {
			requirments = append(requirments, parseNode(l))
		}

		output := parseNode(parts[1])

		result.reactions[output.name] = Reaction{output: output, requirments: requirments}
		result.stock[output.name] = 0
	}

	return result
}

func (f *NanoFactory) create(name string, quantity int) {
	outQuantity := f.reactions[name].output.quantity
	inStock := f.stock[name]
	reqs := f.reactions[name].requirments
	multiple := int(math.Ceil(math.Max(float64(quantity-inStock), 0) / float64(outQuantity)))

	for _, r := range reqs {
		if r.name == "ORE" {
			f.ores += multiple * r.quantity
			continue
		}

		f.create(r.name, multiple*r.quantity)
	}

	inStock += multiple * outQuantity

	f.stock[name] = inStock - quantity
}

func (f *NanoFactory) clear() {
	f.ores = 0
	f.stock = map[string]int{}
}

func main() {
	factory := load("input.txt")

	factory.create("FUEL", 1)

	min := 1000000000000 / factory.ores
	max := min * 10000

	for {
		f := (max + min) / 2

		factory.clear()
		factory.create("FUEL", f)

		if max == min || max-1 == min {
			break
		}

		if factory.ores > 1000000000000 {
			max = f
		} else {
			min = f
		}

		fmt.Println(min, max, f, factory.ores)
	}

	result := min

	factory.clear()
	factory.create("FUEL", result)

	fmt.Println(factory.ores)
}
