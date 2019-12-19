package main

import (
	"fmt"
)

const debug = false

const (
	OpcodeAdd                = 1
	OpcodeMultiply           = 2
	OpcodeGetInput           = 3
	OpcodeWriteOutput        = 4
	OpcodeJumpIfTrue         = 5
	OpcodeJumpIfFalse        = 6
	OpcodeLessThan           = 7
	OpcodeEquals             = 8
	OpcodeAdjustRelativeBase = 9
	OpcodeHalt               = 99
)

const (
	InputModePosition  = 0
	InputModeImmidiate = 1
	InputModeRelative  = 2
)

type Process struct {
	memory          []int
	position        int
	output          []int
	input           []int
	inputPointer    int
	outputPointer   int
	relativeBase    int
	waitingForInput bool

	halted bool
}

func NewProcess(code []int, input []int) *Process {
	memory := make([]int, len(code)+900000)

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

func (p *Process) Write(position int, value int, mode int) error {
	if position >= len(p.memory) || position < 0 {
		return fmt.Errorf("Index out of range")
	}

	switch mode {
	case InputModePosition:
		if debug {
			fmt.Printf(" A %d = %d", p.memory[position], value)
		}
		p.memory[p.memory[position]] = value
	case InputModeRelative:
		if debug {
			fmt.Printf(" R %d = %d", p.memory[position]+p.relativeBase, value)
		}
		p.memory[p.memory[position]+p.relativeBase] = value
	}

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

	case InputModeRelative:
		pointer, err := p.Read(position)
		if err != nil {
			return 0, err
		}

		value, err := p.Read(pointer + p.relativeBase)
		if err != nil {
			return 0, err
		}

		return value, nil
	default:
		return 0, fmt.Errorf("Unknonwn input mode")
	}
}

func (p *Process) Debug(name string, length int) {
	if !debug {
		return
	}

	fmt.Printf("%4s ", name)

	mode1 := (p.memory[p.position] / 100) % 10
	mode2 := (p.memory[p.position] / 1000) % 10

	if length >= 1 {
		fmt.Printf(" %10d ", p.memory[p.position+1])
	} else {
		fmt.Printf(" %9s- ", "")
	}

	if length >= 2 {
		fmt.Printf(" %10d ", p.memory[p.position+2])
	} else {
		fmt.Printf(" %9s- ", "")
	}

	if length >= 3 {
		fmt.Printf(" %10d ", p.memory[p.position+3])
	} else {
		fmt.Printf(" %9s- ", "")
	}

	fmt.Printf(" | ")

	if length >= 1 {
		v, _ := p.LoadParam(p.position+1, mode1)
		fmt.Printf(" %16d ", v)
	} else {
		fmt.Printf(" %15s- ", "")
	}

	if length >= 2 {
		v, _ := p.LoadParam(p.position+2, mode2)

		fmt.Printf(" %16d ", v)
	} else {
		fmt.Printf(" %15s- ", "")
	}

	if length >= 3 {
		v, _ := p.Read(p.position + 3)

		fmt.Printf(" %16d ", v)
	} else {
		fmt.Printf(" %15s- ", "")
	}
}

//
// Run program until halt or error.
//
func (p *Process) RunTilInterupt() error {
	operation, err := p.Read(p.position)
	if err != nil {
		return err
	}

	if debug {
		fmt.Println()
		fmt.Printf("%4d (r %4d): ", p.position, p.relativeBase)
	}

	instruction := operation % 100

	param1Mode := (operation / 100) % 10
	param2Mode := (operation / 1000) % 10
	param3Mode := (operation / 10000) % 10

	if debug {
		fmt.Printf("[%d %d %d %2d] ", param1Mode, param2Mode, param3Mode, instruction)
	}

	switch instruction {
	case OpcodeAdd:
		p.Debug("ADD", 3)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		p.Write(p.position+3, value1+value2, param3Mode)

		p.position += 4

		return p.RunTilInterupt()
	case OpcodeMultiply:
		p.Debug("MUL", 3)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		p.Write(p.position+3, value1*value2, param3Mode)

		p.position += 4

		return p.RunTilInterupt()
	case OpcodeGetInput:
		p.Debug("GET", 1)

		p.waitingForInput = true

		if p.inputPointer == len(p.input) {
			// no input, program needs to complete
			if debug {
				fmt.Printf("Waiting for input")
			}
			return nil
		}

		p.waitingForInput = false

		p.Write(p.position+1, p.input[p.inputPointer], param1Mode)

		p.inputPointer++

		p.position += 2

		return p.RunTilInterupt()
	case OpcodeWriteOutput:
		p.Debug("WRT", 1)

		value, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		p.output = append(p.output, value)

		p.position += 2

		// fmt.Printf("Output written")
		return nil
	case OpcodeJumpIfTrue:
		p.Debug("JIT", 2)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		if value1 != 0 {
			value2, err := p.LoadParam(p.position+2, param2Mode)
			if err != nil {
				return err
			}

			if debug {
				fmt.Printf("true")
			}

			p.position = value2
		} else {
			if debug {
				fmt.Printf("false")
			}

			p.position += 3
		}
		return p.RunTilInterupt()

	case OpcodeJumpIfFalse:
		p.Debug("JIF", 2)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		if value1 == 0 {
			value2, err := p.LoadParam(p.position+2, param2Mode)
			if err != nil {
				return err
			}

			if debug {
				fmt.Printf("true")
			}

			p.position = value2
		} else {
			if debug {
				fmt.Printf("false ")
			}

			p.position += 3
		}

		return p.RunTilInterupt()
	case OpcodeLessThan:
		p.Debug("LT", 3)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		if value1 < value2 {
			if debug {
				fmt.Printf("true")
			}

			p.Write(p.position+3, 1, param3Mode)
		} else {
			if debug {
				fmt.Printf("false")
			}

			p.Write(p.position+3, 0, param3Mode)
		}

		p.position += 4

		return p.RunTilInterupt()
	case OpcodeEquals:
		p.Debug("EQL", 3)

		value1, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		value2, err := p.LoadParam(p.position+2, param2Mode)
		if err != nil {
			return err
		}

		if value1 == value2 {
			if debug {
				fmt.Printf("true")
			}

			p.Write(p.position+3, 1, param3Mode)
		} else {
			if debug {
				fmt.Printf("false")
			}

			p.Write(p.position+3, 0, param3Mode)
		}

		p.position += 4

		return p.RunTilInterupt()

	case OpcodeAdjustRelativeBase:
		p.Debug("ADJ", 1)

		value, err := p.LoadParam(p.position+1, param1Mode)
		if err != nil {
			return err
		}

		p.relativeBase += value

		p.position += 2

		return p.RunTilInterupt()
	case OpcodeHalt:
		if debug {
			fmt.Println("halt")
		}

		p.halted = true

		return nil
	default:
		return fmt.Errorf("Unknonwn opcode %d", operation)
	}

	// fmt.Println(operation)
	// fmt.Println(instruction)
	panic("How we got here?")
}

func (p *Process) RunTilInputNeeded() error {
	for {
		err := p.RunTilInterupt()

		if err != nil {
			return err
		}

		if p.waitingForInput || p.halted {
			return nil
		}
	}
}

func (p *Process) Run() error {
	for {
		err := p.RunTilInterupt()

		if err != nil {
			return err
		}

		if p.halted {
			return nil
		}
	}
}

func (p *Process) AddInput(val int) {
	p.input = append(p.input, val)
}

func (p *Process) NextOutput() int {
	res := p.output[p.outputPointer]
	p.outputPointer++
	return res
}

type Point struct {
	X, Y int
}

const screenX = 42
const screenY = 24

var score = 0
var screen [screenY][screenX]int

func paddleX() int {
	for y := 0; y < screenY; y++ {
		for x := 0; x < screenX; x++ {
			if screen[y][x] == 3 {
				return x
			}
		}
	}

	return -1
}

func ballX() int {
	for y := 0; y < screenY; y++ {
		for x := 0; x < screenX; x++ {
			if screen[y][x] == 4 {
				return x
			}
		}
	}

	return -1
}

func blockCount() int {
	count := 0

	for y := 0; y < screenY; y++ {
		for x := 0; x < screenX; x++ {
			if screen[y][x] == 2 {
				count++
			}
		}
	}

	return count
}

func render() {
	print("\033[H\033[2J")

	for y := 0; y < screenY; y++ {
		for x := 0; x < screenX; x++ {
			if screen[y][x] == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(screen[y][x])
			}
		}

		fmt.Println()
	}
}

func update(p *Process) {
	for i := 0; i < len(p.output); i += 3 {
		x := p.output[i]
		y := p.output[i+1]
		t := p.output[i+2]

		if x == -1 && y == 0 {
			score = t
			continue
		}

		screen[y][x] = t
	}
}

func main() {
	code := []int{2, 380, 379, 385, 1008, 2655, 455702, 381, 1005, 381, 12, 99, 109, 2656, 1101, 0, 0, 383, 1101, 0, 0, 382, 20102, 1, 382, 1, 21002, 383, 1, 2, 21101, 37, 0, 0, 1105, 1, 578, 4, 382, 4, 383, 204, 1, 1001, 382, 1, 382, 1007, 382, 42, 381, 1005, 381, 22, 1001, 383, 1, 383, 1007, 383, 24, 381, 1005, 381, 18, 1006, 385, 69, 99, 104, -1, 104, 0, 4, 386, 3, 384, 1007, 384, 0, 381, 1005, 381, 94, 107, 0, 384, 381, 1005, 381, 108, 1106, 0, 161, 107, 1, 392, 381, 1006, 381, 161, 1101, -1, 0, 384, 1106, 0, 119, 1007, 392, 40, 381, 1006, 381, 161, 1102, 1, 1, 384, 21002, 392, 1, 1, 21102, 1, 22, 2, 21102, 1, 0, 3, 21101, 138, 0, 0, 1106, 0, 549, 1, 392, 384, 392, 21001, 392, 0, 1, 21102, 22, 1, 2, 21102, 3, 1, 3, 21101, 0, 161, 0, 1106, 0, 549, 1102, 0, 1, 384, 20001, 388, 390, 1, 20102, 1, 389, 2, 21102, 180, 1, 0, 1105, 1, 578, 1206, 1, 213, 1208, 1, 2, 381, 1006, 381, 205, 20001, 388, 390, 1, 20101, 0, 389, 2, 21101, 0, 205, 0, 1106, 0, 393, 1002, 390, -1, 390, 1102, 1, 1, 384, 21002, 388, 1, 1, 20001, 389, 391, 2, 21101, 0, 228, 0, 1106, 0, 578, 1206, 1, 261, 1208, 1, 2, 381, 1006, 381, 253, 21002, 388, 1, 1, 20001, 389, 391, 2, 21102, 253, 1, 0, 1105, 1, 393, 1002, 391, -1, 391, 1102, 1, 1, 384, 1005, 384, 161, 20001, 388, 390, 1, 20001, 389, 391, 2, 21101, 0, 279, 0, 1106, 0, 578, 1206, 1, 316, 1208, 1, 2, 381, 1006, 381, 304, 20001, 388, 390, 1, 20001, 389, 391, 2, 21102, 304, 1, 0, 1105, 1, 393, 1002, 390, -1, 390, 1002, 391, -1, 391, 1102, 1, 1, 384, 1005, 384, 161, 20102, 1, 388, 1, 21001, 389, 0, 2, 21101, 0, 0, 3, 21101, 0, 338, 0, 1106, 0, 549, 1, 388, 390, 388, 1, 389, 391, 389, 20101, 0, 388, 1, 20102, 1, 389, 2, 21101, 4, 0, 3, 21102, 365, 1, 0, 1106, 0, 549, 1007, 389, 23, 381, 1005, 381, 75, 104, -1, 104, 0, 104, 0, 99, 0, 1, 0, 0, 0, 0, 0, 0, 268, 19, 19, 1, 1, 21, 109, 3, 21201, -2, 0, 1, 21202, -1, 1, 2, 21102, 0, 1, 3, 21101, 0, 414, 0, 1105, 1, 549, 22101, 0, -2, 1, 22102, 1, -1, 2, 21101, 0, 429, 0, 1105, 1, 601, 1202, 1, 1, 435, 1, 386, 0, 386, 104, -1, 104, 0, 4, 386, 1001, 387, -1, 387, 1005, 387, 451, 99, 109, -3, 2105, 1, 0, 109, 8, 22202, -7, -6, -3, 22201, -3, -5, -3, 21202, -4, 64, -2, 2207, -3, -2, 381, 1005, 381, 492, 21202, -2, -1, -1, 22201, -3, -1, -3, 2207, -3, -2, 381, 1006, 381, 481, 21202, -4, 8, -2, 2207, -3, -2, 381, 1005, 381, 518, 21202, -2, -1, -1, 22201, -3, -1, -3, 2207, -3, -2, 381, 1006, 381, 507, 2207, -3, -4, 381, 1005, 381, 540, 21202, -4, -1, -1, 22201, -3, -1, -3, 2207, -3, -4, 381, 1006, 381, 529, 22102, 1, -3, -7, 109, -8, 2106, 0, 0, 109, 4, 1202, -2, 42, 566, 201, -3, 566, 566, 101, 639, 566, 566, 2101, 0, -1, 0, 204, -3, 204, -2, 204, -1, 109, -4, 2106, 0, 0, 109, 3, 1202, -1, 42, 593, 201, -2, 593, 593, 101, 639, 593, 593, 21001, 0, 0, -2, 109, -3, 2105, 1, 0, 109, 3, 22102, 24, -2, 1, 22201, 1, -1, 1, 21101, 0, 509, 2, 21102, 684, 1, 3, 21102, 1, 1008, 4, 21102, 630, 1, 0, 1106, 0, 456, 21201, 1, 1647, -2, 109, -3, 2106, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 2, 2, 0, 0, 0, 0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 2, 0, 0, 2, 2, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 1, 1, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 2, 0, 0, 2, 2, 2, 0, 2, 0, 2, 0, 1, 1, 0, 2, 2, 2, 0, 0, 2, 0, 2, 0, 2, 2, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 2, 0, 2, 2, 0, 2, 2, 2, 0, 0, 0, 2, 0, 2, 2, 2, 0, 1, 1, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 2, 2, 2, 0, 2, 2, 2, 0, 2, 0, 2, 2, 0, 0, 0, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 1, 1, 0, 2, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 2, 0, 2, 2, 2, 2, 2, 2, 0, 2, 0, 0, 2, 0, 2, 0, 0, 2, 2, 2, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 2, 2, 2, 2, 2, 0, 2, 2, 2, 2, 2, 0, 0, 0, 2, 0, 2, 0, 0, 2, 0, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 2, 2, 2, 0, 2, 0, 0, 1, 1, 0, 2, 0, 0, 2, 2, 2, 2, 0, 2, 2, 2, 0, 0, 0, 0, 2, 0, 2, 0, 0, 2, 0, 0, 2, 2, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 2, 0, 0, 0, 1, 1, 0, 2, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 0, 2, 0, 0, 0, 0, 2, 2, 2, 2, 0, 2, 0, 0, 2, 2, 0, 0, 2, 0, 0, 0, 0, 1, 1, 0, 0, 2, 0, 0, 0, 2, 0, 2, 2, 2, 0, 2, 2, 0, 2, 2, 2, 0, 0, 0, 2, 0, 2, 0, 2, 2, 0, 0, 2, 0, 0, 0, 0, 2, 0, 2, 2, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 0, 0, 2, 2, 0, 2, 0, 2, 0, 2, 2, 2, 2, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 1, 1, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 0, 0, 0, 2, 0, 2, 0, 0, 2, 2, 0, 0, 2, 2, 0, 2, 0, 0, 2, 0, 0, 2, 0, 1, 1, 0, 2, 0, 0, 0, 2, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 2, 0, 0, 2, 2, 0, 2, 0, 0, 2, 0, 0, 2, 2, 2, 0, 1, 1, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 0, 0, 2, 0, 2, 2, 0, 2, 2, 0, 2, 0, 2, 0, 0, 0, 0, 0, 2, 0, 2, 2, 0, 0, 0, 2, 2, 2, 0, 1, 1, 0, 2, 2, 2, 0, 0, 0, 2, 0, 2, 2, 0, 0, 0, 2, 2, 0, 2, 0, 0, 0, 2, 2, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 2, 2, 0, 0, 0, 0, 1, 1, 0, 2, 2, 0, 2, 0, 0, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 2, 2, 0, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 23, 82, 82, 16, 37, 71, 32, 87, 51, 93, 33, 83, 22, 21, 23, 36, 43, 97, 16, 24, 33, 77, 54, 2, 88, 59, 72, 36, 26, 90, 26, 4, 4, 44, 42, 14, 5, 40, 27, 7, 27, 96, 27, 74, 43, 17, 90, 6, 85, 69, 21, 28, 82, 82, 81, 53, 95, 14, 84, 70, 92, 51, 29, 86, 83, 44, 37, 36, 54, 77, 1, 26, 33, 92, 46, 74, 43, 10, 96, 73, 31, 32, 22, 66, 14, 89, 2, 72, 97, 3, 16, 22, 31, 24, 90, 87, 18, 18, 42, 55, 82, 38, 2, 64, 38, 22, 49, 39, 32, 23, 14, 58, 15, 24, 65, 7, 28, 88, 15, 81, 20, 18, 70, 5, 98, 56, 60, 9, 47, 94, 7, 51, 18, 90, 27, 74, 50, 45, 81, 86, 73, 75, 89, 56, 63, 34, 15, 72, 48, 86, 77, 66, 47, 91, 18, 89, 25, 51, 41, 2, 57, 52, 84, 84, 44, 76, 7, 15, 97, 56, 59, 50, 73, 94, 81, 7, 4, 95, 32, 82, 97, 36, 60, 38, 5, 51, 60, 65, 51, 27, 45, 5, 82, 35, 7, 30, 63, 44, 9, 95, 29, 70, 88, 63, 48, 56, 12, 40, 44, 28, 94, 25, 48, 72, 28, 95, 83, 46, 48, 67, 42, 23, 23, 76, 34, 25, 84, 40, 39, 69, 6, 40, 28, 42, 15, 19, 92, 9, 91, 94, 22, 51, 31, 19, 39, 42, 60, 63, 16, 29, 46, 69, 52, 7, 79, 59, 33, 90, 93, 61, 59, 9, 98, 1, 13, 24, 74, 70, 35, 12, 50, 54, 67, 83, 18, 88, 52, 49, 40, 19, 59, 54, 33, 62, 66, 82, 65, 63, 29, 93, 14, 7, 57, 56, 87, 52, 41, 28, 46, 14, 70, 69, 94, 25, 88, 59, 7, 45, 18, 73, 11, 41, 20, 42, 7, 25, 36, 88, 76, 42, 57, 65, 84, 21, 12, 71, 25, 94, 38, 5, 71, 60, 61, 92, 24, 32, 18, 36, 12, 74, 57, 95, 59, 30, 94, 88, 30, 30, 9, 96, 25, 80, 88, 27, 89, 89, 48, 84, 23, 11, 50, 45, 53, 81, 18, 57, 94, 50, 57, 26, 87, 33, 3, 50, 71, 96, 71, 89, 49, 29, 45, 6, 74, 32, 98, 23, 27, 7, 92, 29, 93, 82, 84, 95, 98, 1, 74, 59, 10, 92, 63, 60, 54, 34, 70, 4, 60, 59, 7, 30, 70, 8, 53, 52, 23, 46, 7, 26, 88, 40, 51, 77, 12, 32, 33, 34, 46, 79, 4, 33, 33, 10, 16, 7, 23, 90, 74, 90, 93, 78, 6, 21, 40, 77, 64, 76, 74, 58, 7, 26, 18, 74, 90, 82, 40, 68, 60, 18, 45, 16, 59, 96, 48, 7, 96, 49, 60, 48, 88, 42, 63, 30, 18, 8, 96, 88, 36, 38, 82, 96, 17, 72, 76, 23, 98, 45, 74, 26, 42, 69, 11, 56, 26, 59, 67, 33, 98, 62, 73, 7, 59, 22, 17, 48, 89, 14, 1, 47, 28, 43, 95, 91, 33, 62, 15, 77, 81, 29, 6, 81, 20, 55, 1, 51, 19, 40, 25, 52, 43, 19, 91, 47, 59, 21, 88, 73, 80, 65, 62, 57, 19, 80, 1, 40, 74, 33, 30, 95, 73, 68, 92, 26, 86, 22, 12, 33, 30, 23, 14, 79, 52, 42, 2, 61, 32, 3, 55, 10, 10, 4, 71, 4, 6, 22, 36, 39, 8, 14, 11, 92, 61, 74, 12, 15, 16, 77, 50, 8, 7, 1, 38, 40, 11, 87, 11, 96, 52, 74, 69, 34, 63, 48, 45, 92, 71, 60, 6, 58, 47, 23, 25, 64, 50, 98, 48, 80, 27, 76, 31, 66, 91, 3, 74, 9, 59, 97, 45, 98, 18, 74, 45, 9, 7, 29, 97, 64, 57, 54, 19, 61, 37, 41, 14, 62, 55, 92, 79, 16, 85, 53, 78, 85, 93, 30, 94, 5, 51, 34, 25, 64, 21, 21, 79, 16, 59, 12, 68, 50, 39, 59, 62, 17, 40, 51, 42, 26, 51, 60, 87, 21, 37, 97, 45, 23, 43, 27, 7, 9, 25, 48, 54, 37, 45, 34, 7, 58, 86, 8, 48, 91, 88, 56, 94, 7, 80, 80, 15, 83, 91, 23, 92, 23, 29, 36, 62, 50, 2, 45, 9, 94, 96, 93, 60, 18, 96, 83, 40, 13, 19, 28, 69, 26, 66, 75, 36, 98, 35, 39, 70, 58, 67, 72, 78, 59, 57, 60, 18, 60, 41, 97, 94, 39, 11, 18, 70, 63, 24, 5, 19, 41, 92, 27, 88, 81, 28, 37, 36, 92, 51, 23, 32, 69, 95, 8, 66, 67, 59, 49, 31, 16, 65, 17, 23, 57, 71, 75, 20, 63, 36, 62, 32, 82, 26, 73, 57, 93, 69, 27, 20, 91, 72, 23, 44, 86, 94, 59, 23, 49, 15, 7, 4, 69, 64, 59, 77, 37, 50, 42, 64, 88, 3, 4, 23, 47, 60, 46, 72, 22, 78, 46, 12, 18, 30, 18, 19, 74, 80, 93, 43, 10, 73, 15, 59, 47, 37, 53, 16, 57, 43, 72, 81, 4, 55, 40, 33, 14, 16, 85, 61, 90, 72, 40, 79, 96, 24, 94, 75, 14, 59, 7, 76, 52, 13, 87, 53, 10, 87, 95, 4, 51, 13, 89, 68, 34, 68, 15, 31, 60, 64, 21, 41, 84, 12, 90, 6, 5, 85, 77, 94, 10, 8, 18, 61, 39, 80, 90, 78, 13, 16, 13, 36, 48, 28, 71, 91, 90, 35, 20, 60, 98, 44, 18, 88, 69, 22, 71, 27, 79, 54, 38, 25, 8, 6, 94, 36, 3, 57, 10, 58, 92, 6, 88, 62, 19, 67, 47, 79, 95, 71, 6, 68, 37, 16, 28, 89, 34, 72, 56, 65, 11, 35, 10, 83, 24, 51, 41, 40, 31, 12, 84, 68, 41, 44, 56, 73, 46, 59, 93, 98, 3, 71, 12, 90, 26, 80, 88, 97, 64, 18, 24, 75, 34, 85, 53, 39, 62, 69, 58, 13, 17, 91, 53, 89, 58, 34, 87, 64, 43, 455702}
	p := NewProcess(code, []int{})

	// starting to play

	for {
		p.RunTilInputNeeded()

		update(p)
		render()
		fmt.Println(score)

		bx := ballX()
		px := paddleX()

		if bx < px {
			p.AddInput(-1)
		}

		if bx > px {
			p.AddInput(1)
		}

		if bx == px {
			p.AddInput(0)
		}

		c := blockCount()

		fmt.Println(c)

		if c == 0 {
			break
		}
	}

	fmt.Println(score)
}
