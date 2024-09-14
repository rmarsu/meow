package lexer

import "fmt"

type TokenKind int

const (
	// Special tokens
	ILLEGAL TokenKind = iota
	EOF

	// Identifiers and literals
	IDENT
	INT
	FLOAT
	STRING

	// Operators and delimiters
	DOT
	ASSIGN
	PLUS
	PLUS_EQUALS
	MINUS
	MINUS_EQUALS
	MUL
	MUL_EQUALS
	DIV
	DIV_EQUALS
	LPAR
	RPAR
	SEMICOLON
	COMMA
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS
	EQUALS
	NOT
	NOT_EQUALS

	// Keywords
	FUNCTION
	VAR
	RETURN
	IF
	ELSE
	WHILE
	FOR
	TRUE
	FALSE
	OR
	AND
	IMPORT
	CLASS
	TYPEOF
	MEOW
	CONST
	LBRAK
	RBRAK
	STATIC
	PUBLIC
	PRIVATE

	LCURLY
	RCURLY

	VOID
	EXCLAMINATION_MARK
)

var reserved_lookup map[string]TokenKind = map[string]TokenKind{
	"assign":   ASSIGN,
	"function": FUNCTION,
	"var":      VAR,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"true":     TRUE,
	"false":    FALSE,
	"or":       OR,
	"and":      AND,
	"import":   IMPORT,
	"typeof":   TYPEOF,
	"meow":     MEOW,
	"const":    CONST,
	"class":    CLASS,
	"static":   STATIC,
	"public":   PUBLIC,
	"private":  PRIVATE,
	"void":     VOID,
	"!":        EXCLAMINATION_MARK,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func NewToken(kind TokenKind, value string) Token {
	return Token{Kind: kind, Value: value}
}

func (token Token) Debug() {
	fmt.Printf("%s: %s\n", TokenKindString(token.Kind), token.Value)
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case PLUS_EQUALS:
		return "PLUS_EQUALS"
	case MINUS:
		return "MINUS"
	case MINUS_EQUALS:
		return "MINUS_EQUALS"
	case MUL:
		return "MUL"
	case MUL_EQUALS:
		return "MUL_EQUALS"
	case DIV:
		return "DIV"
	case DIV_EQUALS:
		return "DIV_EQUALS"
	case LPAR:
		return "LPAR"
	case RPAR:
		return "RPAR"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case LESS:
		return "LESS"
	case LESS_EQUALS:
		return "LESS_EQUALS"
	case GREATER:
		return "GREATER"
	case GREATER_EQUALS:
		return "GREATER_EQUALS"
	case EQUALS:
		return "EQUALS"
	case NOT:
		return "NOT"
	case NOT_EQUALS:
		return "NOT_EQUALS"
	case FUNCTION:
		return "FUNCTION"
	case VAR:
		return "VAR"
	case RETURN:
		return "RETURN"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case OR:
		return "OR"
	case AND:
		return "AND"
	case IMPORT:
		return "IMPORT"
	case TYPEOF:
		return "TYPEOF"
	case MEOW:
		return "MEOW"
	case CONST:
		return "CONST"
	case LBRAK:
		return "LBRAK"
	case RBRAK:
		return "RBRAK"
	case CLASS:
		return "CLASS"
	case STATIC:
		return "STATIC"
	case PUBLIC:
		return "PUBLIC"
	case PRIVATE:
		return "PRIVATE"
	case LCURLY:
		return "LCURLY"
	case RCURLY:
		return "RCURLY"
	case VOID:
		return "VOID"
	case EXCLAMINATION_MARK:
		return "EXCLAMINATION_MARK"
	case DOT:
		return "DOT"
	}
	return "UNKNOWN"
}

func GetTokenKind(token Token) TokenKind {
	return token.Kind
}
