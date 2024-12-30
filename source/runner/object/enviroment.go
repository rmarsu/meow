package object

type Environment struct {
	errors []Error
	store  map[string]Object
	outer  *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{
		store:  s,
		outer:  nil,
		errors: make([]Error, 0),
	}
}

func (e *Environment) GetStore() map[string]Object {
	return e.store
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Delete(name string) Object {
	obj, ok := e.store[name]
	if ok {
		delete(e.store, name)
	}
	return obj
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}


func (e *Environment) Errors() []Error {
	return e.errors
}
