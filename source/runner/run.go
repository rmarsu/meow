package runner

import (
	"meow/source/ast"
)

type Runner struct {
	Errors []error
	AST     ast.BlockStatement
	Memory
	Output   string
}

func NewRunner(ast ast.BlockStatement) *Runner {
     return &Runner{AST: ast, Memory:*NewMemory()}
}

func (r *Runner) Run() {
	Tree := r.AST
	r.runBlockStatement(&Tree)
}