package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Orbits map[string]string

func load(filename string) Orbits {
	orbits := Orbits{}

	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // EOF
		}

		parsedOrbit := strings.Split(line, ")")

		orbits[strings.Trim(parsedOrbit[1], "\n")] = parsedOrbit[0]
	}

	return orbits
}

func calculateOrbitDepth(name string, orbits *Orbits, depths *map[string]int) int {
	if name == "COM" {
		return 0
	}

	if (*depths)[name] != 0 {
		return (*depths)[name]
	}

	depth := 1 + calculateOrbitDepth((*orbits)[name], orbits, depths)

	(*depths)[name] = depth

	return depth
}

func listAllOrbits(name string, orbits *Orbits) []string {
	parent := (*orbits)[name]

	if parent == "COM" {
		return []string{"COM"}
	}

	return append(listAllOrbits(parent, orbits), parent)
}

func main() {
	orbits := load("input.txt")
	depths := map[string]int{}

	fmt.Println(orbits)
	fmt.Println(depths)

	for name, _ := range orbits {
		calculateOrbitDepth(name, &orbits, &depths)
	}

	fmt.Println(depths)

	total := 0

	for _, v := range depths {
		total += v
	}

	fmt.Println(total)

	santaParents := listAllOrbits("SAN", &orbits)
	myParents := listAllOrbits("YOU", &orbits)

	fmt.Println(santaParents)
	fmt.Println(myParents)

	commonParentCount := 0

	for i := 0; i < len(santaParents) || i < len(myParents); i++ {
		if santaParents[i] == myParents[i] {
			commonParentCount += 1
			continue
		} else {
			break
		}
	}

	fmt.Println(commonParentCount)
	fmt.Println(len(santaParents))
	fmt.Println(len(myParents))

	fmt.Println(len(santaParents) - commonParentCount + len(myParents) - commonParentCount)
}
