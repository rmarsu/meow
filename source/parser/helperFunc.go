package parser

import (
     "meow/source/lexer"
)

func (p *parser) getCurrToken() lexer.Token {
	return p.tokens[p.currPos]
}



func (p *parser) advance() lexer.Token {
	tmp := p.getCurrToken()
	p.currPos++
	return tmp
}

func (p *parser) hasTokens() bool {
	return p.currPos < len(p.tokens) && lexer.GetTokenKind(p.getCurrToken()) != lexer.EOF 
}

func (p *parser) peek(num1 int) lexer.TokenKind {
	if!p.hasTokens() {
          return 0
	}
	return lexer.GetTokenKind(p.tokens[p.currPos + num1])
}



