package runner

import (
	"fmt"
	"meow/source/ast"
)

func (r *Runner) RunReturnStatement(stmt ast.Statement) {
	return
}

func (r *Runner) RunIfStatement(stmt ast.Statement, packagename string) {
	ifStmt := stmt.(*ast.IfStatement)
	if r.Evaluate(ifStmt.Condition).(bool) {
		r.Run(ifStmt.ThenBlock, packagename)
	} else if ifStmt.ElseBlock != nil {
		r.Run(ifStmt.ElseBlock, packagename)
	}
}

func (r *Runner) RunPrintStatement(stmt ast.Statement) {
	printStmt := stmt.(*ast.PrintStatement)
	fmt.Println(r.Evaluate(printStmt.Input))
}

func (r *Runner) RunExpressionStatement(stmt ast.Statement) {
	return
}

func (r *Runner) RunWhileStatement(stmt ast.Statement) {
	return
}
