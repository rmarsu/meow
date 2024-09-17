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
	case *ast.ArrayDeclaration:
		return r.evaluateArrayDeclaration(e)
	case *ast.ClassInstance:
		return r.evaluateClassInstance(e)
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
	variable := r.GetVariable(r.MainPackage(), e.Value)
	if variable == nil {
		return e.Value
	}
	return r.Evaluate(variable.AssignedValue)
}

func (r *Runner) evaluateMemberInstance(e *ast.MemberInstance) any {
	name := r.Evaluate(e.Instance)
	var pkg *Package
	memberName := (e.MemberName)
	switch n := name.(type) {
	case string:
		pkg = r.GetPackage(n)
		if pkg == nil {
			pkg = r.MainPackage()
		}
	case *ast.ClassInstance:
		switch mB := memberName.(type) {
		case *ast.SymbolExpression:
			return r.Evaluate(n.Fields[mB.Value])
		case *ast.MemberInstance:
			return r.evaluateMemberInstance(mB)
		case *ast.FunctionInstance:
			r.RegisterVariable(r.MainPackage(), &ast.VariableDecStatement{
				Name:          n.ClassName,
				AssignedValue: n,
				Type:          nil,
				IsConstant:    false,
			})
			return r.evaluateFunctionInstance(mB)
		}
		return nil
	}
	switch mB := memberName.(type) {
	case *ast.FunctionInstance:
		function := r.GetFunction(pkg, mB.FunctionName)
		r.RegisterFunction(r.MainPackage(), function)
		return r.Evaluate(mB)
	case *ast.MemberInstance:
		return r.evaluateMemberInstance(mB)
	case *ast.SymbolExpression:
		return r.evaluateSymbolExpression(mB)
	}
	return nil

}

func (r *Runner) evaluateFunctionInstance(e *ast.FunctionInstance) any {
	if e.FunctionName == "len" {
		if len(e.Parameters) > 1 {
			panic("Неверное число аргументов")
		}
		switch param := r.Evaluate(e.Parameters[0]).(type) {
		case string:
			return len([]rune(param))
		case []rune:
			return len(param)
		case []any:
			return float64(len(param))
		}
	}
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
	var variableName string
	switch assign := e.Assigne.(type) {
	case *ast.SymbolExpression:
		variableName = assign.Value
		r.RegisterVariable(r.MainPackage(), r.GetVariable(r.MainPackage(), variableName))
		variable := r.GetVariable(r.MainPackage(), variableName)
		if variable.IsConstant {
			panic("Нельзя изменять константу")
		}
		variable.AssignedValue = &ast.NumberExpression{Value: r.Evaluate(e.Value).(float64)}
	case *ast.MemberInstance:
		variableName = assign.Instance.(*ast.SymbolExpression).Value
		variable := r.GetVariable(r.MainPackage(), variableName)
		variable.AssignedValue.(*ast.ClassInstance).Fields[assign.MemberName.(*ast.SymbolExpression).Value] = e.Value

	}
	return nil
}

func (r *Runner) evaluateArrayInstance(e *ast.ArrayInstance) any {
	underlying := r.Evaluate(e.Underlying)
	var indexes []any
	for _, index := range e.Content {
		indexes = append(indexes, r.Evaluate(index))
	}
	switch underlying := underlying.(type) {
	case string:
		return string(underlying[int(indexes[0].(float64))])
	case []any:
		return underlying[int(indexes[0].(float64))]
	}
	return nil
}

func (r *Runner) evaluateArrayDeclaration(e *ast.ArrayDeclaration) any {
	var array []any
	for i := 0; i < len(e.Elements); i++ {
		array = append(array, r.Evaluate(e.Elements[i]))
	}
	return array
}

func (r *Runner) evaluateClassInstance(e *ast.ClassInstance) any {
	return e
}
