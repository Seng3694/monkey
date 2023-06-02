package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = "🐒 ➤ "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParseErrors(out, parser.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "🙈 Woops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.ByteCode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "🙈 Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "🙈 Parser errors occured:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
