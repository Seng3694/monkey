package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func NewFrame(c *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          c,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
