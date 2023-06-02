package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type OpCode byte

const (
	OpConstant OpCode = iota
	OpAdd
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op OpCode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instrLen := 1
	for _, w := range def.OperandWidths {
		instrLen += w
	}

	instr := make([]byte, instrLen)
	instr[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instr[offset:], uint16(o))
		}
		offset += width
	}
	return instr
}

func (instr Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(instr) {
		def, err := Lookup(instr[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, instr[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, instr.fmtInstruction(def, operands))
		i += 1 + read
	}

	return out.String()
}

func (instr Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func ReadOperands(def *Definition, instr Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instr[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(instr Instructions) uint16 {
	return binary.BigEndian.Uint16(instr)
}