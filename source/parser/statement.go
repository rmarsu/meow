package parser

import (
	"meow/source/ast"
	"meow/source/lexer"
)

func parseStatement(p *parser) ast.Statement {
	statement_func, exists := statement_lu[lexer.GetTokenKind(p.getCurrToken())]
	if exists {
		return statement_func(p)
	}
	expression := parseExpression(p, default_power)
	p.expect(lexer.SEMICOLON)

	return &ast.ExpressionStatement{
		Expression: expression,
	}
}

func parseVariableDeclaration(p *parser) ast.Statement {
	var expilitType ast.Type
	var assigmentValue ast.Expression
	var IsConstant bool
	var names []string
	if p.getCurrToken().Kind == lexer.VAR {
		p.advance()
		IsConstant = false
	} else if p.getCurrToken().Kind == lexer.CONST {
		p.advance()
		IsConstant = true
	}
	for p.hasTokens() && p.getCurrToken().Kind != lexer.ASSIGN {
		name := p.getCurrToken().Value
		names = append(names, name)
		p.advance()
		if p.getCurrToken().Kind != lexer.ASSIGN {
			p.expect(lexer.COMMA)
		}
	}

	if p.getCurrToken().Kind == lexer.ASSIGN {
		p.advance()
		assigmentValue = parseExpression(p, ASSIGN)
	} else if p.getCurrToken().Kind == lexer.SEMICOLON {
		expilitType = parseType(p, ASSIGN)
		assigmentValue = nil
	}

	p.expect(lexer.SEMICOLON)

	return &ast.VariableDecStatement{
		Names:         names,
		IsConstant:    IsConstant,
		AssignedValue: assigmentValue,
		Type:          expilitType,
	}
}

func parseClassDeclaration(p *parser) ast.Statement {
	p.expect(lexer.CLASS)
	var fields = map[string]ast.ClassFieldStatement{}
	var functions = map[string]ast.ClassFunctionStatement{}
	className := p.expect(lexer.IDENT).Value

	p.expect(lexer.LPAR)
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		var isStatic bool
		var fieldName string
		if p.getCurrToken().Kind == lexer.STATIC {
			isStatic = true
			p.expect(lexer.STATIC)
		}
		if p.getCurrToken().Kind == lexer.IDENT {
			fieldType := parseType(p, default_power)
			fieldName = p.expect(lexer.IDENT).Value
			if p.getCurrToken().Kind == lexer.LPAR {
				p.advance()
				parameters := make([]ast.Type, 0)
				for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
					parameterType := parseType(p, PRIMARY)
					parameters = append(parameters, parameterType)
				}
				p.expect(lexer.RPAR)
				functionName := fieldName
				functions[functionName] = ast.ClassFunctionStatement{
					Parameters: parameters,
					ReturnType: fieldType,
					IsStatic:   isStatic,
				}
			}
			p.expect(lexer.COMMA)

			_, exists := fields[fieldName]
			if exists {
				panic("!! Данное поле уже было указано в классе")
			}
			fields[fieldName] = ast.ClassFieldStatement{
				Type:     fieldType,
				IsStatic: isStatic,
			}
			continue
		}
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.SEMICOLON)

	return &ast.ClassDecStatement{
		Name:      className,
		Fields:    fields,
		Functions: functions,
	}
}

func parseFunctionDeclaration(p *parser) ast.Statement {
	p.expect(lexer.VOID)
	functionName := p.expect(lexer.IDENT).Value
	p.expect(lexer.LPAR)
	var params = make([]ast.VariableDecStatement, 0)
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		var paramName string
		var paramType ast.Type
		paramName = p.expect(lexer.IDENT).Value
		paramType = parseType(p, default_power)
		if p.getCurrToken().Kind != lexer.RPAR {
			p.expect(lexer.COMMA)
		}
		var names []string
		names = append(names, paramName)
		params = append(params, ast.VariableDecStatement{
			Names:         names,
			IsConstant:    false,
			AssignedValue: nil,
			Type:          paramType,
		})
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.LPAR)
	var returnValues []ast.Type
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		returnValues = append(returnValues, parseType(p, default_power))
		if p.getCurrToken().Kind != lexer.RPAR {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.LPAR)
	var body []ast.Statement
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		body = append(body, parseStatement(p))
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.SEMICOLON)
	return &ast.FunctionDecStatement{
		Name:       functionName,
		Parameters: params,
		ReturnType: returnValues,
		Body: &ast.BlockStatement{
			Statements: body,
		},
	}
}

func parseReturnStatement(p *parser) ast.Statement {
	p.expect(lexer.RETURN)
	var expressions []ast.Expression
	for p.hasTokens() && p.getCurrToken().Kind != lexer.SEMICOLON {
		expressions = append(expressions, parseExpression(p, default_power))
		if p.getCurrToken().Kind != lexer.SEMICOLON {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.SEMICOLON)
	return &ast.ReturnStatement{
		Expressions: expressions,
	}
}

func parseIfStatement(p *parser) ast.Statement {
	p.expect(lexer.IF)
	p.expect(lexer.LPAR)
	condition := parseExpression(p, LOGICAL)
	p.expect(lexer.RPAR)
	p.expect(lexer.LPAR)
	var thenBranch []ast.Statement
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		thenBranch = append(thenBranch, parseStatement(p))
	}
	p.expect(lexer.RPAR)
	var elseBranch []ast.Statement
	if p.getCurrToken().Kind == lexer.ELSE {
		p.advance()
		p.expect(lexer.LPAR)
		for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
			elseBranch = append(elseBranch, parseStatement(p))
		}
		p.expect(lexer.RPAR)
	}
	p.expect(lexer.SEMICOLON)
	return &ast.IfStatement{
		Condition: condition,
		ThenBlock: &ast.BlockStatement{
			Statements: thenBranch,
		},
		ElseBlock: &ast.BlockStatement{
			Statements: elseBranch,
		},
	}
}

func parseWhileStatement(p *parser) ast.Statement {
	p.expect(lexer.FOR)
	p.expect(lexer.LPAR)
	var conditions []ast.Expression
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		conditions = append(conditions, parseExpression(p, LOGICAL))
		if p.getCurrToken().Kind != lexer.RPAR {
			p.expect(lexer.SEMICOLON)
		}
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.LPAR)
	var body []ast.Statement
	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		body = append(body, parseStatement(p))
	}
	p.expect(lexer.RPAR)
	p.expect(lexer.SEMICOLON)
	return &ast.WhileStatement{
		Conditions: conditions,
		Body: &ast.BlockStatement{
			Statements: body,
		},
	}
}

func parseImportStatement(p *parser) ast.Statement {
	p.expect(lexer.IMPORT)
	name := p.expect(lexer.IDENT).Value
	path := p.expect(lexer.STRING).Value
	p.expect(lexer.SEMICOLON)
	return &ast.ImportStatement{
		ImportName:  name,
		PackagePath: path,
	}
}
