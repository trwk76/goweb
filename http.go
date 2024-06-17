package web

import (
	"fmt"
	"net/http"
	"strconv"
)

func ContentType(h http.Header) string {
	return h.Get(HeaderContentType)
}

func ContentLength(h http.Header) int {
	val, err := strconv.Atoi(h.Get(HeaderContentLength))
	if err != nil {
		val = 0
	}

	return val
}

func SetContentType(h http.Header, value string) {
	if value != "" {
		h.Set(HeaderContentType, value)
	} else {
		delete(h, HeaderContentType)
	}
}

func SetContentLength(h http.Header, value int) {
	if value >= 0 {
		h.Set(HeaderContentLength, strconv.Itoa(value))
	} else {
		delete(h, HeaderContentLength)
	}
}

func SetETag(h http.Header, value string) {
	if value >= "" {
		h.Set(HeaderETag, fmt.Sprintf(`"%s"`, value))
	} else {
		delete(h, HeaderETag)
	}
}

const (
	HeaderContentType   string = "Content-Type"
	HeaderContentLength string = "Content-Length"
	HeaderETag          string = "ETag"
	HeaderLocation      string = "Location"
)

const (
	ContentTypeHTML string = "text/html"
	ContentTypeJSON string = "application/json"
	ContentTypeText string = "text/plain"
	ContentTypeYAML string = "application/yaml"
)
