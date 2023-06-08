package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type EmittedInstruction struct {
	OpCode   code.OpCode
	Position int
}

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "<":
			c.emit(code.OpLesserThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">=":
			c.emit(code.OpGreaterEqual)
		case "<=":
			c.emit(code.OpLesserEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// address will be patched later
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 0x1deadb0b)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstruction.OpCode == code.OpPop {
			c.removeLastInstruction()
		}

		jumpPos := c.emit(code.OpJump, 0x1deadb0b)

		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstruction.OpCode == code.OpPop {
				c.removeLastInstruction()
			}
		}

		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)
	}
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.OpCode, operands ...int) int {
	instr := code.Make(op, operands...)
	pos := c.addInstruction(instr)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) setLastInstruction(op code.OpCode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{OpCode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.OpCode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	for i := 1; i < len(newInstruction); i++ {
		c.instructions[opPos+i] = newInstruction[i]
	}
}

func (c *Compiler) addInstruction(instr []byte) int {
	posNewInstr := len(c.instructions)
	c.instructions = append(c.instructions, instr...)
	return posNewInstr
}
