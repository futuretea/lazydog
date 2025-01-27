package inject

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

const INJECT = `
package main

//import "fmt"

func a(){
	__traceStack()
}
`

type Injector struct {
	ImportSpec *ast.ImportSpec
	Stmt       ast.Stmt
}

func (i *Injector) InjectFunc(f ast.Decl) error {

	fd, ok := f.(*ast.FuncDecl)
	if !ok {
		return fmt.Errorf("not func")
	}
	newList := make([]ast.Stmt, 0, len(fd.Body.List)+1)

	newList = append(newList, i.Stmt)

	newList = append(newList, fd.Body.List...)

	fd.Body.List = newList
	return nil
}

func (i *Injector) InjectFile(path string) error {
	fSet := token.NewFileSet()
	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	index := strings.LastIndex(path, `/`)
	f, err := parser.ParseFile(fSet, path[index+1:], fBytes, 0)
	if err != nil {
		return err
	}

	for _, decl := range f.Decls {
		if err := i.InjectFunc(decl); err != nil {
			return err
		}
	}

	return nil
}

func NewInjector() *Injector {
	i := &Injector{}
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", INJECT, 0)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok {
			i.Stmt = fd.Body.List[0]
		}
	}
	return i
}
