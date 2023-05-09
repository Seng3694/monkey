package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
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

		tok := lexer.NextToken()
		for tok.Type != token.EOF {
			fmt.Fprintf(out, "%+v\n", tok)
			tok = lexer.NextToken()
		}
	}
}
