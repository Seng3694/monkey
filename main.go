package main

import (
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/repl"
	"monkey/vm"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 {
		runScript(os.Args[1])
	} else {
		fmt.Printf(
			"Hello %s! This is the Monkey programming language!\n",
			user.Username)
		fmt.Print("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func runScript(file string) {
	contents, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	lexer := lexer.New(string(contents))
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		printErrors(os.Stderr, "Parser", parser.Errors())
		return
	}

	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		printError(os.Stderr, "Compiler", err)
		return
	}

	machine := vm.New(comp.ByteCode())

	err = machine.Run()
	if err != nil {
		printError(os.Stderr, "VM", err)
		return
	}

	io.WriteString(os.Stdout, machine.LastPoppedStackElement().Inspect())
	io.WriteString(os.Stdout, "\n")
}

func printErrors(out io.Writer, module string, errors []string) {
	io.WriteString(out, fmt.Sprintf("🙈 %s errors occured:\n", module))
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printError(out io.Writer, module string, err error) {
	io.WriteString(out, fmt.Sprintf("🙈 %s error occured: %s\n", module, err))
}
