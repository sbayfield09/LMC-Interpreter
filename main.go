package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MemorySize = 100
	MaxValue   = 999
	MinValue   = -999
)

type Instruction struct {
	OpCode  string
	Address int
}

type Assembler struct {
	labels   map[string]int
	memory   [MemorySize]int
	codeSize int
}

type CPU struct {
	accumulator    int
	programCounter int
	memory         [MemorySize]int
	halted         bool
	scanner        *bufio.Scanner
}

func NewAssembler() *Assembler {
	return &Assembler{
		labels: make(map[string]int),
	}
}

func (a *Assembler) Assemble(lines []string) error {
	instructions := []Instruction{}
	a.labels = make(map[string]int)

	passOne := []struct {
		line    string
		label   string
		mnemonic string
		operand string
	}{}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, ";") {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}

		idx := 0
		label := ""
		if _, err := strconv.Atoi(fields[0]); err != nil && !isMnemonic(fields[0]) {
			if len(fields) > 1 {
				label = fields[0]
				idx = 1
			} else {
				return fmt.Errorf("line %d: invalid syntax: %s", i+1, line)
			}
		}

		if idx >= len(fields) {
			return fmt.Errorf("line %d: missing mnemonic: %s", i+1, line)
		}

		mnemonic := strings.ToUpper(fields[idx])
		idx++

		operand := ""
		if idx < len(fields) {
			operand = fields[idx]
		}

		if label != "" {
			if _, exists := a.labels[label]; exists {
				return fmt.Errorf("line %d: duplicate label: %s", i+1, label)
			}
			a.labels[label] = a.codeSize
		}

		passOne = append(passOne, struct {
			line    string
			label   string
			mnemonic string
			operand string
		}{line: line, label: label, mnemonic: mnemonic, operand: operand})

		if mnemonic != "DAT" {
			a.codeSize++
		} else {
			a.codeSize++
		}

		if a.codeSize > MemorySize {
			return fmt.Errorf("line %d: program exceeds memory size", i+1)
		}
	}

	for i, entry := range passOne {
		instr, err := a.assembleInstruction(entry.mnemonic, entry.operand, i+1)
		if err != nil {
			return err
		}
		instructions = append(instructions, instr)
	}

	a.memory = [MemorySize]int{}
	for i, instr := range instructions {
		if instr.OpCode == "DAT" {
			a.memory[i] = instr.Address
		} else {
			opCode, err := strconv.Atoi(instr.OpCode)
			if err != nil {
				return fmt.Errorf("invalid opcode: %s", instr.OpCode)
			}
			a.memory[i] = opCode*100 + instr.Address
		}
	}

	return nil
}

func (a *Assembler) assembleInstruction(mnemonic, operand string, lineNum int) (Instruction, error) {
	switch mnemonic {
	case "ADD", "1":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "1", Address: addr}, nil
	case "SUB", "2":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "2", Address: addr}, nil
	case "STA", "STO", "3":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "3", Address: addr}, nil
	case "LDA", "5":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "5", Address: addr}, nil
	case "BRA", "6":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "6", Address: addr}, nil
	case "BRZ", "7":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "7", Address: addr}, nil
	case "BRP", "8":
		addr, err := a.resolveAddress(operand, lineNum)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{OpCode: "8", Address: addr}, nil
	case "INP", "IN":
		return Instruction{OpCode: "9", Address: 1}, nil
	case "OUT":
		return Instruction{OpCode: "9", Address: 2}, nil
	case "CHAR":
		return Instruction{OpCode: "9", Address: 3}, nil
	case "HLT", "COB":
		return Instruction{OpCode: "0", Address: 0}, nil
	case "DAT":
		value := 0
		if operand != "" {
			var err error
			value, err = strconv.Atoi(operand)
			if err != nil {
				return Instruction{}, fmt.Errorf("line %d: invalid DAT value: %s", lineNum, operand)
			}
			if value < MinValue || value > MaxValue {
				return Instruction{}, fmt.Errorf("line %d: DAT value out of range: %d", lineNum, value)
			}
		}
		return Instruction{OpCode: "DAT", Address: value}, nil
	default:
		return Instruction{}, fmt.Errorf("line %d: unknown mnemonic: %s", lineNum, mnemonic)
	}
}

func (a *Assembler) resolveAddress(operand string, lineNum int) (int, error) {
	if operand == "" {
		return 0, fmt.Errorf("line %d: missing operand", lineNum)
	}

	addr, err := strconv.Atoi(operand)
	if err == nil {
		if addr < 0 || addr >= MemorySize {
			return 0, fmt.Errorf("line %d: address out of range: %d", lineNum, addr)
		}
		return addr, nil
	}

	addr, exists := a.labels[operand]
	if !exists {
		return 0, fmt.Errorf("line %d: undefined label: %s", lineNum, operand)
	}
	return addr, nil
}

func isMnemonic(s string) bool {
	mnemonics := []string{
		"ADD", "SUB", "STA", "STO", "LDA", "BRA", "BRZ", "BRP",
		"INP", "IN", "OUT", "CHAR", "HLT", "COB", "DAT",
	}
	upper := strings.ToUpper(s)
	for _, m := range mnemonics {
		if upper == m {
			return true
		}
	}
	return false
}

func NewCPU(memory [MemorySize]int, scanner *bufio.Scanner) *CPU {
	return &CPU{
		accumulator:    0,
		programCounter: 0,
		memory:         memory,
		halted:         false,
		scanner:        scanner,
	}
}

func (c *CPU) Run() error {
	maxSteps := 10000
	step := 0

	for !c.halted && step < maxSteps {
		if c.programCounter < 0 || c.programCounter >= MemorySize {
			return fmt.Errorf("program counter out of range: %d", c.programCounter)
		}

		word := c.memory[c.programCounter]
		c.programCounter++

		opCode := word / 100
		address := word % 100

		switch opCode {
		case 0:
			if word == 0 {
				c.halted = true
			}
		case 1:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for ADD: %d", address)
			}
			c.accumulator += c.memory[address]
		case 2:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for SUB: %d", address)
			}
			c.accumulator -= c.memory[address]
		case 3:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for STA: %d", address)
			}
			c.memory[address] = c.accumulator
		case 5:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for LDA: %d", address)
			}
			c.accumulator = c.memory[address]
		case 6:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for BRA: %d", address)
			}
			c.programCounter = address
		case 7:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for BRZ: %d", address)
			}
			if c.accumulator == 0 {
				c.programCounter = address
			}
		case 8:
			if address < 0 || address >= MemorySize {
				return fmt.Errorf("invalid address for BRP: %d", address)
			}
			if c.accumulator >= 0 {
				c.programCounter = address
			}
		case 9:
			if word == 901 {
				var input int
				for {
					fmt.Print("Input: ")
					if !c.scanner.Scan() {
						return fmt.Errorf("error reading input: %v", c.scanner.Err())
					}
					line := c.scanner.Text()
					line = strings.TrimSpace(line)
					val, err := strconv.Atoi(line)
					if err != nil || val < 0 || val > MaxValue {
						fmt.Printf("Please enter a number between 0 and %d\n", MaxValue)
						continue
					}
					input = val
					break
				}
				c.accumulator = input
			} else if word == 902 {
				fmt.Printf("%d\n", c.accumulator)
			} else if word == 903 {
				if c.accumulator < 0 || c.accumulator > 127 {
					return fmt.Errorf("CHAR: accumulator value %d is outside ASCII range (0-127)", c.accumulator)
				}
				fmt.Printf("%c", c.accumulator)
			} else {
				return fmt.Errorf("unknown instruction: %d", word)
			}
		default:
			return fmt.Errorf("unknown opcode: %d", opCode)
		}

		step++
	}

	if step >= maxSteps {
		return fmt.Errorf("execution exceeded maximum steps (possible infinite loop)")
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.lmc>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	assembler := NewAssembler()
	if err := assembler.Assemble(lines); err != nil {
		fmt.Fprintf(os.Stderr, "Assembly error: %v\n", err)
		os.Exit(1)
	}

	cpu := NewCPU(assembler.memory, bufio.NewScanner(os.Stdin))
	if err := cpu.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}
