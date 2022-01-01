package brownfox

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/JodeZer/lazydog/file"
	"github.com/JodeZer/lazydog/inject"
)

type BrownFox struct {
	path   string
	dirs   []string
	depth  int
	inject *inject.Injector
	file.Jumper
}

func NewBrownFox(path string, depth int) *BrownFox {
	return &BrownFox{
		path:   path,
		depth:  depth,
		inject: inject.NewInjector(),
		dirs:   file.TreeDir(path, depth),
	}
}

func (b *BrownFox) Backup() error {
	for _, dir := range b.dirs {
		if err := b.BackupPath(dir); err != nil {
			return err
		}
	}
	return nil
}

func (b *BrownFox) Restore() error {
	for _, dir := range b.dirs {
		if err := b.RestorePath(dir); err != nil {
			return err
		}
		goFiles := file.ListGoFile(dir, false)
		if len(goFiles) == 0 {
			continue
		}
		parser := inject.NewParser(token.NewFileSet(), goFiles[0])
		if err := parser.Parse(); err != nil {
			return err
		}
		dgHelper := inject.NewDogHelper(dir, parser.PkgName())
		if err := dgHelper.EraseDogHelper(); err != nil && !os.IsNotExist(err) {
			return err
		}

	}

	return nil
}

func (b *BrownFox) Inject() error {
	for _, dir := range b.dirs {
		for _, goFile := range file.ListGoFile(dir, false) {
			fmt.Println("Inject: ", goFile)
			fset := token.NewFileSet()
			parser := inject.NewParser(fset, goFile)
			err := parser.Parse()
			if err != nil {
				panic(err)
			}

			// inject
			parser.ForEachDecl(func(decl ast.Decl) {
				b.inject.InjectFunc(decl)
			})

			// write to new file
			var buf bytes.Buffer
			if err := printer.Fprint(&buf, fset, parser.GetAst()); err != nil {
				return err
			}

			//fmt.Println(buf.String())
			if err := ioutil.WriteFile(goFile, buf.Bytes(), os.ModeExclusive); err != nil {
				panic(err)
			}

			// write helper
			dgHelper := inject.NewDogHelper(dir, parser.PkgName())
			dgHelper.WriteDogHelper()
		}
	}

	return nil
}
