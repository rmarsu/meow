package runner

import (
	"fmt"
	"io"
	"meow/source/ast"
	"meow/source/lexer"
	"meow/source/parser"
	"meow/source/runner/object"
	"os"
)

func Evaluate(node ast.Expression, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.ClassInstance:
		return evaluateClassInstance(node, env)
	case *ast.MemberInstance:
		instance := node.Instance.(*ast.SymbolExpression).Value
		member := node.MemberName
		value := evaluateMemberInstance(instance, member, env)
		return value
	case *ast.AssignmentExpression:
		value := Evaluate(node.Value, env)
		assigne := node.Assigne.(*ast.SymbolExpression).Value
		if IsError(value) {
			return value
		}
		env.Set(assigne, value)
		return value
	case *ast.NumberExpression:
		if isWhole(node.Value) {
			return &object.Integer{Value: int64(node.Value)}
		} else {
			return &object.Float{Value: node.Value}
		}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringExpression:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		right := Evaluate(node.RightExpr, env)
		return evaluatePrefixExpression(node, right)
	case *ast.BOExpression:
		left := Evaluate(node.Left, env)
		if IsError(left) {
			return left
		}
		right := Evaluate(node.Right, env)
		if IsError(right) {
			return right
		}
		return evaluateBOExpression(node.Op.Kind, left, right)
	case *ast.SymbolExpression:
		return evaluateSymbolExpression(node, env)
	case *ast.FunctionInstance:
		defaults := checkForDefault(node.FunctionName)
		if defaults {
			switch node.FunctionName {
			case "meow":
				arg := Evaluate(node.Parameters[0], env)
				fmt.Println(arg.Inspect())
				return nil
			case "len":
				value := Evaluate(node.Parameters[0], env)
				switch val := value.(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(val.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(val.Elements))}
				default:
					return newError("Невозможно высчитать длину типа %s", val.Type())

				}
			case "tail":
				args := EvaluateExpressions(node.Parameters, env)
				if len(args) != 2 {
					return newError("Функция tail требует два аргумента")
				}
				if args[0].Type() != object.ARRAY {
					return newError("Первый аргумент функции tail должен быть массивом")
				}
				arr := args[0].(*object.Array)
				if args[1].Type() != arr.ElementsType {
					return newError("Второй аргумент функции tail должен быть %s", arr.ElementsType)
				}
				length := len(arr.Elements)
				newElements := make([]object.Object, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements, ElementsType: arr.ElementsType}
			}
		}
		functionObject, ok := env.Get(node.FunctionName)
		if !ok {
			return newError("Неизвестная функция: %s", node.FunctionName)
		}
		args := EvaluateExpressions(node.Parameters, env)
		if len(args) == 1 {
			return args[0]
		}
		return applyFunction(functionObject, args)
	case *ast.ArrayDeclaration:
		elements := EvaluateExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		_type := elements[0].Type()
		for _, elem := range elements {
			if elem.Type() != _type {
				return newError("Все элементы массива должны быть одного типа")
			}
		}
		return &object.Array{Elements: elements, ElementsType: _type}
	case *ast.ArrayInstance:
		left := Evaluate(node.Underlying, env)
		if IsError(left) {
			return left
		}
		index := Evaluate(node.Content[0], env)
		if IsError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}
	return NULL
}

func evaluateMemberInstance(instance string, member ast.Expression, env *object.Environment) object.Object {
	instanceVal, ok := env.Get(instance)
	if !ok {
		return newError("Объект %s не найден", instance)
	}
	var memberName string
	switch member := member.(type) {
	case *ast.SymbolExpression:
		memberName = member.Value
	case *ast.FunctionInstance:
		memberName = member.FunctionName
	}
	switch instanceVal.Type() {
	case object.CLASS:
		class := instanceVal.(*object.Class)
		field, ok := class.Fields[memberName]
		if !ok {
			function, ok := class.Functions[memberName]
			if !ok {
				return newError("Функция %s не найдена в классе %s", memberName, class.Name)
			}
			params := EvaluateExpressions(member.(*ast.FunctionInstance).Parameters, env)
			actualClass, _ := env.Get(class.Name)
			env.Set(class.Name, instanceVal)
			result := applyFunction(function, params)
			env.Set(class.Name, actualClass)

			return result
		}
		return field
	case object.MODULE:
		module := instanceVal.(*object.Module)
		field, ok := module.Environment.Get(memberName)
		if !ok {
			return newError(" %s не найдено в модуле %s", memberName, module.Name)
		}
		if field.Type() == object.FUNCTION {
			params := EvaluateExpressions(member.(*ast.FunctionInstance).Parameters, env)
			result := applyFunction(field, params)
			return result
		}
		return field
	}
	return newError("Невозможно получить доступ к полю %s у объекта %s", member, instanceVal.Type())
}

func evaluateClassInstance(node *ast.ClassInstance, env *object.Environment) object.Object {
	parentClass, ok := env.Get(node.ClassName)
	if !ok {
		return newError("Класс %s не найден", node.ClassName)
	}
	class := parentClass.(*object.Class)
	var fields = make(map[string]object.Object)
	for index, expr := range node.Fields {
		if expr == nil {
			continue
		}
		value := Evaluate(expr, env)
		if IsError(value) {
			return value
		}
		if class.Fields[index].Type() != value.Type() {
			return newError("Невозможно присвоить полю %s объекта %s неверного типа", index, class.Name)
		}
		fields[index] = value
	}
	return &object.Class{
		Name:      node.ClassName,
		Fields:    fields,
		Functions: class.Functions,
	}

}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	}
	return newError("Невозможно получить доступ по индексу для объекта %s", left.Type())
}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	indexValue := index.(*object.Integer).Value
	arrayValue := array.(*object.Array)
	if indexValue < 0 || indexValue >= int64(len(arrayValue.Elements)) {
		return newError("Индекс выходит за границы массива")
	}
	return arrayValue.Elements[int(indexValue)]
}

func evalStringIndexExpression(stringObject object.Object, index object.Object) object.Object {
	indexValue := index.(*object.Integer).Value
	stringValue := stringObject.(*object.String)
	if indexValue < 0 || int64(indexValue) >= int64(len(stringValue.Value)) {
		return newError("Индекс выходит за границы строки")
	}
	return &object.String{Value: []rune(string(stringValue.Value[int(indexValue)]))}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.FunctionLiteral)
	if !ok {
		return newError("Не является функцией")
	}
	extendedEnv := extendFunctionEnv(function, args)
	executed := Execute(function.Body, extendedEnv)
	if IsError(executed) {
		return executed
	}
	return unwrapReturn(executed, function.ReturnType)
}

func extendFunctionEnv(fn *object.FunctionLiteral, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for index, param := range fn.Parameters {
		env.Set(param.Names[0], args[index])
	}
	return env
}

func unwrapReturn(obj object.Object, _types []object.ObjectType) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		for i := 0; i < len(returnValue.Values); i++ {
			if returnValue.Values[i].Type() != _types[i] {
				return newError("Невозможно привести возвращаемое значение к %s", _types[i])
			}
		}
		return returnValue
	}
	return obj
}

func EvaluateExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	evaluated := []object.Object{}
	for _, expr := range exprs {
		evaluated = append(evaluated, Evaluate(expr, env))
	}
	return evaluated
}

func evaluatePrefixExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch node.Op.Kind {
	case lexer.MINUS:
		return evalMinusOperatorExpr(right)
	default:
		return newError("Неизвестный оператор: %s, %s", node.Op.Value, right.Type())
	}
}

func evaluateSymbolExpression(node *ast.SymbolExpression, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)
	if !ok {
		return newError("Неизвестная переменная: %s", node.Value)
	}
	return value
}

func evaluateBOExpression(operator lexer.TokenKind, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerBOExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringBOExpression(operator, left, right)
	}
	return newError("Невозможно бинарное действие типов %s, %s", left.Type(), right.Type())

}

func evalIntegerBOExpression(operator lexer.TokenKind, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("Нельзя выполнить операцию с разными типами: %s, %s", left.Type(), right.Type())
	}
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch {
	case operator == lexer.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case operator == lexer.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case operator == lexer.MUL:
		return &object.Integer{Value: leftVal * rightVal}
	case operator == lexer.DIV:
		if rightVal == 0 {
			return newError("Деление на ноль")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case operator == lexer.LESS:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case operator == lexer.GREATER:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case operator == lexer.GREATER_EQUALS:
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case operator == lexer.LESS_EQUALS:
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case operator == lexer.NOT_EQUALS:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case operator == lexer.EQUALS:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	return newError("Неизвестный оператор")
}

func evalStringBOExpression(operator lexer.TokenKind, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return newError("Нельзя выполнить операцию с разными типами: %s, %s", left.Type(), right.Type())
	}
	leftVal := string(left.(*object.String).Value)
	rightVal := string(right.(*object.String).Value)
	switch operator {
	case lexer.PLUS:
		return &object.String{Value: []rune(leftVal + rightVal)}
	case lexer.EQUALS:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	return newError("Неизвестный оператор")
}

func evalMinusOperatorExpr(right object.Object) object.Object {
	if right.Type() != (object.INTEGER) {
		panic(fmt.Sprintf("unknown operand type for '-' operator: %s", right.Type()))
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func Execute(node ast.Statement, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.ImportStatement:
		return ExecuteImportStat(*node, env)
	case *ast.ClassDecStatement:
		return ExecuteClassDec(*node, env)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)
	case *ast.BlockStatement:
		return ExecuteBlock(*node, env)
	case *ast.IfStatement:
		return EvaluateIf(*node, env)
	case *ast.ReturnStatement:
		val := EvaluateExpressions(node.Expressions, env)
		return &object.ReturnValue{Values: val}
	case *ast.VariableDecStatement:
		val := Evaluate(node.AssignedValue, env)
		if val.Type() == object.RETURN_VALUE {
			for i := 0; i < len(val.(*object.ReturnValue).Values); i++ {
				env.Set(node.Names[i], val.(*object.ReturnValue).Values[i])
			}
		} else {
			env.Set(node.Names[0], val)
		}
	case *ast.FunctionDecStatement:
		params := node.Parameters
		body := node.Body
		var returnTypes []object.ObjectType
		for _, _type := range node.ReturnType {
			if _type.(*ast.SymbolType).Name == "string" {
				returnTypes = append(returnTypes, object.STRING)
			} else if _type.(*ast.SymbolType).Name == "bool" {
				returnTypes = append(returnTypes, object.BOOLEAN)
			} else if _type.(*ast.SymbolType).Name == "int" {
				returnTypes = append(returnTypes, object.INTEGER)
			} else if _type.(*ast.SymbolType).Name == "array" {
				returnTypes = append(returnTypes, object.ARRAY)
			}
		}
		function := &object.FunctionLiteral{
			Env:        env,
			Parameters: params,
			Body:       body,
			ReturnType: returnTypes,
		}
		functionFromEnv := env.Set(node.Name, function)
		return functionFromEnv

	case *ast.WhileStatement:
		conditions := EvaluateExpressions(node.Conditions, env)
		for isAllTruthy(conditions) {
			ExecuteBlock(*node.Body, env)
			conditions = EvaluateExpressions(node.Conditions, env)
		}
	}
	return NULL
}

func ExecuteClassDec(node ast.ClassDecStatement, env *object.Environment) object.Object {
	className := node.Name
	var variables = make(map[string]object.Object)
	var functions = make(map[string]object.Object)
	for index, variable := range node.Fields {
		tmp := EvaluateClassField(variable, env)
		variables[index] = tmp
	}
	for index, fn := range node.Functions {
		tmp := EvaluateFunctionField(fn, env, index)
		functions[index] = tmp

	}
	class := &object.Class{
		Name:      className,
		Fields:    variables,
		Functions: functions,
	}
	env.Set(className, class)
	return class
}

func EvaluateClassField(variable ast.ClassFieldStatement, env *object.Environment) object.Object {
	switch variable.Type.(type) {
	case *ast.ArrayType:
		return &object.Array{}
	}
	switch variable.Type.(*ast.SymbolType).Name {
	case "string":
		return &object.String{}
	case "int":
		return &object.Integer{}
	case "bool":
		return &object.Boolean{}
	}
	return newError("Неизвестный тип поля: %s", variable.Type.(*ast.SymbolType).Name)
}

func EvaluateFunctionField(fn ast.ClassFunctionStatement, env *object.Environment, index string) object.Object {
	var params []ast.VariableDecStatement
	for _, param := range fn.Parameters {
		tmp := ast.VariableDecStatement{
			Names:         []string{index},
			AssignedValue: nil,
			Type:          param,
			IsConstant:    false,
		}
		params = append(params, tmp)
	}
	var returnTypes []object.ObjectType
	for _, _type := range fn.ReturnTypes {
		switch _type.(type) {
		case *ast.ArrayType:
			returnTypes = append(returnTypes, object.ARRAY)
		}
		if _type.(*ast.SymbolType).Name == "string" {
			returnTypes = append(returnTypes, object.STRING)
		} else if _type.(*ast.SymbolType).Name == "bool" {
			returnTypes = append(returnTypes, object.BOOLEAN)
		} else if _type.(*ast.SymbolType).Name == "int" {
			returnTypes = append(returnTypes, object.INTEGER)
		} else if _type.(*ast.SymbolType).Name == "array" {
			returnTypes = append(returnTypes, object.ARRAY)
		}
	}
	function, ok := env.Get(index)
	if !ok {
		return newError("Не найдено определение функции: %s", index)
	}
	body := function.(*object.FunctionLiteral).Body
	return &object.FunctionLiteral{
		Env:        env,
		Parameters: params,
		Body:       body,
		ReturnType: returnTypes,
	}
}

func EvaluateIf(node ast.IfStatement, env *object.Environment) object.Object {
	condition := Evaluate(node.Condition, env)
	if isTruthy(condition) {
		return Execute(node.ThenBlock, env)
	} else if node.ElseBlock != nil {
		return Execute(node.ElseBlock, env)
	} else {
		return NULL
	}
}

func ExecuteProgram(statements ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range statements.Statements {
		result = Execute(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}
	return result
}

func ExecuteBlock(block ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Execute(statement, env)

		if result != nil && result.Type() == object.RETURN_VALUE {
			return result
		}
	}
	return result
}

func ExecuteImportStat(stat ast.ImportStatement, env *object.Environment) object.Object {
	modulePath := stat.PackagePath
	file, err := os.Open(modulePath)
	if err != nil {
		return newError("Ошибка при нахождении пакета: %s", modulePath)
	}
	defer file.Close()
	input, err := io.ReadAll(file)
	if err != nil {
		return newError("Ошибка при чтении файла: %s", modulePath)
	}
	tokens := lexer.Tokenize(string(input))
	ast, err := parser.Parse(tokens)
	if err != nil {
		return newError("Ошибка при парсинге файла: %s", err)
	}
	enviroment := object.NewEnvironment()
	for _, stmt := range ast.Statements {
		Execute(stmt, enviroment)
	}
	module := &object.Module{
		Name:        stat.ImportName,
		Environment: *enviroment,
	}
	env.Set(stat.ImportName, module)

	return nil
}