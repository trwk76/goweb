package openapi

import "log"

type (
	MediaTypes []MediaType

	MediaType interface {
		Name() string
		ContentType() string
	}
)

func NewMediaTypes(items ...MediaType) MediaTypes {
	res := make(MediaTypes, 0)

	for _, item := range items {
		res = res.Add(item)
	}

	return res
}

func (m MediaTypes) Name(name string) MediaType {
	for _, itm := range m {
		if itm.Name() == name {
			return itm
		}
	}

	return nil
}

func (m MediaTypes) ContentType(ctype string) MediaType {
	for _, itm := range m {
		if itm.ContentType() == ctype {
			return itm
		}
	}

	return nil
}

func (m MediaTypes) Add(item MediaType) MediaTypes {
	if m.Name(item.Name()) != nil {
		log.Fatalf("another media type is already registered under the name '%s'", item.Name())
	}

	if m.ContentType(item.ContentType()) != nil {
		log.Fatalf("another media type is already registered under the content type '%s'", item.ContentType())
	}

	return append(m, item)
}
