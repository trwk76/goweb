package spec

type (
	NamedPaths map[string]PathItem

	PathItem struct {
		Summary     string     `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string     `json:"description,omitempty" yaml:"description,omitempty"`
		GET         *Operation `json:"get,omitempty" yaml:"get,omitempty"`
		PUT         *Operation `json:"put,omitempty" yaml:"put,omitempty"`
		POST        *Operation `json:"post,omitempty" yaml:"post,omitempty"`
		DELETE      *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
		OPTIONS     *Operation `json:"options,omitempty" yaml:"options,omitempty"`
		HEAD        *Operation `json:"head,omitempty" yaml:"head,omitempty"`
		PATCH       *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
		TRACE       *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
		Parameters  Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}

	Operation struct {
		Responses   Responses            `json:"responses" yaml:"responses"`
		OperationID string               `json:"operationId,omitempty" yaml:"operationId,omitempty"`
		Summary     string               `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string               `json:"description,omitempty" yaml:"description,omitempty"`
		Deprecated  bool                 `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		Parameters  Parameters           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		RequestBody *RequestBody         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
		Security    SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
		Tags        []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
	}

	NamedResponses map[string]Response
	Responses      map[string]Response

	Response struct {
		Ref         string       `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string       `json:"description,omitempty" yaml:"description,omitempty"`
		Headers     NamedHeaders `json:"headers,omitempty" yaml:"headers,omitempty"`
		Content     MediaTypes   `json:"content,omitempty" yaml:"content,omitempty"`
	}

	NamedRequestBodies map[string]RequestBody

	RequestBody struct {
		Ref         string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string     `json:"description,omitempty" yaml:"description,omitempty"`
		Content     MediaTypes `json:"content,omitempty" yaml:"content,omitempty"`
		Required    bool       `json:"required,omitempty" yaml:"required,omitempty"`
	}

	MediaTypes map[string]MediaType

	MediaType struct {
		Schema   Schema        `json:"schema,omitempty" yaml:"schema,omitempty"`
		Example  any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	NamedParameters map[string]Parameter
	Parameters      []Parameter

	Parameter struct {
		Ref         string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
		In          ParameterIn   `json:"in,omitempty" yaml:"in,omitempty"`
		Description string        `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool          `json:"required,omitempty" yaml:"required,omitempty"`
		Deprecated  bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		Schema      Schema        `json:"schema,omitempty" yaml:"schema,omitempty"`
		Example     any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples    NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	ParameterIn string

	NamedHeaders map[string]Header

	Header struct {
		Ref         string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string        `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool          `json:"required,omitempty" yaml:"required,omitempty"`
		Deprecated  bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		Schema      Schema        `json:"schema,omitempty" yaml:"schema,omitempty"`
		Example     any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples    NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	NamedExamples map[string]Example

	Example struct {
		Ref         string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Summary     string `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
		Value       any    `json:"value,omitempty" yaml:"value,omitempty"`
	}
)

const (
	ParameterInCookie ParameterIn = "cookie"
	ParameterInHeader ParameterIn = "header"
	ParameterInPath   ParameterIn = "path"
	ParameterInQuery  ParameterIn = "query"
)
