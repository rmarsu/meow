package runner

import (
	"meow/source/ast"
	"meow/source/lexer"
)

func (r *Runner) Evaluate(expr ast.Expression) any {
	switch e := expr.(type) {
	case *ast.NumberExpression:
		return r.evaluateNumberExpression(e)
	case *ast.StringExpression:
		return r.evaluateStringExpression(e)
	case *ast.BOExpression:
		return r.evaluateBinaryExpression(e)
	case *ast.SymbolExpression:
		return r.evaluateSymbolExpression(e)
	case *ast.ImportedVariableExpression:
		return r.evaluateImportedVariableExpression(e)
	}
	return nil
}

func (r *Runner) evaluateNumberExpression(e *ast.NumberExpression) float64 {
	return e.Value
}

func (r *Runner) evaluateStringExpression(e *ast.StringExpression) string {
	return e.Value
}

func (r *Runner) evaluateBinaryExpression(e *ast.BOExpression) any {
	left := r.Evaluate(e.Left)
	right := r.Evaluate(e.Right)

	switch e.Op.Kind {
	case lexer.PLUS:
		return left.(float64) + right.(float64)
	case lexer.MINUS:
		return left.(float64) - right.(float64)
	case lexer.MUL:
		return left.(float64) * right.(float64)
	case lexer.DIV:
		return left.(float64) / right.(float64)
	case lexer.LESS:
		return left.(float64) < right.(float64)
	case lexer.LESS_EQUALS:
		return left.(float64) <= right.(float64)
	case lexer.GREATER:
		return left.(float64) > right.(float64)
	case lexer.GREATER_EQUALS:
		return left.(float64) >= right.(float64)
	case lexer.EQUALS:
		return left.(float64) == right.(float64)
	case lexer.NOT_EQUALS:
		return left.(float64)!= right.(float64)
	}
	return nil
}

func (r *Runner) evaluateSymbolExpression(e *ast.SymbolExpression) any {
	return r.Evaluate(r.GetVariable(r.Packages["main"], e.Value).AssignedValue)
}

func (r *Runner) evaluateImportedVariableExpression(e *ast.ImportedVariableExpression) any {
	pkg := r.GetPackage(e.ImportName)
    return r.Evaluate(r.GetVariable(pkg, e.Variable).AssignedValue)
}