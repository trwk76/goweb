package spec

type (
	OpenAPI struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Info    *Info  `json:"info" yaml:"info"`
	}

	Info struct {
		Title       string `json:"title" yaml:"title"`
		Version     string `json:"version" yaml:"version"`
		Summary     string `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}
)
