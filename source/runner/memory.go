package runner

import "meow/source/ast"

type Memory struct {
	Variables map[string]interface{}
	Functions map[string]ast.FunctionDecStatement
	Classes   map[string]ast.ClassDecStatement
}

func NewMemory() *Memory {
	return &Memory{
		Variables: make(map[string]any),
		Functions: make(map[string]ast.FunctionDecStatement),
		Classes:   make(map[string]ast.ClassDecStatement),
	}
}

func (m *Memory) GetVariable(name string) (interface{}, bool) {
	return m.Variables[name], true
}

func (m *Memory) GetFunction(name string) (ast.FunctionDecStatement, bool) {
	return m.Functions[name], true
}

func (m *Memory) GetClass(name string) (interface{}, bool) {
	return m.Classes[name], true
}

func (m *Memory) SetVariable(name string, value any) {
	m.Variables[name] = value
}

func (m *Memory) SetFunction(name string, value ast.FunctionDecStatement) {
	m.Functions[name] = value
}

func (m *Memory) SetClass(name string, value ast.ClassDecStatement) {
	m.Classes[name] = value
}

func (m *Memory) GetAll() any {
	variables := make(map[string]any)
	for name, value := range m.Variables {
		variables[name] = value
	}
	for name, value := range m.Functions {
        variables[name] = value
    }
	for name, value := range m.Classes {
        variables[name] = value
    }

	return variables
}