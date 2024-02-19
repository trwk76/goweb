package gen

import (
	"fmt"
	"os"
	"path/filepath"

	code "github.com/trwk76/gocode"
	"github.com/trwk76/gocode/golang"
	"github.com/trwk76/goweb/api/spec"
)

func New(path string, dir string, pkgName string) *APISpec {
	return &APISpec{
		path: path,
		dir:  dir,
		file: golang.NewFile("", pkgName),
	}
}

type (
	APISpec struct {
		Info      spec.Info
		Auths     AuthSpecs
		Paths     PathSpecs
		Schemas   SchemaSpecs
		Params    ParamSpecs
		ReqBodies ReqBodySpecs
		Resps     RespSpecs
		Security  spec.SecurityRequirements
		Tags      spec.Tags

		path string
		dir  string
		file golang.File
	}
)

func (s *APISpec) Generate() error {
	if err := os.MkdirAll(s.dir, os.FileMode(0777)); err != nil {
		return fmt.Errorf("error creating target directory '%s': %s", s.dir, err.Error())
	}

	stmts := make([]golang.Stmt, 0)

	stmts = append(stmts, golang.AssignStmt{
		Op:   golang.DeclAssign,
		Dest: golang.Exprs{golang.SymbolExpr{Name: "root"}},
		Src: golang.Exprs{golang.CallExpr{
			Func: golang.MemberExpr{
				Expr: golang.SymbolExpr{Name: "srv"},
				Name: "Path",
			},
			Args: golang.Args{
				golang.StringExpr(s.path),
			},
		}},
	})

	stmts = append(stmts, s.Paths.generate()...)

	s.file.Add(golang.FuncDecl{
		Name:   "InitializeAPI",
		Params: golang.Params{{Name: "srv", Type: golang.PtrType{Target: s.file.NamedType("Server", webPkg, nil)}}},
		Body:   golang.BlockStmt(stmts),
	})

	spc := s.spec()

	if err := os.WriteFile(filepath.Join(s.dir, "openapi.json"), spc.JSON(), os.FileMode(0644)); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(s.dir, "openapi.yaml"), spc.YAML(), os.FileMode(0644)); err != nil {
		return err
	}

	if err := code.WriteFile(filepath.Join(s.dir, "api_gen.go"), "\t", func(w *code.Writer) { s.file.Render(w) }); err != nil {
		return err
	}

	return nil
}

func (s *APISpec) spec() spec.Spec {
	return spec.Spec{
		OpenAPI: spec.OpenAPI,
		Info:    s.Info,
		Paths:   s.Paths.spec(),
		Components: &spec.Components{
			Schemas:         s.Schemas.spec(),
			Parameters:      s.Params.spec(),
			RequestBodies:   s.ReqBodies.spec(),
			Responses:       s.Resps.spec(),
			SecuritySchemes: s.Auths.spec(),
		},
		Servers:  spec.Servers{{URL: s.path, Description: "Current server."}},
		Security: s.Security,
		Tags:     s.Tags,
	}
}
