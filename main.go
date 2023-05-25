package main

import (
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
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
		printParseErrors(os.Stderr, parser.Errors())
		return
	}

	evaluated := evaluator.Eval(program, object.NewEnvironment())
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "🙈 Parser errors occured:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
