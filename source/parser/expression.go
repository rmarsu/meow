package parser

import (
	"errors"
	"fmt"
	"meow/source/ast"
	"meow/source/ast/helper"
	"meow/source/lexer"
	"strconv"
)

func parseExpression(p *parser, bp binding_power) ast.Expression {
	tokenKind := lexer.GetTokenKind(p.getCurrToken())
	nud_func, exists := nud_lu[tokenKind]
	if !exists {
		p.errors = append(p.errors, fmt.Errorf("ожидалось значение: %s", lexer.TokenKindString(tokenKind)))
	}
	left := nud_func(p)
	for bp_lu[lexer.GetTokenKind(p.getCurrToken())] > bp {
		tokenKind = lexer.GetTokenKind(p.getCurrToken())
		led_func, exists := led_lu[tokenKind]
		if !exists {
			p.errors = append(p.errors, fmt.Errorf("ожидался оператор: %s", lexer.TokenKindString(tokenKind)))
			break
		}

		left = led_func(p, left, bp_lu[p.getCurrToken().Kind])
	}
	return left
}

func parsePrimaryExpressions(p *parser) ast.Expression {
	switch lexer.GetTokenKind(p.getCurrToken()) {
	case lexer.INT:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return &ast.NumberExpression{Value: number}
	case lexer.STRING:
		return &ast.StringExpression{Value: []rune(p.advance().Value)}
	case lexer.IDENT:
		return &ast.SymbolExpression{Value: p.advance().Value}
	default:
		p.errors = append(p.errors, errors.New("невозможно создать первичное выражение"))
		return nil
	}
}

func parseBinaryExpressions(p *parser, left ast.Expression, bp binding_power) ast.Expression {
	operator := p.advance()
	right := parseExpression(p, bp)
	return &ast.BOExpression{
		Left:  left,
		Op:    operator,
		Right: right,
	}
}

func parseAssignmentExpressions(p *parser, left ast.Expression, bp binding_power) ast.Expression {
	operatorToken := p.advance()
	rhs := parseExpression(p, bp)
	return &ast.AssignmentExpression{
		Assigne: left,
		Op:      operatorToken,
		Value:   rhs,
	}
}

func parsePrefixExpressions(p *parser) ast.Expression {
	operToken := p.advance()
	rhs := parsePrimaryExpressions(p)
	return &ast.PrefixExpression{
		Op:        operToken,
		RightExpr: rhs,
	}
}

func parseGroupingExpressions(p *parser) ast.Expression {
	p.advance()
	expr := parseExpression(p, default_power)
	p.expect(lexer.RPAR)
	return expr
}

func parseClassInstanceExpressions(p *parser) ast.Expression {
	p.expect(lexer.EXCLAMINATION_MARK)
	p.expect(lexer.EXCLAMINATION_MARK)
	var structName = p.expect(lexer.IDENT).Value
	var fiels = map[string]ast.Expression{}
	p.expect(lexer.LPAR)

	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		fieldName := p.expect(lexer.IDENT).Value
		p.expect(lexer.ASSIGN)
		expr := parseExpression(p, LOGICAL)

		fiels[fieldName] = expr

		if p.getCurrToken().Kind != lexer.RPAR {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.RPAR)
	return &ast.ClassInstance{
		ClassName: structName,
		Fields:    fiels,
	}
}

func parseArrayInstanceExpressions(p *parser, left ast.Expression, bp binding_power) ast.Expression {
	var content = []ast.Expression{}
	p.expect(lexer.LBRAK)
	if p.getCurrToken().Kind != lexer.RBRAK {
		for p.hasTokens() && p.getCurrToken().Kind != lexer.RBRAK {
			content = append(content, parseExpression(p, LOGICAL))
			if p.getCurrToken().Kind != lexer.RBRAK {
				p.expect(lexer.COMMA)
			}
		}
	} else {
	}
	p.expect(lexer.RBRAK)
	return &ast.ArrayInstance{
		Underlying: left,
		Content:    content,
	}
}

func parseFunctionInstanceExpression(p *parser, left ast.Expression, bp binding_power) ast.Expression {
	var functionName = helper.ExpectType[*ast.SymbolExpression](left).Value
	var parameters = []ast.Expression{}
	p.expect(lexer.LPAR)

	for p.hasTokens() && p.getCurrToken().Kind != lexer.RPAR {
		parameters = append(parameters, parseExpression(p, LOGICAL))
		if p.getCurrToken().Kind != lexer.RPAR {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.RPAR)
	return &ast.FunctionInstance{
		FunctionName: functionName,
		Parameters:   parameters,
	}
}

func parseMemberInstanceExpression(p *parser, left ast.Expression, bp binding_power) ast.Expression {
	p.expect(lexer.DOT)
	memberName := p.expect(lexer.IDENT).Value
	return &ast.MemberInstance{
		Instance:   left,
		MemberName: memberName,
	}
}
