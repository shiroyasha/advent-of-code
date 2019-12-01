package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadModuleMasses(filename string) []int {
	result := []int{}

	inputFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	check(err)
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break // EOF
		}

		mass, err := strconv.ParseInt(line[0:len(line)-1], 10, 64)
		check(err)

		result = append(result, int(mass))
	}

	return result
}

//
// The fuel requeirment based on the mass that needs to be transported.
//
func fuelForMass(mass int) int {
	return mass/3 - 2
}

//
// Fuel itself requires fuel just like a module - take its mass, divide by
// three, round down, and subtract 2. However, that fuel also requires fuel,
// and that fuel requires fuel, and so on. Any mass that would require negative
// fuel should instead be treated as if it requires zero fuel; the remaining
// mass, if any, is instead handled by wishing really hard, which has no mass
// and is outside the scope of this calculation.
//
func fuelForFuel(fuel int) int {
	result := fuelForMass(fuel)

	if result <= 0 {
		return 0
	}

	return result + fuelForFuel(result)
}

func main() {
	fuelForModules := 0
	totalFuel := 0

	moduleMasses := loadModuleMasses("input1.txt")

	for _, m := range moduleMasses {
		fuel := fuelForMass(m)

		fuelForModules += fuel
		totalFuel += fuel + fuelForFuel(fuel)
	}

	fmt.Printf("Fuel for modules: %d\n", fuelForModules)
	fmt.Printf("Total Fuel: %d\n", totalFuel)
}
