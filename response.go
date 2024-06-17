package web

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type (
	ResponseData struct {
		StatusCode int
		Header     http.Header
		Body       io.Reader
	}

	ResponseBuffer struct {
		hdr    http.Header
		status int
		body   bytes.Buffer
	}

	Response interface {
		Status() int
		Size() int
		WriteResponse(w http.ResponseWriter)
	}
)

func NewTextResponse(status int, text string) ResponseData {
	res := ResponseData{StatusCode: status, Body: strings.NewReader(text)}

	SetContentType(res.Header, ContentTypeText)
	SetContentLength(res.Header, len(text))

	return res
}

func NewDefaultResponse(status int) ResponseData {
	return NewTextResponse(status, http.StatusText(status))
}

func NewResponseData(status int) ResponseData {
	return ResponseData{
		StatusCode: status,
		Header:     make(http.Header),
	}
}

func (r ResponseData) Status() int {
	return r.StatusCode
}

func (r ResponseData) Size() int {
	return ContentLength(r.Header)
}

func (r ResponseData) WriteResponse(w http.ResponseWriter) {
	if len(r.Header) > 0 {
		for key, vals := range r.Header {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}
	}

	w.WriteHeader(r.StatusCode)

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

func (r *ResponseBuffer) Status() int {
	return r.status
}

func (r *ResponseBuffer) Size() int {
	return ContentLength(r.hdr)
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
	_ Response            = ResponseData{}
	_ Response            = (*ResponseBuffer)(nil)
	_ http.ResponseWriter = (*ResponseBuffer)(nil)
)
