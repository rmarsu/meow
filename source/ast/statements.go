package ast

type BlockStatement struct {
	Statements []Statement
}

func (bs BlockStatement) statement() {

}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) statement() {

}

type VariableDecStatement struct {
	Names          []string
	IsConstant    bool
	Type          Type
	AssignedValue Expression
}

func (vds VariableDecStatement) statement() {}

type ClassFieldStatement struct {
	IsStatic bool
	Type     Type
}

type ClassFunctionStatement struct {
	Parameters []Type
	ReturnType Type
	IsStatic   bool
}

type ClassDecStatement struct {
	Name      string
	Fields    map[string]ClassFieldStatement
	Functions map[string]ClassFunctionStatement
}

func (cds ClassDecStatement) statement() {}

type FunctionDecStatement struct {
	Name       string
	Parameters []VariableDecStatement
	ReturnType []Type
	Body       *BlockStatement
}

func (fds FunctionDecStatement) statement() {}

type ReturnStatement struct {
	Expressions []Expression
}

func (rs ReturnStatement) statement() {}

type IfStatement struct {
	Condition Expression
	ThenBlock *BlockStatement
	ElseBlock *BlockStatement
}

func (is IfStatement) statement() {}

type WhileStatement struct {
	Conditions []Expression
	Body       *BlockStatement
}

func (ws WhileStatement) statement() {}


type ImportStatement struct {
	ImportName  string
	PackagePath string
}

func (is ImportStatement) statement() {}
