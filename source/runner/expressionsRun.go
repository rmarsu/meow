package runner

import (
	"meow/source/ast"
	"meow/source/lexer"
	"fmt"
)

func (r *Runner) evaluate(expr ast.Expression) any {
	switch expr := expr.(type) {
     case *ast.SymbolExpression:
          variable, _ := r.Memory.GetVariable(expr.Value)
          return variable
     case *ast.AssignmentExpression:
          r.Memory.SetVariable(expr.Assigne.(*ast.SymbolExpression).Value, r.evaluate(expr.Value))

          variable, _ := r.Memory.GetVariable(expr.Assigne.(*ast.SymbolExpression).Value)
          return variable
     case *ast.NumberExpression:
          return expr.Value
     case *ast.StringExpression:
          return expr.Value
     case *ast.PrefixExpression:
          switch expr.Op.Kind {
               case lexer.MINUS:
                    return r.evaluate(expr.RightExpr).(float64) * -1
               case lexer.PLUS:
                    return r.evaluate(expr.RightExpr)
               default:
                    r.Errors = append(r.Errors, fmt.Errorf("неподдерживаемый префиксный оператор: %s", expr.Op.Value))
                    return ""
          }
     
     case *ast.BOExpression:
          var left float64
          var right float64
          left = r.evaluate(expr.Left).(float64)
          right = r.evaluate(expr.Right).(float64)
          switch expr.Op.Kind {
          case lexer.PLUS:
               return left + right
          case lexer.MINUS:
               return left - right
          case lexer.MUL:
               return left * right
          case lexer.DIV:
               if right == 0 {
                    r.Errors = append(r.Errors, fmt.Errorf("деление на ноль"))
                    return ""
               }
               return left / right
          case lexer.NOT_EQUALS:
               return left != right
          case lexer.EQUALS:
               return left == right
          case lexer.LESS:
               return left < right
          case lexer.LESS_EQUALS:
               return left <= right
          case lexer.GREATER:
               return left > right
          case lexer.GREATER_EQUALS:
               return left >= right
          }
     case *ast.FunctionInstance:
          function, _ := r.Memory.GetFunction(expr.FunctionName)
          if function.Name == "" {
               r.Errors = append(r.Errors, fmt.Errorf("функция %s не найдена", expr.FunctionName))
               return ""
          }
          parameters := expr.Parameters
          calledParams := function.Parameters
          if len(parameters)!= len(calledParams) {
               r.Errors = append(r.Errors, fmt.Errorf("неверное число аргументов для функции %s", expr.FunctionName))
               return ""
          }
          for i := 0; i < len(parameters); i++ {
               r.Memory.SetVariable(calledParams[i].Name, r.evaluate(parameters[i]))
          }
          result := r.runBlockStatement(function.Body)
          if len(result) == 1 {
               return result[0]
           }
          if result != nil {
               return result
          }
     case *ast.ClassInstance:
          class, _ := r.Memory.GetClass(expr.ClassName)
          if class.(ast.ClassDecStatement).Name == "" {
               r.Errors = append(r.Errors, fmt.Errorf("класс %s не найден", expr.ClassName))
               return ""
          }
          instance := make(map[string]any)

          for key, value := range expr.Fields {
               instance[key] = r.evaluate(value)
          }
          r.Memory.SetVariable(expr.ClassName, instance)
          return instance
     }

     return nil
} 

func CompareTypes(value1 string , value2 any) bool {
     switch value2.(type) {
     case float64:
          if value1 == "float64" {
               return true
          }
          return false
     case string:
          if value1 == "string" {
               return true
          }
          return false
     case int:
          if value1 == "int" {
               return true
          }
          return false
     case bool:
          if value1 == "bool" {
               return true
          }
          return false
     }
     return false 
}