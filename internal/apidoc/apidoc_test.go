// Package apidoc checks that the public API is documented.
package apidoc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportedIdentifiersAreDocumented fails for an exported identifier of pkg/ without a doc comment, and for a
// package without a package comment.
func TestExportedIdentifiersAreDocumented(t *testing.T) {
	fset := token.NewFileSet()
	packageDoc := map[string]bool{}
	err := filepath.Walk(filepath.Join("..", "..", "pkg"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		dir := filepath.Dir(path)
		packageDoc[dir] = packageDoc[dir] || file.Doc != nil
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Name.IsExported() && d.Doc == nil && exportedReceiver(d) {
					t.Errorf("%s: %s has no doc comment", fset.Position(d.Pos()), d.Name.Name)
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() && s.Doc == nil && d.Doc == nil {
							t.Errorf("%s: %s has no doc comment", fset.Position(s.Pos()), s.Name.Name)
						}
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.IsExported() && s.Doc == nil && d.Doc == nil {
								t.Errorf("%s: %s has no doc comment", fset.Position(name.Pos()), name.Name)
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(packageDoc) == 0 {
		t.Fatal("no packages found")
	}
	for dir, documented := range packageDoc {
		if !documented {
			t.Errorf("%s: the package has no package comment", dir)
		}
	}
}

// exportedReceiver reports whether a function is not a method or is a method of an exported type.
func exportedReceiver(d *ast.FuncDecl) bool {
	if d.Recv == nil {
		return true
	}
	t := d.Recv.List[0].Type
	if star, ok := t.(*ast.StarExpr); ok {
		t = star.X
	}
	if index, ok := t.(*ast.IndexExpr); ok {
		t = index.X
	}
	id, ok := t.(*ast.Ident)
	return !ok || id.IsExported()
}
