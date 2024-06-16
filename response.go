package web

import (
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

	response interface {
		Write(w http.ResponseWriter)
	}
)

func NewTextResponse(status int, text string) Response {
	res := Response{Status: status, Body: strings.NewReader(text)}

	res.SetContentType(ContentTypeText)
	res.SetContentLength(len(text))

	return res
}

func NewDefaultResponse(status int) Response {
	return NewTextResponse(status, http.StatusText(status))
}

func (r *Response) SetContentType(value string) {
	r.ensureHeader()
	r.Header.Set(HeaderContentType, value)
}

func (r *Response) SetContentLength(value int) {
	r.ensureHeader()
	r.Header.Set(HeaderContentLength, strconv.Itoa(value))
}

func (r Response) Write(w http.ResponseWriter) {
	if len(r.Header) > 0 {
		for key, vals := range r.Header {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}
	}

	w.WriteHeader(r.Status)

	if r.Body != nil {
		if cl, ok := r.Body.(io.Closer); ok {
			defer cl.Close()
		}

		io.Copy(w, r.Body)
	}
}

func (r *Response) ensureHeader() {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
}

var _ response = Response{}
