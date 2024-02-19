package web

type (
	SecurityRequirements []SecurityRequirement
	SecurityRequirement  map[string][]string

	Principals map[string]Principal

	Principal interface {
		HasRole(role string) bool
	}
)
