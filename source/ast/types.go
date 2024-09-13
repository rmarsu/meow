package ast

type SymbolType struct {
	Name string
}

func (st SymbolType) func_type() {}

type ArrayType struct {
	Underlying Type
}

func (at ArrayType) func_type() {}
