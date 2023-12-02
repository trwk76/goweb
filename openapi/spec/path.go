package spec

type (
	Paths map[string]PathItem

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
		Servers     Servers    `json:"servers,omitempty" yaml:"servers,omitempty"`
		Parameters  Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}

	Operation struct {
		OperationID  string               `json:"operationId,omitempty" yaml:"operationId,omitempty"`
		Summary      string               `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description  string               `json:"description,omitempty" yaml:"description,omitempty"`
		Tags         []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
		Deprecated   bool                 `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		ExternalDocs *ExternalDoc         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
		Security     SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
		Servers      Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
		Parameters   Parameters           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		RequestBody  *RequestBody         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
		Responses    Responses            `json:"responses,omitempty" yaml:"responses,omitempty"`
	}

	Parameters []Parameter

	Parameter struct {
		Ref             string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Name            string        `json:"name" yaml:"name"`
		In              ParameterIn   `json:"in" yaml:"in"`
		Description     string        `json:"description,omitempty" yaml:"description,omitempty"`
		Required        bool          `json:"required,omitempty" yaml:"required,omitempty"`
		Deprecated      bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		AllowEmptyValue bool          `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
		Style           string        `json:"style,omitempty" yaml:"style,omitempty"`
		Explode         bool          `json:"explode,omitempty" yaml:"explode,omitempty"`
		AllowReserved   bool          `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
		Schema          *Schema       `json:"schema,omitempty" yaml:"schema,omitempty"`
		Example         any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples        NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	ParameterIn string

	RequestBody struct {
		Ref         string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Content     MediaTypes `json:"content" yaml:"content"`
		Description string     `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool       `json:"required,omitempty" yaml:"required,omitempty"`
	}

	Responses map[string]Response

	Response struct {
		Ref         string       `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string       `json:"description" yaml:"description"`
		Headers     NamedHeaders `json:"headers,omitempty" yaml:"headers,omitempty"`
		Content     MediaTypes   `json:"content,omitempty" yaml:"content,omitempty"`
	}

	Header struct {
		Ref             string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description     string        `json:"description,omitempty" yaml:"description,omitempty"`
		Required        bool          `json:"required,omitempty" yaml:"required,omitempty"`
		Deprecated      bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		AllowEmptyValue bool          `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
		Style           string        `json:"style,omitempty" yaml:"style,omitempty"`
		Explode         bool          `json:"explode,omitempty" yaml:"explode,omitempty"`
		AllowReserved   bool          `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
		Schema          *Schema       `json:"schema,omitempty" yaml:"schema,omitempty"`
		Example         any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples        NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	MediaTypes map[string]MediaType

	MediaType struct {
		Schema   Schema        `json:"schema" yaml:"schema"`
		Example  any           `json:"example,omitempty" yaml:"example,omitempty"`
		Examples NamedExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	}
)

const (
	ParameterInCookie ParameterIn = "cookie"
	ParameterInHeader ParameterIn = "header"
	ParameterInPath   ParameterIn = "path"
	ParameterInQuery  ParameterIn = "query"
)
