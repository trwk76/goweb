package gen

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

func genInit(dir string, pkgName string, a *API) {
	path := filepath.Join(dir, "init_gen.go")

	init := make([]string, 0, len(a.ops))

	for key := range a.ops {
		init = append(init, key)
	}

	sort.Slice(init, func(i, j int) bool { return init[i] < init[j] })

	for idx, opID := range init {
		init[idx] = fmt.Sprintf("\tm.HandleFunc(%s, %s)", strconv.Quote(a.ops[opID]), opID)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Errorf("error creating init file '%s': %s", path, err.Error()))
	}

	w := bufio.NewWriter(file)

	defer func() {
		w.Flush()
		file.Close()
	}()

	template.Must(template.ParseFS(templFS, "init.templ")).Execute(w, initData{
		PkgName:  pkgName,
		Title:    a.info.Title,
		BasePath: a.path,
		Init:     strings.Join(init, "\n"),
	})
}

type (
	initData struct {
		PkgName  string
		Title    string
		BasePath string
		Init     string
	}
)

//go:embed init.templ
var templFS embed.FS
