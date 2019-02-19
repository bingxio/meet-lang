package repl

import (
	"bufio"
	"fmt"
	"io"
	"meet-lang/src/environment"
	"meet-lang/src/interpreter"
	"meet-lang/src/lexer"
	"meet-lang/src/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, env *environment.Environment) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println("Good bye, Have a wonderful day !")
			break
		} else if line == "" {
			continue
		}

		tokens := lexer.New(line)

		parser := parser.New(tokens)
		ast := parser.ParseProgram()

		interpreter.Eval(ast, env)
	}
}
