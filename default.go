package goweb

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

func DefaultResponse(w http.ResponseWriter, r *http.Request, status int) {
	text := http.StatusText(status)

	w.Header().Set(HeaderContentType, ContentTypeText)
	w.Header().Set(HeaderContentLength, strconv.Itoa(len(text)))
	w.WriteHeader(status)

	if r.Method != http.MethodHead {
		io.Copy(w, strings.NewReader(text))
	}
}
