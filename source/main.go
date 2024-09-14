package start

import (
	"io"
	"meow/source/lexer"
	"meow/source/parser"
	"meow/source/runner"
	"os"

	"github.com/sanity-io/litter"
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
	runner := runner.NewRunner()
	runner.Run(&ast, "main")

}

func DebugTree(filepath string) {
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
	litter.Dump(ast)
}
