package runner

import (
	"fmt"
	"meow/source/ast"
)

func (r *Runner) RunReturnStatement(stmt ast.Statement) []any {
	returnStmt := stmt.(*ast.ReturnStatement)
	var values []any
	for _, expr := range returnStmt.Expressions {
		values = append(values, r.Evaluate(expr))
	}
	return values
}

func (r *Runner) RunIfStatement(stmt ast.Statement, packagename string) []any {
	var returned []any
	ifStmt := stmt.(*ast.IfStatement)
	if r.Evaluate(ifStmt.Condition).(bool) {
		returned = r.Run(ifStmt.ThenBlock, packagename)
	} else if ifStmt.ElseBlock != nil {
		returned = r.Run(ifStmt.ElseBlock, packagename)
	}
	return returned
}

func (r *Runner) RunPrintStatement(stmt ast.Statement) {
	printStmt := stmt.(*ast.PrintStatement)
	fmt.Println(r.Evaluate(printStmt.Input))
}

func (r *Runner) RunExpressionStatement(stmt ast.Statement) {
	exprStmt := stmt.(*ast.ExpressionStatement)
	r.Evaluate(exprStmt.Expression)
}

func (r *Runner) RunWhileStatement(stmt ast.Statement) {
	whileStmt := stmt.(*ast.WhileStatement)
	for r.Evaluate(whileStmt.Conditions[0]).(bool) {
		r.Run(whileStmt.Body, "main")
	}
}
