package ast

import "meow/source/lexer"

// number
type NumberExpression struct {
	Value float64
}

func (n NumberExpression) expression() {

}

// string
type StringExpression struct {
	Value []rune
}

func (s StringExpression) expression() {

}

// symbol
type SymbolExpression struct {
	Value string
}

func (s SymbolExpression) expression() {

}

// binary expression
type BOExpression struct {
	Left  Expression
	Right Expression
	Op    lexer.Token
}

func (bo BOExpression) expression() {

}

type PrefixExpression struct {
	Op        lexer.Token
	RightExpr Expression
}

func (pe PrefixExpression) expression() {}

type AssignmentExpression struct {
	Assigne Expression
	Op      lexer.Token
	Value   Expression
}

func (ae AssignmentExpression) expression() {}

type ClassInstance struct {
	ClassName string
	Fields    map[string]Expression
}

func (ci ClassInstance) expression() {}

type ArrayInstance struct {
	Underlying Expression
	Content    []Expression
}

func (ai ArrayInstance) expression() {}

type FunctionInstance struct {
	FunctionName string
	Parameters   []Expression
}

func (fi FunctionInstance) expression() {}

type MemberInstance struct {
	Instance Expression
	MemberName string
}

func (mi MemberInstance) expression() {}