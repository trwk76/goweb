package openapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/trwk76/goweb"
)

func WriteJSON(w http.ResponseWriter, status int, body any) {
	raw, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	w.Header().Set(goweb.HeaderContentType, goweb.ContentTypeJSON)
	w.Header().Set(goweb.HeaderContentLength, strconv.Itoa(len(raw)))
	w.WriteHeader(status)

	io.Copy(w, bytes.NewReader(raw))
}
