package doc

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	web "github.com/trwk76/goweb"
)

func Init(p web.Path, apis ...API) {
	web.GET_HEAD(p, web.Redirect(http.StatusMovedPermanently, p.Path()))

	ents, err := docFS.ReadDir(".")
	if err != nil {
		panic(fmt.Errorf("error reading swaggerui embedded fs: %s", err.Error()))
	}

	for _, ent := range ents {
		if ent.Name() == "swagger-initializer.js" {
			raw, err := fs.ReadFile(docFS, ent.Name())
			if err != nil {
				panic(err)
			}

			urlsJSON, err := json.Marshal(apis)
			if err != nil {
				panic(fmt.Errorf("error marshaling api specifications: %s", err.Error()))
			}

			raw = []byte(strings.ReplaceAll(string(raw), "url: \"https://petstore.swagger.io/v2/swagger.json\"", "urls: "+string(urlsJSON)))

			web.GET_HEAD(p, web.Content(ent.Name(), web.ContentTypeJSON, Version, time.Now(), bytes.NewReader(raw)))
		} else {
			web.GET_HEAD(p, web.File(docFS, ent.Name(), Version))
		}
	}
}

type (
	API struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url"`
	}
)

const Version string = "5.17.14"

//go:embed *.css *.html *.js *.map
var docFS embed.FS
