package parser

import (
	"meow/source/ast"
	"meow/source/lexer"
	"fmt"
)

type parser struct {
	errors []error
	tokens []lexer.Token
	currPos int
}

func NewParser(tokens []lexer.Token) *parser {
	createTokenLookups()
	createTokenTypeLookups()
	return &parser{tokens: tokens, currPos: 0}
}


func Parse(tokens []lexer.Token) (ast.BlockStatement , error) {
	body := make([]ast.Statement, 0)
	parser := NewParser(tokens)
	for parser.hasTokens() {
		body = append(body, parseStatement(parser))
	}
	return ast.BlockStatement{
		Statements: body,
	}, nil
}

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.getCurrToken()
	kind := token.Kind
	if kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Ожидался '%s', но получен '%s' возле %s", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind), p.getCurrToken().Value)
		}
		panic(err)
		
     }
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
     return p.expectError(expectedKind, nil)
}