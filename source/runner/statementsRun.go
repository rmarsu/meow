package runner

import (
	"meow/source/ast"
	"fmt"
)


func (r *Runner) runVariableDecStatement(stmt ast.Statement) {
	varDecStmt, _ := stmt.(*ast.VariableDecStatement)
	if varDecStmt.AssignedValue != nil {
		r.Memory.SetVariable(varDecStmt.Name, r.evaluate(varDecStmt.AssignedValue))
	}

}

func (r *Runner) runExpressionStatement(stmt ast.Statement) {
	expStmt, _ := stmt.(*ast.ExpressionStatement)
	r.evaluate(expStmt.Expression)
}

func (r *Runner) runClassDecStatement(stmt ast.Statement) {
	classDecStmt, _ := stmt.(*ast.ClassDecStatement)
	r.Memory.SetClass(classDecStmt.Name, *classDecStmt)
}

func (r *Runner) runFunctionDecStatement(stmt ast.Statement) {
	funcDecStmt, _ := stmt.(*ast.FunctionDecStatement)
	r.Memory.SetFunction(funcDecStmt.Name, *funcDecStmt)
}

func (r *Runner) runReturnStatement(stmt ast.Statement) []any {
	returnStmt, _ := stmt.(*ast.ReturnStatement)
	var result []any
    for _, expr := range returnStmt.Expressions {
		result = append(result, r.evaluate(expr))
	}
	return result
}  

func (r *Runner) runIfStatement(stmt ast.Statement) {
	ifStmt, _ := stmt.(*ast.IfStatement)
    if r.evaluate(ifStmt.Condition).(bool) {
        r.runBlockStatement(ifStmt.ThenBlock)
    } else {
		if ifStmt.ElseBlock!= nil {
            r.runBlockStatement(ifStmt.ElseBlock)
        }
	}
}

func (r *Runner) runWhileStatement(stmt ast.Statement) {
    whileStmt, _ := stmt.(*ast.WhileStatement)
    for r.evaluate(whileStmt.Conditions[0]).(bool) {
        r.runBlockStatement(whileStmt.Body)
    }
}

func (r *Runner) runBlockStatement(block *ast.BlockStatement) []any {
	for _, stmt := range block.Statements {
        switch stmt.(type) {
        case *ast.ReturnStatement:
                return r.runReturnStatement(stmt)
        }
        r.runStatement(stmt)
    }
    return nil
}

func (r *Runner) runPrintStatement(stmt ast.Statement) {
    printStmt, _ := stmt.(*ast.PrintStatement)
    fmt.Println(r.evaluate(printStmt.Input))
}

func (r *Runner) runStatement(stmt ast.Statement) {
    if len(r.Errors) > 0 {
        fmt.Println(r.Errors)
        return
    }
	switch stmt.(type) {
        case *ast.PrintStatement:
            r.runPrintStatement(stmt)
        case *ast.BlockStatement:
            r.runBlockStatement(stmt.(*ast.BlockStatement))
        case *ast.VariableDecStatement:
            r.runVariableDecStatement(stmt)
        case *ast.ExpressionStatement:
            r.runExpressionStatement(stmt)
        case *ast.ClassDecStatement:
            r.runClassDecStatement(stmt)
        case *ast.FunctionDecStatement:
            r.runFunctionDecStatement(stmt)
        case *ast.ReturnStatement:
            r.runReturnStatement(stmt)
        case *ast.IfStatement:
            r.runIfStatement(stmt)
        case *ast.WhileStatement:
            r.runWhileStatement(stmt)
        default:
            fmt.Println(' ')
    }
}