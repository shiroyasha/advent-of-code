package main

import (
	"fmt"
	"log"
)

// [1 2 3 4 5  1 2 3 4 5]
//  1 0 1 0 1  0 1 0 1 0

func load(str string) []int {
	result := make([]int, len(str))

	for i, c := range []byte(str) {
		result[i] = int(c - '0')
	}

	return result
}

// [ 1  2   3   4  1   2   3  4  1  2   3  4]
// [ 1  0  -1   0  1   0  -1  0  1  0  -1  0]
// [ 0  1   1   0  0  -1  -1  -1  0  1  0  -1  0]
// [ 0  0   1   1   0  0  -1  -1  -1  0  1  0  -1  0]

// [0 0  1 1 1 0 0 0 -1 -1 -1 0]

// 10

// cycleIndex =

// 11 % 12 = 11

// patternIndex = 11 % 3 rep = 2

// func genPattern(basePattern []int, position int, length int) []int {
// 	result := []int{}

// 	patternIndex := 0

// 	for {
// 		times := position

// 		for {
// 			times -= 1
// 			result = append(result, basePattern[patternIndex])

// 			if times == 0 || len(result) == length+1 {
// 				break
// 			}
// 		}

// 		patternIndex += 1
// 		patternIndex = patternIndex % len(basePattern)

// 		if len(result) == length+1 {
// 			break
// 		}
// 	}

// 	return result[1:len(result)]
// }

func prod(input, pattern []int, repetition int) int {
	result := 0

	// start := time.Now()

	pointer := -1

	for {
		for _, i := range []int{0, 1, 0, -1} {
			if i == 0 {
				pointer += repetition

				continue
			}

			for j := repetition; j > 0; j-- {

				if pointer >= len(input) {
					break
				}

				// fmt.Print(pointer, " ")

				result += input[pointer] * i
				pointer += 1
			}

			if pointer >= len(input) {
				break
			}
		}

		if pointer >= len(input) {
			break
		}
	}

	// fmt.Println("p", time.Now().Sub(start))

	return result
}

func calc(input []int, basePattern []int) {
	prev := make([]int, len(input)*10000)
	next := make([]int, len(input)*10000)

	for i := 0; i < 100; i++ {
		log.Println("Started")

		for j := 0; j >= len(next); j++ {
			result := 0

			for k := 0; k < 10000; k++ {
				result += prod(prev, basePattern, j+1)
			}

			if result < 0 {
				result = -result
			}

			next[j] = result % 10

			fmt.Println(j, len(input))
		}

		input = next
		next = make([]int, len(input)*10000)
		log.Println(input[0:8])
	}

	log.Println(input[0:8])
}

func main() {
	// baseInput := "59769638638635227792873839600619296161830243411826562620803755357641409702942441381982799297881659288888243793321154293102743325904757198668820213885307612900972273311499185929901117664387559657706110034992786489002400852438961738219627639830515185618184324995881914532256988843436511730932141380017180796681870256240757580454505096230610520430997536145341074585637105456401238209187118397046373589766408080120984817035699228422366952628344235542849850709181363703172334788744537357607446322903743644673800140770982283290068502972397970799328249132774293609700245065522290562319955768092155530250003587007804302344866598232236645453817273744027537630"

	// input := load(baseInput)
	pattern := []int{0, 1, 0, -1}

	// calc(input, pattern)

	for repetition := 1; repetition < 50; repetition++ {
		for i := 1; i < 100; i++ {
			cycleIndex := i % (4 * repetition) / repetition

			fmt.Printf("%3d", pattern[cycleIndex])
		}
		fmt.Println()
	}
}
