package runner

import "meow/source/ast"

func (r *Runner) GetPackage(filepath string) *Package {
    return r.Packages[filepath]
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
