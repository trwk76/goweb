package web

import (
	"fmt"
	"net/http"
	"strconv"
)

func SetHeaderContentType(h http.Header, value string) {
	if value != "" {
		h.Set(HeaderContentType, value)
	} else {
		delete(h, HeaderContentType)
	}
}

func SetHeaderContentLength(h http.Header, value int) {
	if value >= 0 {
		h.Set(HeaderContentLength, strconv.Itoa(value))
	} else {
		delete(h, HeaderContentLength)
	}
}

func SetHeaderETag(h http.Header, value string) {
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
)

const (
	ContentTypeHTML string = "text/html"
	ContentTypeJSOM string = "application/json"
	ContentTypeText string = "text/plain"
	ContentTypeYAML string = "application/yaml"
)
