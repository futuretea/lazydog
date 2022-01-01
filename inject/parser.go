package inject

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type Parser struct {
	filename string
	astF     *ast.File
	tokenFS  *token.FileSet
}

func NewParser(tokenFS *token.FileSet, fn string) *Parser {
	return &Parser{
		filename: fn,
		tokenFS:  tokenFS,
	}
}

func (p *Parser) Parse() error {
	fBytes, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}
	index := strings.LastIndex(p.filename, `/`)
	f, err := parser.ParseFile(p.tokenFS, string(p.filename[index+1:]), fBytes, 0)
	if err != nil {
		return err
	}
	p.astF = f
	return nil
}

func (p *Parser) PkgName() string {
	return p.astF.Name.Name
}

func (p *Parser) ForEachDecl(f func(ast.Decl)) {
	for _, decl := range p.astF.Decls {
		f(decl)
	}
}

func (p *Parser) GetAst() *ast.File {
	return p.astF
}
