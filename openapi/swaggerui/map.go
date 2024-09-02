package swaggerui

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/trwk76/goweb"
)

func Map(m *http.ServeMux, path string, apis ...API) {
	if !strings.HasSuffix(path, "/") {
		m.Handle("GET "+path, http.RedirectHandler(path+"/", http.StatusFound))
		path += "/"
	}

	ents, err := distFS.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("error reading embedded swaggerui directory: %s", err.Error()))
	}

	for _, ent := range ents {
		if ent.Name() == "index.html" {
			m.HandleFunc("GET "+path, contentHandler(ent.Name()))
		} else if ent.Name() == "swagger-initializer.js" {
			urlsJSON, err := json.Marshal(apis)
			if err != nil {
				panic(fmt.Errorf("error marshaling api info to json: %s", err.Error()))
			}

			raw, err := distFS.ReadFile(ent.Name())
			if err != nil {
				panic(fmt.Errorf("error reading embedded swaggerui '%s' file: %s", ent.Name(), err.Error()))
			}

			raw = bytes.ReplaceAll(raw, []byte(`url: "https://petstore.swagger.io/v2/swagger.json"`), append([]byte("urls: "), urlsJSON...))

			m.HandleFunc("GET "+path+ent.Name(), func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(goweb.HeaderContentType, mime.TypeByExtension(filepath.Ext(ent.Name())))
				w.Header().Set(goweb.HeaderContentLength, strconv.Itoa(len(raw)))
				w.WriteHeader(http.StatusOK)

				if r.Method != http.MethodHead {
					io.Copy(w, bytes.NewReader(raw))
				}
			})
		} else {
			m.HandleFunc("GET "+path+ent.Name(), contentHandler(ent.Name()))
		}
	}
}

type (
	API struct {
		URL  string `json:"url"`
		Name string `json:"name,omitempty"`
	}
)

//go:embed *.css *.html *.js *.map *.png
var distFS embed.FS

func contentHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check precondition
		for _, etag := range r.Header.Values(goweb.HeaderIfNoneMatch) {
			if etag == Version {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		f, err := distFS.Open(name)
		if err != nil {
			panic(fmt.Errorf("error reading embedded swaggerui '%s' file: %s", name, err.Error()))
		}

		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			panic(fmt.Errorf("error getting embedded swaggerui '%s' file info: %s", name, err.Error()))
		}

		w.Header().Set(goweb.HeaderContentType, mime.TypeByExtension(filepath.Ext(name)))
		w.Header().Set(goweb.HeaderContentLength, strconv.Itoa(int(info.Size())))
		w.Header().Set(goweb.HeaderETag, Version)
		w.WriteHeader(http.StatusOK)

		if r.Method != http.MethodHead {
			io.Copy(w, f)
		}
	}
}

const Version string = `"5.17.14"`
