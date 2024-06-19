package gen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/trwk76/goweb/openapi/spec"
	"gopkg.in/yaml.v3"
)

func New(path string, info spec.Info, sec spec.SecurityRequirements) *API {
	res := &API{
		spc: spec.OpenAPI{
			OpenAPI: spec.Version,
			Info:    info,
			Servers: []spec.Server{{URL: path, Description: "Current server"}},
			Paths:   make(spec.Paths),
			Components: &spec.Components{
				Schemas:         make(spec.NamedSchemas),
				Responses:       make(spec.NamedResponseOrRefs),
				Parameters:      make(spec.NamedParameterOrRefs),
				Examples:        make(spec.NamedExampleOrRefs),
				RequestBodies:   make(spec.NamedRequestBodyOrRefs),
				Headers:         make(spec.NamedHeaderOrRefs),
				SecuritySchemes: make(spec.NamedSecuritySchemeOrRefs),
			},
			Security: sec,
		},
		sch: make(schemas),
		paths: paths{
			named: make(map[string]*Path),
		},
	}

	return res
}

func (a *API) Generate(dir string, initName string, initPkg string) error {
	if err := os.MkdirAll(dir, os.FileMode(0777)); err != nil {
		return fmt.Errorf("error while creating output directory '%s': %s", dir, err.Error())
	}

	file, err := os.Create(filepath.Join(dir, initName))
	if err != nil {
		return fmt.Errorf("error creating output file '%s': %s", filepath.Join(dir, initName), err.Error())
	}

	w := bufio.NewWriter(file)

	defer func() {
		w.Flush()
		file.Close()
	}()


	writeSpec(dir, "json", a.spc, json.Marshal)
	writeSpec(dir, "yaml", a.spc, yaml.Marshal)

	return nil
}

type (
	API struct {
		spc   spec.OpenAPI
		sch   schemas
		paths paths
	}

	marshal func(v any) ([]byte, error)
)

func writeSpec(dir string, ext string, spc spec.OpenAPI, m marshal) error {
	path := filepath.Join(dir, "openapi." + ext)

	raw, err := m(spc)
	if err != nil {
		return fmt.Errorf("error marshaling openapi specification to %s: %s", ext, err.Error())
	}

	return os.WriteFile(path, raw, os.FileMode(0662))
}
