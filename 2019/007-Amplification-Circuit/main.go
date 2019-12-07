package main

import (
	"fmt"
)

const (
	OpcodeAdd         = 1
	OpcodeMultiply    = 2
	OpcodeGetInput    = 3
	OpcodeWriteOutput = 4
	OpcodeJumpIfTrue  = 5
	OpcodeJumpIfFalse = 6
	OpcodeLessThan    = 7
	OpcodeEquals      = 8
	OpcodeHalt        = 99
)

const (
	InputModePosition  = 0
	InputModeImmidiate = 1
)

type Process struct {
	memory       []int
	position     int
	output       []int
	input        []int
	inputPointer int

	halted bool
}

func NewProcess(code []int, input []int) *Process {
	memory := make([]int, len(code))

	copy(memory, code)

	return &Process{
		memory:       memory,
		position:     0,
		input:        input,
		inputPointer: 0,
		halted:       false,
	}
}

//
// Read & write to program memory
//
func (p *Process) Read(position int) (int, error) {
	if position >= len(p.memory) || position < 0 {
		return 0, fmt.Errorf("Index %d out of range", position)
	}

	return p.memory[position], nil
}

func (p *Process) Write(position int, value int) error {
	if position >= len(p.memory) || position < 0 {
		return fmt.Errorf("Index out of range")
	}

	p.memory[position] = value

	return nil
}

func (p *Process) DumpMemory() {
	for i, v := range p.memory {
		if i == p.position {
			fmt.Printf("[%d] ", v)
		} else {
			fmt.Printf("%d ", v)
		}
	}

	fmt.Print("\n")
}

func (p *Process) LoadParam(position int, mode int) (int, error) {
	switch mode {
	case InputModePosition:
		pointer, err := p.Read(position)
		if err != nil {
			return 0, err
		}

		value, err := p.Read(pointer)
		if err != nil {
			return 0, err
		}

		return value, nil
	case InputModeImmidiate:
		value, err := p.Read(position)
		if err != nil {
			return 0, err
		}

		return value, nil
	default:
		return 0, fmt.Errorf("Unknonwn input mode")
	}
}

//
// Run program until halt or error.
//
func (p *Process) Run() error {
	operation, err := p.Read(p.position)
	if err != nil {
		return err
	}

	// fmt.Printf("%d: ", p.position)

	instruction := operation % 100

	param1Mode := (operation / 100) % 10
	param2Mode := (operation / 1000) % 10

	switch instruction {
	case OpcodeAdd:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		resultPointer, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		p.Write(resultPointer, value1+value2)

		// fmt.Printf("ADD %d %d %d ", value1, value2, resultPointer)

		p.position += 4

		return p.Run()
	case OpcodeMultiply:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		resultPointer, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		p.Write(resultPointer, value1*value2)

		// fmt.Printf("MUL %d %d %d ", value1, value2, resultPointer)

		p.position += 4

		return p.Run()
	case OpcodeGetInput:
		pointer, err := p.Read(p.position + 1)
		if err != nil {
			return err
		}

		if p.inputPointer == len(p.input) {
			// no input, program needs to complete
			fmt.Printf("Waiting for input")
			return nil
		}

		p.Write(pointer, p.input[p.inputPointer])
		p.inputPointer++

		// fmt.Printf("GET %d ", pointer)

		p.position += 2

		return p.Run()
	case OpcodeWriteOutput:
		value, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		p.output = append(p.output, value)

		// fmt.Printf("WRT %d ", value)

		p.position += 2

		fmt.Printf("Output written")
		return nil
	case OpcodeJumpIfTrue:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		// fmt.Printf("JMPT %d ", value1)

		if value1 != 0 {
			value2, err := p.LoadParam(p.position+2, param2Mode)
			if err != nil {
				return err
			}

			// fmt.Printf("%d true", value2)

			p.position = value2
		} else {
			// fmt.Printf("false")

			p.position += 3
		}
		return p.Run()

	case OpcodeJumpIfFalse:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		// fmt.Printf("JMPF %d ", value1)

		if value1 == 0 {
			value2, err := p.LoadParam(p.position+2, param2Mode)
			if err != nil {
				return err
			}

			// fmt.Printf("%d true", value2)

			p.position = value2
		} else {
			// fmt.Printf("false ")

			p.position += 3
		}

		return p.Run()
	case OpcodeLessThan:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		pos, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		fmt.Printf("LT %d %d %d ", value1, value2, pos)

		if value1 < value2 {
			// fmt.Printf("true")

			p.Write(pos, 1)
		} else {
			// fmt.Printf("false")

			p.Write(pos, 0)
		}

		p.position += 4

		return p.Run()
	case OpcodeEquals:
		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		pos, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		fmt.Printf("EQ %d %d %d ", value1, value2, pos)

		if value1 == value2 {
			// fmt.Printf("true")

			p.Write(pos, 1)
		} else {
			// fmt.Printf("false")

			p.Write(pos, 0)
		}

		p.position += 4

		return p.Run()
	case OpcodeHalt:
		// fmt.Println("halt")

		p.halted = true

		return nil
	default:
		return fmt.Errorf("Unknonwn opcode %d", operation)
	}

	panic("How we got here?")
}

func allHalted(processes []*Process) bool {
	for _, p := range processes {
		if !p.halted {
			return false
		}
	}

	return true
}

func signalStrength(code []int, phaseSettings []int) int {
	amplifiers := []*Process{
		NewProcess(code, []int{phaseSettings[0]}),
		NewProcess(code, []int{phaseSettings[1]}),
		NewProcess(code, []int{phaseSettings[2]}),
		NewProcess(code, []int{phaseSettings[3]}),
		NewProcess(code, []int{phaseSettings[4]}),
	}

	input := 0

	for {
		if allHalted(amplifiers) {
			break
		}

		for i, a := range amplifiers {
			fmt.Printf("Running ampf %d\n", i)

			a.input = append(a.input, input)

			err := a.Run()
			if err != nil {
				panic(err)
			}

			input = a.output[len(a.output)-1]
		}
	}

	return input
}

func variationsRec(arr []int, i int, f func([]int)) {
	if i > len(arr) {
		f(arr)
		return
	}

	variationsRec(arr, i+1, f)

	for j := i + 1; j < len(arr); j++ {
		arr[i], arr[j] = arr[j], arr[i]

		variationsRec(arr, i+1, f)

		arr[i], arr[j] = arr[j], arr[i]
	}
}

func variations(arr []int, f func([]int)) {
	variationsRec(arr, 0, f)
}

func main() {
	code := []int{3, 8, 1001, 8, 10, 8, 105, 1, 0, 0, 21, 46, 59, 80, 105, 122, 203, 284, 365, 446, 99999, 3, 9, 102, 3, 9, 9, 1001, 9, 5, 9, 102, 2, 9, 9, 1001, 9, 3, 9, 102, 4, 9, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 101, 5, 9, 9, 1002, 9, 3, 9, 1001, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 99, 3, 9, 1002, 9, 4, 9, 1001, 9, 2, 9, 102, 4, 9, 9, 101, 3, 9, 9, 102, 2, 9, 9, 4, 9, 99, 3, 9, 102, 5, 9, 9, 101, 4, 9, 9, 102, 3, 9, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 99, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 99, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99}

	maxStrength := 0

	variations([]int{5, 6, 7, 8, 9}, func(phaseSettings []int) {
		fmt.Println(phaseSettings)

		strength := signalStrength(code, phaseSettings)

		if strength > maxStrength {
			maxStrength = strength
		}

		fmt.Println(strength)
	})

	fmt.Println(maxStrength)
}
