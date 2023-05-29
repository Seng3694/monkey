package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type OpCode byte

const (
	OpConstant OpCode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
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
