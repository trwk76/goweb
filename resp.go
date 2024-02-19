package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type (
	Response struct {
		Status int
		Header http.Header
		Body   io.Reader
	}
)

func DefaultResponse(status int) Response {
	txt := http.StatusText(status)
	return TextResponse(status, txt)
}

func TextResponse(status int, val string) Response {
	return FixedResponse(status, TextContentType, len(val), strings.NewReader(val))
}

func JSONResponse(status int, val any) Response {
	raw, _ := json.Marshal(val)
	return FixedResponse(status, JSONContentType, len(raw), bytes.NewReader(raw))
}

func Redirect(url string) Response {
	return Response{
		Status: http.StatusMovedPermanently,
		Header: http.Header{
			"Location": {url},
		},
	}
}

func FixedResponse(status int, ctype string, clen int, body io.Reader) Response {
	return Response{
		Status: status,
		Header: http.Header{
			"Content-Type":   {ctype},
			"Content-Length": {strconv.Itoa(clen)},
		},
		Body: body,
	}
}

const (
	JSONContentType string = "application/json"
	TextContentType string = "text/plain"
	YAMLContentType string = "application/yaml"
)

func (r Response) write(w http.ResponseWriter) error {
	hdr := w.Header()

	if r.Header != nil {
		for key, vals := range r.Header {
			for _, val := range vals {
				hdr.Add(key, val)
			}
		}
	}

	w.WriteHeader(r.Status)

	if r.Body != nil {
		_, err := io.Copy(w, r.Body)
		return err
	}

	return nil
}
