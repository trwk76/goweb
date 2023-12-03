package openapi

type (
	MediaTypes []MediaType

	MediaType struct {
		n  string
		ct string
	}
)

func (t MediaType) Name() string {
	return t.n
}

func (t MediaType) ContentType() string {
	return t.ct
}
