package object

import (
	"bytes"
	"fmt"
	"meow/source/ast"
	"strings"

	"github.com/sanity-io/litter"
)

type ObjectType string

const (
	INTEGER      ObjectType = "INTEGER"
	BOOLEAN      ObjectType = "BOOLEAN"
	NULL         ObjectType = "NULL"
	STRING       ObjectType = "STRING"
	RETURN_VALUE ObjectType = "RETURN_VALUE"
	ERROR        ObjectType = "ERROR"
	FUNCTION     ObjectType = "FUNCTION"
	ARRAY        ObjectType = "ARRAY"
	CLASS        ObjectType = "CLASS"
	MODULE       ObjectType = "MODULE"
	FLOAT        ObjectType = "FLOAT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f *Float) Type() ObjectType {
	return FLOAT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "^^null^^"
}

type String struct {
	Value []rune
}

func (s *String) Type() ObjectType {
	return STRING
}

func (s *String) Inspect() string {
	return string(s.Value)
}

type ReturnValue struct {
	Values []Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE
}

func (r *ReturnValue) Inspect() string {
	var out bytes.Buffer
	for _, v := range r.Values {
		out.WriteString(v.Inspect())
		out.WriteString(" ")
	}
	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

type FunctionLiteral struct {
	Env        *Environment
	Parameters []ast.VariableDecStatement
	ReturnType []ObjectType
	Body       *ast.BlockStatement
}

func (fl *FunctionLiteral) Type() ObjectType {
	return FUNCTION
}

func (fl *FunctionLiteral) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.Names...)
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(litter.Sdump(fl.Body))
	out.WriteString("}")

	return out.String()
}

type Array struct {
	Elements     []Object
	ElementsType ObjectType
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Class struct {
	OriginName string
	Name       string
	Fields     map[string]Object
	Functions  map[string]Object
}

func (c *Class) Type() ObjectType {
	return CLASS
}

func (c *Class) Inspect() string {
	var out bytes.Buffer
	out.WriteString("class " + c.Name + " {\n")
	for _, f := range c.Fields {
		out.WriteString(f.Inspect() + "\n")
	}
	for _, f := range c.Functions {
		out.WriteString(f.Inspect() + "\n")
	}
	out.WriteString("}")

	return out.String()
}

type Field struct {
	FieldType ObjectType
	IsStatic  bool
}

func (f Field) Type() ObjectType {
	return f.FieldType
}

func (f Field) Inspect() string {
	return fmt.Sprintf("%s: %t", f.FieldType, f.IsStatic)
}

type Instance struct {
	obj Object
}

func (i Instance) Type() ObjectType {
	return i.obj.Type()
}

func (i Instance) Inspect() string {
	return i.obj.Inspect()
}

type Module struct {
	Name string
	Environment
}

func (m Module) Inspect() string {
	return "Module"
}

func (m Module) Type() ObjectType {
	return MODULE
}
