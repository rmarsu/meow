package parser

import (
	"meow/source/ast"
	"meow/source/lexer"
)

type binding_power int

const (
	default_power binding_power = iota
	COMMA
	ASSIGN
	LOGICAL
	RELATIONAL
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	CALL
	MEMBER
	PRIMARY
)

type statementHandler func(p *parser) ast.Statement
type nudHandler func(p *parser) ast.Expression
type ledHandler func(p *parser, left ast.Expression, bp binding_power) ast.Expression

// lookup tables
type statementLookupTable map[lexer.TokenKind]statementHandler
type nudLookupTable map[lexer.TokenKind]nudHandler
type ledLookupTable map[lexer.TokenKind]ledHandler
type bpLookupTable map[lexer.TokenKind]binding_power

var bp_lu = bpLookupTable{}
var nud_lu = nudLookupTable{}
var led_lu = ledLookupTable{}
var statement_lu = statementLookupTable{}

func led(kind lexer.TokenKind, bp binding_power, led_function ledHandler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_function
}

func nud(kind lexer.TokenKind, nud_function nudHandler) {
	nud_lu[kind] = nud_function
}

func statement(kind lexer.TokenKind, statement_function statementHandler) {
	statement_lu[kind] = statement_function
}

func createTokenLookups() {
	nud(lexer.MINUS, parsePrefixExpressions)
	led(lexer.ASSIGN, ASSIGN, parseAssignmentExpressions)
	led(lexer.PLUS_EQUALS, ASSIGN, parseAssignmentExpressions)
	led(lexer.MINUS_EQUALS, ASSIGN, parseAssignmentExpressions)
	led(lexer.MUL_EQUALS, ASSIGN, parseAssignmentExpressions)
	led(lexer.DIV_EQUALS, ASSIGN, parseAssignmentExpressions)

	led(lexer.AND, LOGICAL, parseBinaryExpressions)
	led(lexer.OR, LOGICAL, parseBinaryExpressions)

	led(lexer.LESS, RELATIONAL, parseBinaryExpressions)
	led(lexer.LESS_EQUALS, RELATIONAL, parseBinaryExpressions)
	led(lexer.GREATER, RELATIONAL, parseBinaryExpressions)
	led(lexer.GREATER_EQUALS, RELATIONAL, parseBinaryExpressions)
	led(lexer.EQUALS, RELATIONAL, parseBinaryExpressions)
	led(lexer.NOT_EQUALS, RELATIONAL, parseBinaryExpressions)

	led(lexer.PLUS, ADDITIVE, parseBinaryExpressions)
	led(lexer.MINUS, ADDITIVE, parseBinaryExpressions)
	led(lexer.MUL, MULTIPLICATIVE, parseBinaryExpressions)
	led(lexer.DIV, MULTIPLICATIVE, parseBinaryExpressions)

	nud(lexer.INT, parsePrimaryExpressions)
	nud(lexer.STRING, parsePrimaryExpressions)
	nud(lexer.IDENT, parsePrimaryExpressions)
	nud(lexer.LPAR, parseGroupingExpressions)

	nud(lexer.LBRAK, parseArrayInstanceExpressions)
	nud(lexer.EXCLAMINATION_MARK, parseClassInstanceExpressions)
	led(lexer.LPAR, CALL, parseFunctionInstanceExpression)

	statement(lexer.CONST, parseVariableDeclaration)
	statement(lexer.VAR, parseVariableDeclaration)
	statement(lexer.CLASS, parseClassDeclaration)
	statement(lexer.VOID, parseFunctionDeclaration)
	statement(lexer.RETURN, parseReturnStatement)
	statement(lexer.IF, parseIfStatement)
	statement(lexer.FOR, parseWhileStatement)
	statement(lexer.MEOW, parsePrintStatement)
	statement(lexer.IMPORT, parseImportStatement)
}
