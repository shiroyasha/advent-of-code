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

func (p *Process) inputString(str string) {
	fmt.Println("Entering: ", str)
	for _, b := range str {
		p.AddInput(int(b))
	}

	p.AddInput(int('\n'))
}

func (p *Process) readOutput() string {
	result := []byte{}

	for _, b := range p.output {
		result = append(result, byte(b))
	}

	return string(result)
}

// ---------------------------------------------------------------------------

func part1() {
	code := []int{3, 62, 1001, 62, 11, 10, 109, 2243, 105, 1, 0, 1555, 2097, 668, 1728, 833, 1425, 2029, 2060, 1798, 1631, 1136, 864, 1988, 1330, 938, 1957, 1169, 969, 2140, 2208, 571, 899, 1660, 1293, 1454, 1485, 1858, 633, 1390, 1691, 802, 1829, 1031, 1596, 998, 1262, 730, 1761, 1520, 2171, 1231, 699, 1062, 765, 604, 1893, 1105, 1200, 1361, 1922, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 64, 1008, 64, -1, 62, 1006, 62, 88, 1006, 61, 170, 1106, 0, 73, 3, 65, 20101, 0, 64, 1, 20102, 1, 66, 2, 21102, 1, 105, 0, 1106, 0, 436, 1201, 1, -1, 64, 1007, 64, 0, 62, 1005, 62, 73, 7, 64, 67, 62, 1006, 62, 73, 1002, 64, 2, 133, 1, 133, 68, 133, 102, 1, 0, 62, 1001, 133, 1, 140, 8, 0, 65, 63, 2, 63, 62, 62, 1005, 62, 73, 1002, 64, 2, 161, 1, 161, 68, 161, 1102, 1, 1, 0, 1001, 161, 1, 169, 1001, 65, 0, 0, 1102, 1, 1, 61, 1102, 1, 0, 63, 7, 63, 67, 62, 1006, 62, 203, 1002, 63, 2, 194, 1, 68, 194, 194, 1006, 0, 73, 1001, 63, 1, 63, 1105, 1, 178, 21102, 1, 210, 0, 106, 0, 69, 1201, 1, 0, 70, 1102, 1, 0, 63, 7, 63, 71, 62, 1006, 62, 250, 1002, 63, 2, 234, 1, 72, 234, 234, 4, 0, 101, 1, 234, 240, 4, 0, 4, 70, 1001, 63, 1, 63, 1105, 1, 218, 1105, 1, 73, 109, 4, 21102, 0, 1, -3, 21101, 0, 0, -2, 20207, -2, 67, -1, 1206, -1, 293, 1202, -2, 2, 283, 101, 1, 283, 283, 1, 68, 283, 283, 22001, 0, -3, -3, 21201, -2, 1, -2, 1106, 0, 263, 22102, 1, -3, -3, 109, -4, 2106, 0, 0, 109, 4, 21102, 1, 1, -3, 21101, 0, 0, -2, 20207, -2, 67, -1, 1206, -1, 342, 1202, -2, 2, 332, 101, 1, 332, 332, 1, 68, 332, 332, 22002, 0, -3, -3, 21201, -2, 1, -2, 1105, 1, 312, 21201, -3, 0, -3, 109, -4, 2105, 1, 0, 109, 1, 101, 1, 68, 358, 21002, 0, 1, 1, 101, 3, 68, 367, 20101, 0, 0, 2, 21102, 1, 376, 0, 1105, 1, 436, 22102, 1, 1, 0, 109, -1, 2105, 1, 0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728, 268435456, 536870912, 1073741824, 2147483648, 4294967296, 8589934592, 17179869184, 34359738368, 68719476736, 137438953472, 274877906944, 549755813888, 1099511627776, 2199023255552, 4398046511104, 8796093022208, 17592186044416, 35184372088832, 70368744177664, 140737488355328, 281474976710656, 562949953421312, 1125899906842624, 109, 8, 21202, -6, 10, -5, 22207, -7, -5, -5, 1205, -5, 521, 21102, 1, 0, -4, 21102, 0, 1, -3, 21101, 0, 51, -2, 21201, -2, -1, -2, 1201, -2, 385, 470, 21002, 0, 1, -1, 21202, -3, 2, -3, 22207, -7, -1, -5, 1205, -5, 496, 21201, -3, 1, -3, 22102, -1, -1, -5, 22201, -7, -5, -7, 22207, -3, -6, -5, 1205, -5, 515, 22102, -1, -6, -5, 22201, -3, -5, -3, 22201, -1, -4, -4, 1205, -2, 461, 1105, 1, 547, 21101, 0, -1, -4, 21202, -6, -1, -6, 21207, -7, 0, -5, 1205, -5, 547, 22201, -7, -6, -7, 21201, -4, 1, -4, 1106, 0, 529, 22101, 0, -4, -7, 109, -8, 2105, 1, 0, 109, 1, 101, 1, 68, 563, 21001, 0, 0, 0, 109, -1, 2106, 0, 0, 1101, 85199, 0, 66, 1102, 1, 2, 67, 1102, 1, 598, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1101, 602, 0, 72, 1105, 1, 73, 0, 0, 0, 0, 49, 96146, 1102, 1, 79357, 66, 1102, 1, 1, 67, 1101, 631, 0, 68, 1102, 556, 1, 69, 1102, 1, 0, 71, 1101, 0, 633, 72, 1106, 0, 73, 1, 1115, 1102, 102593, 1, 66, 1102, 1, 3, 67, 1101, 0, 660, 68, 1101, 302, 0, 69, 1102, 1, 1, 71, 1102, 1, 666, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 69404, 1102, 1, 54751, 66, 1102, 1, 1, 67, 1102, 1, 695, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 697, 0, 72, 1106, 0, 73, 1, -176, 23, 44753, 1102, 84163, 1, 66, 1101, 0, 1, 67, 1101, 726, 0, 68, 1101, 556, 0, 69, 1102, 1, 1, 71, 1102, 1, 728, 72, 1105, 1, 73, 1, 12, 39, 237092, 1101, 0, 10159, 66, 1102, 1, 1, 67, 1101, 0, 757, 68, 1101, 556, 0, 69, 1102, 1, 3, 71, 1102, 759, 1, 72, 1105, 1, 73, 1, 10, 33, 2293, 37, 135068, 12, 125182, 1101, 42197, 0, 66, 1102, 1, 4, 67, 1102, 1, 792, 68, 1102, 253, 1, 69, 1102, 1, 1, 71, 1102, 1, 800, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 10, 91297, 1101, 0, 36373, 66, 1101, 1, 0, 67, 1102, 829, 1, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 831, 0, 72, 1105, 1, 73, 1, 1613, 33, 4586, 1101, 0, 78787, 66, 1101, 1, 0, 67, 1102, 860, 1, 68, 1102, 1, 556, 69, 1101, 1, 0, 71, 1102, 862, 1, 72, 1105, 1, 73, 1, 13, 42, 572971, 1102, 1, 63079, 66, 1101, 0, 1, 67, 1102, 891, 1, 68, 1102, 1, 556, 69, 1101, 3, 0, 71, 1101, 893, 0, 72, 1106, 0, 73, 1, 3, 42, 327412, 7, 79873, 23, 134259, 1102, 1, 17351, 66, 1102, 1, 5, 67, 1102, 926, 1, 68, 1101, 253, 0, 69, 1102, 1, 1, 71, 1102, 1, 936, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 47797, 1101, 1399, 0, 66, 1102, 1, 1, 67, 1102, 1, 965, 68, 1101, 0, 556, 69, 1102, 1, 1, 71, 1102, 1, 967, 72, 1105, 1, 73, 1, 8, 39, 118546, 1101, 0, 28751, 66, 1101, 1, 0, 67, 1102, 1, 996, 68, 1101, 0, 556, 69, 1102, 1, 0, 71, 1101, 0, 998, 72, 1105, 1, 73, 1, 1683, 1101, 0, 79279, 66, 1101, 0, 2, 67, 1102, 1025, 1, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1102, 1029, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 43, 84394, 1101, 37087, 0, 66, 1102, 1, 1, 67, 1102, 1, 1058, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 0, 1060, 72, 1105, 1, 73, 1, 107, 7, 239619, 1101, 0, 81853, 66, 1102, 1, 7, 67, 1101, 0, 1089, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 1103, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 43, 168788, 1102, 103079, 1, 66, 1102, 1, 1, 67, 1101, 0, 1132, 68, 1101, 0, 556, 69, 1101, 0, 1, 71, 1101, 1134, 0, 72, 1105, 1, 73, 1, 9, 23, 89506, 1102, 1, 91297, 66, 1101, 2, 0, 67, 1101, 0, 1163, 68, 1101, 0, 351, 69, 1101, 1, 0, 71, 1101, 0, 1167, 72, 1105, 1, 73, 0, 0, 0, 0, 255, 29569, 1102, 1, 26687, 66, 1101, 1, 0, 67, 1102, 1196, 1, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 1198, 0, 72, 1105, 1, 73, 1, 8867, 42, 245559, 1101, 81869, 0, 66, 1102, 1, 1, 67, 1101, 0, 1227, 68, 1101, 556, 0, 69, 1101, 1, 0, 71, 1101, 1229, 0, 72, 1105, 1, 73, 1, 233, 7, 159746, 1101, 59753, 0, 66, 1102, 1, 1, 67, 1102, 1, 1258, 68, 1101, 0, 556, 69, 1101, 1, 0, 71, 1101, 1260, 0, 72, 1106, 0, 73, 1, 11, 42, 163706, 1102, 1, 78467, 66, 1101, 1, 0, 67, 1102, 1, 1289, 68, 1102, 1, 556, 69, 1101, 1, 0, 71, 1102, 1291, 1, 72, 1105, 1, 73, 1, 131, 27, 102593, 1102, 1, 44753, 66, 1101, 4, 0, 67, 1101, 1320, 0, 68, 1101, 0, 302, 69, 1101, 1, 0, 71, 1102, 1, 1328, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 19, 1373, 1102, 28051, 1, 66, 1101, 1, 0, 67, 1101, 0, 1357, 68, 1101, 0, 556, 69, 1101, 1, 0, 71, 1101, 0, 1359, 72, 1106, 0, 73, 1, 160, 12, 187773, 1102, 1, 50359, 66, 1102, 1, 1, 67, 1102, 1, 1388, 68, 1102, 556, 1, 69, 1101, 0, 0, 71, 1101, 0, 1390, 72, 1105, 1, 73, 1, 1232, 1101, 93251, 0, 66, 1102, 1, 1, 67, 1102, 1, 1417, 68, 1101, 0, 556, 69, 1102, 3, 1, 71, 1102, 1, 1419, 72, 1106, 0, 73, 1, 7, 42, 81853, 38, 278583, 23, 179012, 1102, 1, 53759, 66, 1102, 1, 1, 67, 1101, 1452, 0, 68, 1101, 556, 0, 69, 1102, 0, 1, 71, 1101, 0, 1454, 72, 1106, 0, 73, 1, 1370, 1101, 19597, 0, 66, 1101, 1, 0, 67, 1101, 1481, 0, 68, 1102, 556, 1, 69, 1102, 1, 1, 71, 1102, 1, 1483, 72, 1105, 1, 73, 1, 2311, 27, 205186, 1102, 18253, 1, 66, 1101, 3, 0, 67, 1101, 1512, 0, 68, 1102, 1, 302, 69, 1101, 0, 1, 71, 1102, 1518, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 43, 126591, 1101, 0, 92861, 66, 1102, 1, 3, 67, 1101, 1547, 0, 68, 1102, 1, 302, 69, 1102, 1, 1, 71, 1102, 1, 1553, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 34, 158558, 1102, 29569, 1, 66, 1101, 1, 0, 67, 1101, 1582, 0, 68, 1101, 556, 0, 69, 1102, 6, 1, 71, 1101, 1584, 0, 72, 1106, 0, 73, 1, 20982, 34, 79279, 19, 2746, 19, 4119, 25, 18253, 25, 36506, 25, 54759, 1102, 2293, 1, 66, 1102, 1, 3, 67, 1101, 0, 1623, 68, 1101, 302, 0, 69, 1101, 0, 1, 71, 1101, 0, 1629, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 52053, 1101, 88259, 0, 66, 1101, 0, 1, 67, 1102, 1, 1658, 68, 1101, 556, 0, 69, 1101, 0, 0, 71, 1101, 1660, 0, 72, 1106, 0, 73, 1, 1672, 1101, 0, 12379, 66, 1101, 0, 1, 67, 1102, 1, 1687, 68, 1101, 556, 0, 69, 1102, 1, 1, 71, 1102, 1689, 1, 72, 1105, 1, 73, 1, -3333, 21, 34702, 1102, 39569, 1, 66, 1101, 1, 0, 67, 1102, 1718, 1, 68, 1101, 556, 0, 69, 1102, 4, 1, 71, 1101, 1720, 0, 72, 1106, 0, 73, 1, 1, 7, 319492, 33, 6879, 20, 170398, 27, 307779, 1102, 1, 47797, 66, 1102, 2, 1, 67, 1101, 1755, 0, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 1, 1759, 72, 1105, 1, 73, 0, 0, 0, 0, 38, 185722, 1102, 33767, 1, 66, 1102, 1, 4, 67, 1101, 0, 1788, 68, 1101, 0, 302, 69, 1101, 0, 1, 71, 1101, 0, 1796, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 12, 312955, 1102, 1, 55381, 66, 1102, 1, 1, 67, 1101, 1825, 0, 68, 1102, 556, 1, 69, 1102, 1, 1, 71, 1102, 1827, 1, 72, 1105, 1, 73, 1, 192, 39, 59273, 1101, 80953, 0, 66, 1101, 0, 1, 67, 1101, 0, 1856, 68, 1101, 556, 0, 69, 1101, 0, 0, 71, 1102, 1, 1858, 72, 1106, 0, 73, 1, 1204, 1102, 1, 96451, 66, 1101, 0, 1, 67, 1101, 0, 1885, 68, 1101, 556, 0, 69, 1101, 0, 3, 71, 1102, 1, 1887, 72, 1105, 1, 73, 1, 5, 37, 67534, 37, 101301, 12, 62591, 1101, 0, 54361, 66, 1102, 1, 1, 67, 1101, 0, 1920, 68, 1101, 0, 556, 69, 1102, 1, 0, 71, 1101, 0, 1922, 72, 1106, 0, 73, 1, 1692, 1101, 48073, 0, 66, 1101, 0, 3, 67, 1102, 1949, 1, 68, 1102, 302, 1, 69, 1102, 1, 1, 71, 1102, 1955, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 86755, 1102, 23671, 1, 66, 1102, 1, 1, 67, 1102, 1984, 1, 68, 1102, 556, 1, 69, 1101, 1, 0, 71, 1101, 0, 1986, 72, 1105, 1, 73, 1, 125, 37, 33767, 1102, 1, 62591, 66, 1101, 0, 6, 67, 1101, 0, 2015, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1101, 2027, 0, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 182594, 1102, 57493, 1, 66, 1102, 1, 1, 67, 1101, 2056, 0, 68, 1101, 556, 0, 69, 1101, 0, 1, 71, 1101, 2058, 0, 72, 1106, 0, 73, 1, -23027, 20, 85199, 1102, 1, 79873, 66, 1102, 1, 4, 67, 1102, 1, 2087, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 2095, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 21, 17351, 1102, 63487, 1, 66, 1102, 1, 1, 67, 1102, 1, 2124, 68, 1101, 0, 556, 69, 1101, 0, 7, 71, 1102, 1, 2126, 72, 1105, 1, 73, 1, 2, 39, 177819, 42, 409265, 49, 48073, 49, 144219, 38, 92861, 12, 250364, 12, 375546, 1101, 36691, 0, 66, 1102, 1, 1, 67, 1101, 2167, 0, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1102, 1, 2169, 72, 1106, 0, 73, 1, 256, 3, 95594, 1101, 59273, 0, 66, 1102, 4, 1, 67, 1102, 1, 2198, 68, 1101, 302, 0, 69, 1102, 1, 1, 71, 1101, 2206, 0, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 42, 491118, 1101, 1373, 0, 66, 1101, 0, 3, 67, 1101, 0, 2235, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1101, 2241, 0, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 43, 42197}
	computers := []*Process{}

	for i := 0; i < 50; i++ {
		p := NewProcess(code, []int{i})
		computers = append(computers, p)
	}

	result := 0

	for {
		for index, p := range computers {
			fmt.Println("Processing", index)

			if len(p.input) > p.inputPointer {
				p.AddInput(-1)
			}

			p.RunTilInterupt()

			for len(p.output) > p.outputPointer {
				address := p.NextOutput()

				x := p.NextOutput()
				y := p.NextOutput()

				fmt.Println("Sending to address", address, x, y)

				if address == 255 {
					result = y
					goto finally
				}

				computers[address].AddInput(x)
				computers[address].AddInput(y)
			}
		}
	}

finally:
	fmt.Println(result)
}

func main() {
	part1()
}
