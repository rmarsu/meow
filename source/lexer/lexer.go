package lexer

import (
	"fmt"
	"regexp"
)

type regexPattern struct {
	pattern *regexp.Regexp
	handler regexHandler
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advance(len(value))
		lex.addTokens(NewToken(kind, value))
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.getReminder())
	lex.addTokens(NewToken(INT, match))
	lex.advance(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.getReminder())
	lex.advance(match[1])
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.getReminder())
	stringLiteral := lex.getReminder()[match[0]+1 : match[1]-1]
	lex.addTokens(NewToken(STRING, stringLiteral))
	lex.advance(len(stringLiteral) + 2)
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	value := regex.FindString(lex.getReminder())

	if kind, exists := reserved_lookup[value]; exists {
		lex.addTokens(NewToken(kind, value))
	} else {
		lex.addTokens(NewToken(IDENT, value))
	}
	lex.advance(len(value))
}

var defaultPatterns = []regexPattern{
	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
	{regexp.MustCompile(`\s+`), skipHandler},
	{regexp.MustCompile(`"[^"]*"`), stringHandler},
	{regexp.MustCompile(`\#\#.*`), skipHandler},
	{regexp.MustCompile(`[.]`), defaultHandler(DOT, ".")},
	{regexp.MustCompile(`\(`), defaultHandler(LPAR, "(")},
	{regexp.MustCompile(`\)`), defaultHandler(RPAR, ")")},
	{regexp.MustCompile(`\^`), defaultHandler(EXCLAMINATION_MARK, "!")},
	{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
	{regexp.MustCompile(`;`), defaultHandler(SEMICOLON, ";")},
	{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
	{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
	{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
	{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
	{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
	{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
	{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
	{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
	{regexp.MustCompile(`\-`), defaultHandler(MINUS, "-")},
	{regexp.MustCompile(`\-=`), defaultHandler(MINUS_EQUALS, "-=")},
	{regexp.MustCompile(`\*`), defaultHandler(MUL, "*")},
	{regexp.MustCompile(`\*\*=`), defaultHandler(MUL_EQUALS, "*=")},
	{regexp.MustCompile(`/`), defaultHandler(DIV, "/")},
	{regexp.MustCompile(`/=`), defaultHandler(DIV_EQUALS, "/=")},
	{regexp.MustCompile(`\=`), defaultHandler(ASSIGN, "=")},
	{regexp.MustCompile(`import`), defaultHandler(IMPORT, "import")},
	{regexp.MustCompile(`typeof`), defaultHandler(TYPEOF, "typeof")},
	{regexp.MustCompile(`function`), defaultHandler(FUNCTION, "function")},
	{regexp.MustCompile(`var`), defaultHandler(VAR, "var")},
	{regexp.MustCompile(`return`), defaultHandler(RETURN, "return")},
	{regexp.MustCompile(`if`), defaultHandler(IF, "if")},
	{regexp.MustCompile(`else`), defaultHandler(ELSE, "else")},
	{regexp.MustCompile(`while`), defaultHandler(WHILE, "while")},
	{regexp.MustCompile(`for`), defaultHandler(FOR, "for")},
	{regexp.MustCompile(`true`), defaultHandler(TRUE, "true")},
	{regexp.MustCompile(`false`), defaultHandler(FALSE, "false")},
	{regexp.MustCompile(`or`), defaultHandler(OR, "or")},
	{regexp.MustCompile(`and`), defaultHandler(AND, "and")},
	{regexp.MustCompile(`const`), defaultHandler(CONST, "const")},
	{regexp.MustCompile(`\n`), skipHandler},
	{regexp.MustCompile(`\[`), defaultHandler(LBRAK, "[")},
	{regexp.MustCompile(`\]`), defaultHandler(RBRAK, "]")},
	{regexp.MustCompile(`class`), defaultHandler(CLASS, "class")},
	{regexp.MustCompile(`typeof`), defaultHandler(TYPEOF, "typeof")},
	{regexp.MustCompile(`\{`), defaultHandler(LCURLY, "{")},
	{regexp.MustCompile(`\}`), defaultHandler(RCURLY, "}")},
	{regexp.MustCompile(`static`), defaultHandler(STATIC, "static")},
	{regexp.MustCompile(`void`), defaultHandler(VOID, "void")},
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	input    string
	currPos  int
}

func Tokenize(input string) []Token {
	lexer := NewLexer(defaultPatterns)
	lexer.input = input

	for !lexer.atTheEnd() {
		match := false
		for _, pattern := range lexer.patterns {
			loc := pattern.pattern.FindStringIndex(lexer.getReminder())
			if loc != nil && loc[0] == 0 {
				match = true
				pattern.handler(lexer, pattern.pattern)
				break
			}
		}
		if !match {
			panic(fmt.Sprintf("Нераспознанный символ возле: %d: %s", lexer.currPos, lexer.getReminder()))
		}
	}
	lexer.addTokens(NewToken(EOF, "EOF"))
	return lexer.Tokens
}

func NewLexer(patterns []regexPattern) *lexer {
	return &lexer{
		Tokens:   make([]Token, 0),
		patterns: patterns,
		currPos:  0,
		input:    "",
	}
}

func (lex *lexer) advance(count int) {
	lex.currPos += count
}

func (lex *lexer) addTokens(tokens Token) {
	lex.Tokens = append(lex.Tokens, tokens)
}

func (lex *lexer) getReminder() string {
	return lex.input[lex.currPos:]
}

func (lex *lexer) GetReminder() string {
	return lex.input[lex.currPos:]
}

func (lex *lexer) atTheEnd() bool {
	return lex.currPos >= len(lex.input)
}
