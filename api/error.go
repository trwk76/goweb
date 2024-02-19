package api

type (
	ParseError struct {
		Items []ParseErrorItem
	}

	ParseErrorItem struct {
		Path    string
		Message string
	}
)

func (e *ParseError) Add(path string, message string) {
	e.Items = append(e.Items, ParseErrorItem{
		Path:    path,
		Message: message,
	})
}

func (e *ParseError) AddRequired(path string) {
	e.Add(path, "requires a value but none provided")
}

func (e ParseError) Error() string {
	return "one or more errors occurred while parsing request"
}
