package parser

import (
	"meow/source/ast"
	"meow/source/lexer"
	// "fmt"
)

type type_nudHandler func(p *parser) ast.Type
type type_ledHandler func(p *parser, left ast.Type, bp binding_power) ast.Type

// lookup tables
type type_nudLookupTable map[lexer.TokenKind]type_nudHandler
type type_ledLookupTable map[lexer.TokenKind]type_ledHandler
type type_bpLookupTable map[lexer.TokenKind]binding_power

var type_bp_lu = bpLookupTable{}
var type_nud_lu = type_nudLookupTable{}
var type_led_lu = type_ledLookupTable{}

func type_led(kind lexer.TokenKind, bp binding_power, led_function type_ledHandler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_function
}

func type_nud(kind lexer.TokenKind, nud_function type_nudHandler) {
	type_nud_lu[kind] = nud_function
}

func createTokenTypeLookups() {
	type_nud(lexer.IDENT, parseSymbolType)
	type_nud(lexer.LBRAK, parseArrayType)

}

func parseSymbolType(p *parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.IDENT).Value,
	}
}

func parseArrayType(p *parser) ast.Type {
     p.advance() // eat '['
	p.expect(lexer.RBRAK)
     innerType := parseType(p, PRIMARY)
	// fmt.Println("returned")
     return ast.ArrayType{
          Underlying: innerType,
     }
}

func parseType(p *parser, bp binding_power) ast.Type {
	tokenKind := lexer.GetTokenKind(p.getCurrToken())
	// fmt.Println(lexer.TokenKindString(tokenKind))
	nud_func, exists := type_nud_lu[tokenKind]
	if !exists {
		// panic("duplicate type")
		return nil
	}
	left := nud_func(p)
	for bp_lu[lexer.GetTokenKind(p.getCurrToken())] > bp {
		tokenKind = lexer.GetTokenKind(p.getCurrToken())
		led_func, exists := type_led_lu[tokenKind]
		if !exists {
			return nil
		}

		left = led_func(p, left, bp_lu[p.getCurrToken().Kind])
	}
	return left
}