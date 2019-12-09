package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

const Width = 25
const Height = 6

type Layer [Height][Width]int
type Image []Layer

func load(filename string) *Image {
	data, _ := ioutil.ReadFile(filename)

	img := Image{}

	for i := 0; i < len(data)-1; i += Height * Width {
		layer := Layer{}

		fmt.Println(i)

		for j := 0; j < Height; j++ {
			for k := 0; k < Width; k++ {
				index := i + j*Width + k

				digit, err := strconv.Atoi(string(data[index]))
				if err != nil {
					panic(err)
				}

				layer[j][k] = digit
			}
		}

		fmt.Println()

		img = append(img, layer)
	}

	return &img
}

func (l *Layer) Count(digit int) int {
	result := 0

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			if l[i][j] == digit {
				result++
			}
		}
	}

	return result
}

func PrintLayer(l Layer) {
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			if l[i][j] == 2 {
				fmt.Printf(" ")
			} else if l[i][j] == 0 {
				fmt.Printf("_")
			} else if l[i][j] == 1 {
				fmt.Printf("#")
			}
		}
		fmt.Println("")
	}
}

func MergeLayers(img Image) Layer {
	layer := Layer{}

	// Set every pixel to transparent
	for j := 0; j < Height; j++ {
		for k := 0; k < Width; k++ {
			layer[j][k] = 2
		}
	}

	// merge
	for i := 0; i < len(img); i++ {
		fmt.Printf("Layer %d\n", i)
		PrintLayer(layer)
		PrintLayer(img[i])
		fmt.Println()

		for j := 0; j < Height; j++ {
			for k := 0; k < Width; k++ {
				if layer[j][k] == 2 {
					layer[j][k] = img[i][j][k]
				}

			}
		}
	}

	return layer
}

func main() {
	img := load("input.txt")

	minLayer := (*img)[0]
	minZeroCount := math.MaxInt32

	for _, l := range *img {
		zeroCount := l.Count(0)

		if zeroCount <= minZeroCount {
			minZeroCount = zeroCount
			minLayer = l
		}
	}

	fmt.Println(minLayer.Count(1) * minLayer.Count(2))

	PrintLayer(MergeLayers(*img))
}
