package runner

import "meow/source/ast"

func (r *Runner) initPackage(packagename string) *Package {
	var pkg *Package
	if packagename == "main" {
		pkg = r.Packages["main"]
		if pkg == nil {
			pkg = &Package{
				IsMain: true,
				Memory: Memory{
					Variables:        make(map[string]*ast.VariableDecStatement),
					Functions:        make(map[string]*ast.FunctionDecStatement),
					Classes:          make(map[string]*ast.ClassDecStatement),
					ClassesInstances: make(map[string]*ast.ClassInstance),
				},
			}
			r.Packages["main"] = pkg
		}
	} else {
		pkg = r.Packages[packagename]
		if pkg == nil {
			pkg = &Package{
				IsMain: false,
				Memory: Memory{
					Variables:        make(map[string]*ast.VariableDecStatement),
					Functions:        make(map[string]*ast.FunctionDecStatement),
					Classes:          make(map[string]*ast.ClassDecStatement),
					ClassesInstances: make(map[string]*ast.ClassInstance),
				},
			}
			r.Packages[packagename] = pkg
		}
	}
	return pkg
}

func (r *Runner) MainPackage() *Package {
	return r.Packages["main"]
}
