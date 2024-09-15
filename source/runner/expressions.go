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
	case *ast.MemberInstance:
		return r.evaluateMemberInstance(e)
	case *ast.FunctionInstance:
		return r.evaluateFunctionInstance(e)
	case *ast.PrefixExpression:
		return r.evaluatePrefixExpression(e)
	case *ast.AssignmentExpression:
		return r.evaluateAssignmentExpression(e)
	case *ast.ArrayInstance:
		return r.evaluateArrayInstance(e)
	}
	return nil
}

func (r *Runner) evaluateNumberExpression(e *ast.NumberExpression) float64 {
	return e.Value
}

func (r *Runner) evaluateStringExpression(e *ast.StringExpression) string {
	return string(e.Value)
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
		return left.(float64) != right.(float64)
	}
	return nil
}

func (r *Runner) evaluateSymbolExpression(e *ast.SymbolExpression) any {
	variable := r.GetVariable(r.Packages["main"], e.Value)
	if variable == nil {
		return e.Value
	}
	return r.Evaluate(variable.AssignedValue)
}

func (r *Runner) evaluateMemberInstance(e *ast.MemberInstance) any {
	pkg := r.GetPackage(r.Evaluate(e.Instance).(string))
	if pkg == nil {
		class := r.GetClassInstance(r.MainPackage(), r.Evaluate(e.Instance).(string))
		return r.Evaluate(class.Fields[e.MemberName])
	}

	return r.Evaluate(r.GetVariable(pkg, e.MemberName).AssignedValue)
}

func (r *Runner) evaluateFunctionInstance(e *ast.FunctionInstance) any {
	function := r.GetFunction(r.MainPackage(), e.FunctionName)
	if function.Name == "" {
		panic("функция не найдена")
	}
	if len(e.Parameters) != len(function.Parameters) {
		panic("Неверное число аргументов")
	}
	for i := range function.Parameters {
		function.Parameters[i].AssignedValue = e.Parameters[i]
		r.RegisterVariable(r.MainPackage(), &function.Parameters[i])

	}
	result := r.Run(function.Body, "main")
	if len(function.ReturnType) != len(result) {
		panic("Неверное число возвращаемых значений")
	}
	if len(result) == 1 {
		return result[0]
	}
	return result
}

func (r *Runner) evaluatePrefixExpression(e *ast.PrefixExpression) any {
	right := r.Evaluate(e.RightExpr)

	switch e.Op.Kind {
	case lexer.MINUS:
		return -right.(float64)
	case lexer.PLUS:
		return +right.(float64)
	case lexer.NOT:
		return !right.(bool)
	}
	return nil
}

func (r *Runner) evaluateAssignmentExpression(e *ast.AssignmentExpression) any {
	r.RegisterVariable(r.MainPackage(), r.GetVariable(r.MainPackage(), e.Assigne.(*ast.SymbolExpression).Value))
	variable := r.GetVariable(r.MainPackage(), e.Assigne.(*ast.SymbolExpression).Value)
	if variable.IsConstant {
		panic("Нельзя изменять константу")
	}
	variable.AssignedValue = &ast.NumberExpression{Value: r.Evaluate(e.Value).(float64)}
	return variable
}

func (r *Runner) evaluateArrayInstance(e *ast.ArrayInstance) any {
	underlying := r.Evaluate(e.Underlying).(string)
	var indexes []any
	for _, index := range e.Content {
		indexes = append(indexes, r.Evaluate(index))
	}
	return string(underlying[int(indexes[0].(float64))])
}
