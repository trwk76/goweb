package web

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type (
	Response struct {
		Status int
		Header http.Header
		Body   io.Reader
	}

	ResponseBuffer struct {
		hdr    http.Header
		status int
		body   bytes.Buffer
	}

	response interface {
		WriteResponse(w http.ResponseWriter)
	}
)

func NewTextResponse(status int, text string) Response {
	res := Response{Status: status, Body: strings.NewReader(text)}

	SetHeaderContentType(res.Header, ContentTypeText)
	SetHeaderContentLength(res.Header, len(text))

	return res
}

func NewDefaultResponse(status int) Response {
	return NewTextResponse(status, http.StatusText(status))
}

func NewResponse(status int) Response {
	return Response{
		Status: status,
		Header: make(http.Header),
	}
}

func (r Response) WriteResponse(w http.ResponseWriter) {
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

func NewRespBuffer() ResponseBuffer {
	return ResponseBuffer{hdr: make(http.Header)}
}

func (r *ResponseBuffer) Header() http.Header {
	return r.hdr
}

func (r *ResponseBuffer) WriteHeader(status int) {
	r.status = status
}

func (r *ResponseBuffer) Write(p []byte) (int, error) {
	return r.body.Write(p)
}

func (r *ResponseBuffer) WriteResponse(w http.ResponseWriter) {
	for key, vals := range r.hdr {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(r.status)
	io.Copy(w, &r.body)
}

var (
	_ response            = Response{}
	_ http.ResponseWriter = (*ResponseBuffer)(nil)
)
