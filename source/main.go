package start

import (
	"io"
	"meow/source/lexer"
	"meow/source/parser"
	"meow/source/runner"
	"os"
)

func Start(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	input, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	tokens := lexer.Tokenize(string(input))
	ast, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}
	compiler := runner.NewRunner(ast)
	compiler.Run()

}
