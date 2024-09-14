package runner

import (
	"io"
	"meow/source/ast"
	"meow/source/lexer"
	"meow/source/parser"
	"os"
)



func (r *Runner) RunImportStatement(stmt *ast.ImportStatement) {
	importName := stmt.ImportName
	packagePath := stmt.PackagePath
	file, err := os.Open(packagePath)
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
	r.Run(&ast, importName)
}
