package runner

import (
	"fmt"
	"meow/source/ast"
)

func (r *Runner) RunReturnStatement(stmt ast.Statement) {
	return;
} 

func (r *Runner) RunIfStatement(stmt ast.Statement) {
    return;
}

func (r *Runner) RunPrintStatement(stmt ast.Statement) {
	printStmt := stmt.(*ast.PrintStatement)
	fmt.Println(r.Evaluate(printStmt.Input))
}

func (r *Runner) RunExpressionStatement(stmt ast.Statement) {
    return;
}

func (r *Runner) RunWhileStatement(stmt ast.Statement) {
    return;
}

