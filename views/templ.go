package views

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewFSTemplates(f fs.FS, partsDir string) FSTemplates {
	if !strings.HasSuffix(partsDir, "/") {
		partsDir += "/"
	}

	return FSTemplates{f: f, parts: partsDir}
}

func NewCachedTemplates(f fs.FS, partsDir string) CachedTemplates {
	res := CachedTemplates{p: make(map[string]*template.Template)}

	parseTemplateDir(f, "", partsDir, res.p)

	return res
}

type (
	Templates interface {
		Get(path string) *template.Template
	}

	FSTemplates struct {
		f     fs.FS
		parts string
	}

	CachedTemplates struct {
		p map[string]*template.Template
	}
)

func (t FSTemplates) Get(path string) *template.Template {
	return template.Must(template.ParseFS(t.f, path+".html", t.parts+"*.html"))
}

func (t CachedTemplates) Get(path string) *template.Template {
	return t.p[path]
}

func parseTemplateDir(f fs.FS, path string, partsDir string, res map[string]*template.Template) {
	ents, err := fs.ReadDir(f, path)
	if err != nil {
		panic(err)
	}

	for _, ent := range ents {
		if ent.IsDir() {
			path := filepath.Join(path, ent.Name())

			if path != partsDir {
				parseTemplateDir(f, path, partsDir, res)
			}
		} else {
			parseTemplateFile(f, filepath.Join(path, ent.Name()), partsDir, res)
		}
	}
}

func parseTemplateFile(f fs.FS, path string, partsDir string, res map[string]*template.Template) {
	t := template.Must(template.ParseFS(f, path, filepath.Join(partsDir, "*.html")))

	res[strings.TrimSuffix(path, ".html")] = t
}

var (
	_ Templates = FSTemplates{}
	_ Templates = CachedTemplates{}
)
