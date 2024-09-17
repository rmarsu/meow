package runner

import (
	"meow/source/ast"
)

type Package struct {
	IsMain bool
	Memory
}

type Runner struct {
	Packages map[string]*Package
	TemporaryMemory Memory
}

type Memory struct {
	Variables        map[string]*ast.VariableDecStatement
	Functions        map[string]*ast.FunctionDecStatement
	Classes          map[string]*ast.ClassDecStatement
	ClassesInstances map[string]*ast.ClassInstance
}

func NewRunner() *Runner {
	return &Runner{
		Packages: make(map[string]*Package),
		TemporaryMemory: Memory{},
	}
}

func (r *Runner) Run(tree *ast.BlockStatement, packagename string) []any {
	for _, stmt := range tree.Statements {
		returning := r.Execute(stmt, packagename)
		if returning != nil {
			return returning
		}
	}
	return nil
}

func (r *Runner) Execute(s ast.Statement, packagename string) []any {
	pkg := r.initPackage(packagename)
	switch stmt := s.(type) {
	case *ast.ImportStatement:
		r.RunImportStatement(stmt)
	case *ast.ClassDecStatement:
		r.RegisterClass(pkg, stmt)
	case *ast.FunctionDecStatement:
		r.RegisterFunction(pkg, stmt)
	case *ast.VariableDecStatement:
		r.RegisterVariable(pkg, stmt)
	case *ast.BlockStatement:
		r.Run(stmt, packagename)
	case *ast.ReturnStatement:
		return r.RunReturnStatement(stmt)
	case *ast.IfStatement:
		r.RunIfStatement(stmt, packagename)
	case *ast.PrintStatement:
		r.RunPrintStatement(stmt)
	case *ast.ExpressionStatement:
		r.RunExpressionStatement(stmt)
	case *ast.WhileStatement:
		r.RunWhileStatement(stmt)
	}
	return nil
}
