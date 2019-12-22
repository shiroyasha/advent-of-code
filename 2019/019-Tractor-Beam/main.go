package main

import (
	"fmt"
)

const debug = true

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
	memory        []int
	position      int
	output        []int
	input         []int
	inputPointer  int
	outputPointer int
	relativeBase  int

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

		if p.inputPointer == len(p.input) {
			// no input, program needs to complete
			// fmt.Printf("Waiting for input")
			return nil
		}

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

		return p.RunTilInterupt()
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

// --------------------------------------------

type Pos struct {
	X, Y int
}

func (p *Pos) Add(p2 Pos) {
	p.X += p2.X
	p.Y += p2.Y
}

func (p *Pos) Substract(p2 Pos) {
	p.X -= p2.X
	p.Y -= p2.Y
}

type Robot struct {
	process *Process
	pos     Pos
}

// func check(directions []int) int {
// 	code := []int{3, 1033, 1008, 1033, 1, 1032, 1005, 1032, 31, 1008, 1033, 2, 1032, 1005, 1032, 58, 1008, 1033, 3, 1032, 1005, 1032, 81, 1008, 1033, 4, 1032, 1005, 1032, 104, 99, 101, 0, 1034, 1039, 1001, 1036, 0, 1041, 1001, 1035, -1, 1040, 1008, 1038, 0, 1043, 102, -1, 1043, 1032, 1, 1037, 1032, 1042, 1106, 0, 124, 1001, 1034, 0, 1039, 1001, 1036, 0, 1041, 1001, 1035, 1, 1040, 1008, 1038, 0, 1043, 1, 1037, 1038, 1042, 1106, 0, 124, 1001, 1034, -1, 1039, 1008, 1036, 0, 1041, 102, 1, 1035, 1040, 101, 0, 1038, 1043, 102, 1, 1037, 1042, 1105, 1, 124, 1001, 1034, 1, 1039, 1008, 1036, 0, 1041, 101, 0, 1035, 1040, 1001, 1038, 0, 1043, 101, 0, 1037, 1042, 1006, 1039, 217, 1006, 1040, 217, 1008, 1039, 40, 1032, 1005, 1032, 217, 1008, 1040, 40, 1032, 1005, 1032, 217, 1008, 1039, 1, 1032, 1006, 1032, 165, 1008, 1040, 3, 1032, 1006, 1032, 165, 1101, 0, 2, 1044, 1105, 1, 224, 2, 1041, 1043, 1032, 1006, 1032, 179, 1102, 1, 1, 1044, 1106, 0, 224, 1, 1041, 1043, 1032, 1006, 1032, 217, 1, 1042, 1043, 1032, 1001, 1032, -1, 1032, 1002, 1032, 39, 1032, 1, 1032, 1039, 1032, 101, -1, 1032, 1032, 101, 252, 1032, 211, 1007, 0, 45, 1044, 1105, 1, 224, 1101, 0, 0, 1044, 1106, 0, 224, 1006, 1044, 247, 1002, 1039, 1, 1034, 1002, 1040, 1, 1035, 1001, 1041, 0, 1036, 1002, 1043, 1, 1038, 102, 1, 1042, 1037, 4, 1044, 1106, 0, 0, 7, 39, 95, 7, 98, 8, 11, 47, 17, 33, 19, 4, 29, 41, 87, 34, 59, 22, 75, 5, 1, 46, 41, 29, 32, 11, 55, 25, 53, 41, 77, 27, 52, 33, 41, 65, 72, 24, 43, 83, 72, 3, 14, 92, 2, 43, 82, 30, 87, 19, 94, 47, 91, 10, 8, 67, 24, 4, 68, 85, 63, 4, 93, 29, 55, 34, 23, 65, 40, 3, 36, 90, 57, 97, 37, 2, 65, 8, 1, 16, 83, 93, 67, 44, 71, 97, 27, 70, 76, 20, 40, 90, 36, 73, 27, 89, 57, 13, 66, 37, 95, 76, 26, 84, 33, 48, 34, 86, 85, 30, 81, 6, 61, 33, 83, 84, 22, 21, 67, 27, 11, 49, 28, 69, 41, 60, 98, 6, 69, 41, 54, 82, 18, 37, 65, 10, 42, 47, 41, 2, 72, 16, 66, 39, 93, 37, 2, 41, 52, 49, 20, 78, 30, 7, 38, 15, 40, 81, 21, 14, 82, 44, 48, 7, 96, 33, 36, 70, 52, 18, 71, 1, 81, 66, 47, 1, 38, 78, 80, 38, 63, 53, 80, 16, 58, 55, 93, 31, 89, 36, 36, 78, 65, 71, 34, 83, 4, 55, 60, 29, 10, 30, 84, 15, 59, 31, 96, 16, 21, 58, 26, 38, 35, 58, 50, 16, 46, 25, 26, 82, 59, 12, 11, 98, 4, 17, 42, 66, 83, 72, 23, 14, 92, 22, 9, 5, 87, 5, 79, 85, 19, 87, 71, 28, 61, 32, 56, 92, 56, 19, 78, 94, 39, 24, 73, 58, 28, 37, 81, 11, 99, 25, 46, 73, 44, 5, 22, 41, 76, 55, 84, 31, 16, 36, 65, 84, 40, 29, 81, 66, 16, 94, 23, 54, 23, 29, 51, 20, 25, 23, 69, 44, 23, 18, 99, 80, 55, 39, 10, 71, 7, 33, 63, 94, 93, 62, 26, 35, 25, 50, 61, 39, 84, 38, 54, 43, 56, 23, 67, 17, 70, 34, 23, 90, 93, 24, 46, 60, 31, 46, 33, 53, 81, 10, 62, 23, 89, 86, 43, 39, 73, 82, 38, 9, 61, 42, 66, 68, 30, 28, 95, 4, 25, 54, 22, 21, 80, 32, 61, 13, 6, 66, 47, 59, 4, 31, 59, 17, 87, 72, 30, 72, 51, 30, 30, 62, 43, 53, 88, 42, 48, 13, 21, 80, 8, 30, 61, 14, 77, 22, 27, 60, 87, 30, 65, 14, 33, 76, 67, 9, 95, 26, 84, 40, 21, 52, 11, 86, 23, 30, 86, 57, 28, 6, 69, 4, 11, 63, 21, 2, 65, 51, 39, 58, 82, 16, 51, 96, 23, 3, 44, 21, 62, 31, 38, 47, 73, 30, 29, 94, 24, 14, 88, 1, 51, 72, 42, 57, 48, 63, 33, 95, 78, 15, 17, 68, 64, 61, 10, 31, 58, 68, 36, 15, 52, 19, 13, 26, 38, 72, 41, 66, 15, 56, 88, 18, 98, 87, 15, 43, 89, 96, 3, 94, 55, 25, 26, 27, 6, 48, 3, 29, 90, 88, 6, 18, 29, 88, 90, 43, 3, 81, 61, 16, 31, 93, 42, 26, 46, 31, 56, 66, 17, 76, 37, 15, 50, 33, 81, 16, 10, 83, 87, 37, 39, 92, 80, 62, 6, 59, 77, 9, 32, 91, 61, 97, 24, 44, 62, 61, 11, 36, 94, 59, 54, 34, 23, 67, 18, 86, 31, 39, 77, 73, 44, 67, 27, 57, 5, 54, 65, 29, 21, 81, 2, 65, 39, 24, 82, 6, 55, 33, 97, 72, 35, 16, 85, 19, 28, 57, 94, 21, 15, 86, 5, 52, 53, 39, 69, 20, 32, 52, 5, 86, 95, 44, 47, 77, 9, 57, 14, 62, 49, 54, 7, 70, 29, 16, 42, 87, 99, 30, 36, 67, 68, 14, 42, 73, 4, 87, 97, 39, 61, 18, 11, 39, 77, 83, 17, 83, 27, 1, 72, 30, 21, 95, 38, 35, 96, 15, 78, 27, 66, 40, 4, 95, 90, 94, 4, 20, 63, 71, 19, 54, 11, 28, 96, 46, 13, 42, 94, 84, 9, 22, 79, 37, 14, 50, 13, 58, 64, 90, 30, 69, 18, 20, 90, 4, 21, 31, 95, 88, 22, 81, 36, 20, 11, 82, 59, 95, 38, 43, 72, 3, 78, 38, 33, 62, 48, 36, 22, 16, 3, 87, 53, 91, 37, 12, 19, 49, 18, 25, 14, 67, 78, 79, 9, 70, 88, 34, 98, 38, 8, 90, 98, 56, 13, 26, 34, 82, 77, 40, 97, 82, 63, 32, 57, 26, 58, 53, 29, 56, 3, 62, 17, 78, 67, 69, 33, 49, 62, 47, 36, 60, 9, 81, 12, 96, 6, 78, 86, 98, 34, 70, 41, 87, 86, 47, 15, 46, 36, 49, 20, 76, 31, 48, 1, 68, 19, 96, 0, 0, 21, 21, 1, 10, 1, 0, 0, 0, 0, 0, 0}
// 	p := NewProcess(code, directions)

// 	p.RunTilInterupt()

// 	return p.output[len(p.output)-1]
// }

var movements = []Pos{
	Pos{X: 0, Y: -1},
	Pos{X: 0, Y: 1},
	Pos{X: -1, Y: 0},
	Pos{X: 1, Y: 0},
}

var oposites = []int{2, 1, 4, 3}

type Visits []Pos

func visited(visits Visits, pos Pos) bool {
	for _, v := range visits {
		if v == pos {
			return true
		}
	}

	return false
}

func move(pos Pos, dir int) Pos {
	switch dir {
	case 1:
		return Pos{X: pos.X, Y: pos.Y - 1}
	case 2:
		return Pos{X: pos.X, Y: pos.Y + 1}
	case 3:
		return Pos{X: pos.X - 1, Y: pos.Y}
	case 4:
		return Pos{X: pos.X + 1, Y: pos.Y}
	}

	panic("unknown direction")
}

const size = 6

var m = [size][size]int{}

func draw() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			switch m[i][j] {
			case 0:
				fmt.Print(".")
			case 1:
				fmt.Print("#")
			}
		}

		fmt.Println()
	}

}

func check(pos Pos) int {
	code := []int{109, 424, 203, 1, 21101, 0, 11, 0, 1106, 0, 282, 21102, 18, 1, 0, 1106, 0, 259, 2101, 0, 1, 221, 203, 1, 21102, 1, 31, 0, 1105, 1, 282, 21101, 0, 38, 0, 1106, 0, 259, 20102, 1, 23, 2, 22101, 0, 1, 3, 21102, 1, 1, 1, 21101, 57, 0, 0, 1105, 1, 303, 2102, 1, 1, 222, 20101, 0, 221, 3, 21002, 221, 1, 2, 21101, 0, 259, 1, 21102, 1, 80, 0, 1105, 1, 225, 21102, 125, 1, 2, 21102, 1, 91, 0, 1106, 0, 303, 2101, 0, 1, 223, 21002, 222, 1, 4, 21102, 1, 259, 3, 21102, 225, 1, 2, 21102, 225, 1, 1, 21101, 0, 118, 0, 1106, 0, 225, 20102, 1, 222, 3, 21101, 0, 69, 2, 21102, 1, 133, 0, 1106, 0, 303, 21202, 1, -1, 1, 22001, 223, 1, 1, 21102, 148, 1, 0, 1106, 0, 259, 1201, 1, 0, 223, 20101, 0, 221, 4, 21001, 222, 0, 3, 21102, 1, 22, 2, 1001, 132, -2, 224, 1002, 224, 2, 224, 1001, 224, 3, 224, 1002, 132, -1, 132, 1, 224, 132, 224, 21001, 224, 1, 1, 21102, 195, 1, 0, 106, 0, 108, 20207, 1, 223, 2, 20101, 0, 23, 1, 21102, -1, 1, 3, 21101, 0, 214, 0, 1105, 1, 303, 22101, 1, 1, 1, 204, 1, 99, 0, 0, 0, 0, 109, 5, 1202, -4, 1, 249, 21202, -3, 1, 1, 22102, 1, -2, 2, 21201, -1, 0, 3, 21101, 250, 0, 0, 1106, 0, 225, 22102, 1, 1, -4, 109, -5, 2105, 1, 0, 109, 3, 22107, 0, -2, -1, 21202, -1, 2, -1, 21201, -1, -1, -1, 22202, -1, -2, -2, 109, -3, 2106, 0, 0, 109, 3, 21207, -2, 0, -1, 1206, -1, 294, 104, 0, 99, 22101, 0, -2, -2, 109, -3, 2106, 0, 0, 109, 5, 22207, -3, -4, -1, 1206, -1, 346, 22201, -4, -3, -4, 21202, -3, -1, -1, 22201, -4, -1, 2, 21202, 2, -1, -1, 22201, -4, -1, 1, 22102, 1, -2, 3, 21101, 0, 343, 0, 1106, 0, 303, 1105, 1, 415, 22207, -2, -3, -1, 1206, -1, 387, 22201, -3, -2, -3, 21202, -2, -1, -1, 22201, -3, -1, 3, 21202, 3, -1, -1, 22201, -3, -1, 2, 22102, 1, -4, 1, 21101, 384, 0, 0, 1106, 0, 303, 1106, 0, 415, 21202, -4, -1, -4, 22201, -4, -3, -4, 22202, -3, -2, -2, 22202, -2, -4, -4, 22202, -3, -2, -3, 21202, -4, -1, -2, 22201, -3, -2, 1, 21202, 1, 1, -4, 109, -5, 2105, 1, 0}
	p := NewProcess(code, []int{pos.X, pos.Y})

	p.Run()

	return p.output[len(p.output)-1]
}

func main() {
	res := 0

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			pos := Pos{X: j, Y: i}

			m[i][j] = check(pos)

			if m[i][j] == 1 {
				res += 1
			}
		}
	}

	draw()

	fmt.Println(res)
}
