package runner

import "meow/source/ast"

func (r *Runner) GetPackage(filepath string) *Package {
	return r.Packages[filepath]
}

func (r *Runner) RegisterPackage(filepath string) *Package {
	pkg := &Package{
		Memory: Memory{
			Variables:        make(map[string]*ast.VariableDecStatement),
			Functions:        make(map[string]*ast.FunctionDecStatement),
			Classes:          make(map[string]*ast.ClassDecStatement),
			ClassesInstances: make(map[string]*ast.ClassInstance),
		},
	}
	r.Packages[filepath] = pkg
	return pkg
}

func (r *Runner) RegisterVariable(pkg *Package, varDec *ast.VariableDecStatement) {
	pkg.Memory.Variables[varDec.Name] = varDec
}

func (r *Runner) GetVariable(pkg *Package, name string) *ast.VariableDecStatement {
	return pkg.Memory.Variables[name]
}

func (r *Runner) RegisterFunction(pkg *Package, funcDec *ast.FunctionDecStatement) {
	pkg.Memory.Functions[funcDec.Name] = funcDec
}

func (r *Runner) GetFunction(pkg *Package, name string) *ast.FunctionDecStatement {
	return pkg.Memory.Functions[name]
}

func (r *Runner) RegisterClass(pkg *Package, classDec *ast.ClassDecStatement) {
	pkg.Memory.Classes[classDec.Name] = classDec
}

func (r *Runner) GetClass(pkg *Package, name string) *ast.ClassDecStatement {
	return pkg.Memory.Classes[name]
}

func (r *Runner) RegisterClassInstance(pkg *Package, classDec *ast.ClassInstance) {
	pkg.Memory.ClassesInstances[classDec.ClassName] = classDec
}

func (r *Runner) GetClassInstance(pkg *Package, name string) *ast.ClassInstance {
	return pkg.Memory.ClassesInstances[name]
}

func (r *Runner) DeleteFromTempVariable(name string) {
	delete(r.TemporaryMemory.Variables, name)
}

func (r *Runner) DeleteFromTempFunction(name string) {
	delete(r.TemporaryMemory.Functions, name)
}

func (r *Runner) DeleteFromTempClass(name string) {
	delete(r.TemporaryMemory.Classes, name)
}

func (r *Runner) DeleteFromTempClassInstance(name string) {
	delete(r.TemporaryMemory.ClassesInstances, name)
}
