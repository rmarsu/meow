package start

import (
	"fmt"
	"io"
	"meow/source/lexer"
	"meow/source/parser"
	"meow/source/runner"
	"os"
	"path/filepath"

	"github.com/sanity-io/litter"
)

func Start(_filepath string) {
	fileExtension := filepath.Ext(_filepath)
	if fileExtension != ".meow" {
		panic(fmt.Sprintf("Файлы языка имеют расширение .meow, не %s", fileExtension))
	}
	file, err := os.Open(_filepath)
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
