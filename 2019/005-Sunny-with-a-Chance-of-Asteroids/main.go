package main

import (
	"fmt"
)

const (
	OpcodeAdd      = 1
	OpcodeMultiply = 2
	OpcodeHalt     = 99
)

type Program struct {
	memory   []int
	position int
}

func NewProgram(code []int) *Program {
	memory := make([]int, len(code))

	copy(memory, code)

	return &Program{
		memory:   memory,
		position: 0,
	}
}

//
// Read & write to program memory
//
func (p *Program) Read(position int) (int, error) {
	if position >= len(p.memory) || position < 0 {
		return 0, fmt.Errorf("Index %d out of range", position)
	}

	return p.memory[position], nil
}

func (p *Program) Write(position int, value int) error {
	if position >= len(p.memory) || position < 0 {
		return fmt.Errorf("Index out of range")
	}

	p.memory[position] = value

	return nil
}

func (p *Program) DumpMemory() {
	for i, v := range p.memory {
		if i == p.position {
			fmt.Printf("[%d] ", v)
		} else {
			fmt.Printf("%d ", v)
		}
	}

	fmt.Print("\n")
}

//
// Run program until halt or error.
//
func (p *Program) Run() error {
	operation, err := p.Read(p.position)
	if err != nil {
		return err
	}

	switch operation {
	case OpcodeAdd:
		inputPointer1, err := p.Read(p.position + 1)
		if err != nil {
			return err
		}

		inputPointer2, err := p.Read(p.position + 2)
		if err != nil {
			return err
		}

		input1, err := p.Read(inputPointer1)
		if err != nil {
			return err
		}

		input2, err := p.Read(inputPointer2)
		if err != nil {
			return err
		}

		resultPointer, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		p.Write(resultPointer, input1+input2)
	case OpcodeMultiply:
		inputPointer1, err := p.Read(p.position + 1)
		if err != nil {
			return err
		}

		inputPointer2, err := p.Read(p.position + 2)
		if err != nil {
			return err
		}

		input1, err := p.Read(inputPointer1)
		if err != nil {
			return err
		}

		input2, err := p.Read(inputPointer2)
		if err != nil {
			return err
		}

		resultPointer, err := p.Read(p.position + 3)
		if err != nil {
			return err
		}

		p.Write(resultPointer, input1*input2)
	case OpcodeHalt:
		return nil
	default:
		return fmt.Errorf("Unknonwn opcode %d", operation)
	}

	p.position += 4

	return p.Run()
}

func main() {
	programCode := []int{
		1, 0, 0, 3,
		1, 1, 2, 3,
		1, 3, 4, 3,
		1, 5, 0, 3,
		2, 9, 1, 19,
		1, 19, 5, 23,
		1, 9, 23, 27,
		2, 27, 6, 31,
		1, 5, 31, 35,
		2, 9, 35, 39,
		2, 6, 39, 43,
		2, 43, 13, 47,
		2, 13, 47, 51,
		1, 10, 51, 55,
		1, 9, 55, 59,
		1, 6, 59, 63,
		2, 63, 9, 67,
		1, 67, 6, 71,
		1, 71, 13, 75,
		1, 6, 75, 79,
		1, 9, 79, 83,
		2, 9, 83, 87,
		1, 87, 6, 91,
		1, 91, 13, 95,
		2, 6, 95, 99,
		1, 10, 99, 103,
		2, 103, 9, 107,
		1, 6, 107, 111,
		1, 10, 111, 115,
		2, 6, 115, 119,
		1, 5, 119, 123,
		1, 123, 13, 127,
		1, 127, 5, 131,
		1, 6, 131, 135,
		2, 135, 13, 139,
		1, 139, 2, 143,
		1, 143, 10, 0,
		99,
		2, 0, 14, 0,
	}

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			fmt.Printf("Trying %d and %d => ", verb, noun)

			p := NewProgram(programCode)

			p.Write(1, noun)
			p.Write(2, verb)

			err := p.Run()
			if err != nil {
				fmt.Println(err)
				continue
			}

			value, err := p.Read(0)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(value)

			if value == 19690720 {
				fmt.Println(100*noun + verb)
				return
			}
		}
	}
}
