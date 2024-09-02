package gen

import "fmt"

func uniqueName[T any](m map[string]T, base string) string {
	if _, fnd := m[base]; !fnd {
		return base
	}
	
	i := 1
	name := fmt.Sprintf("%s%d", base, i)
	
	for _, fnd := m[name]; fnd; {
		i++
		name = fmt.Sprintf("%s%d", base, i)
	}

	return name
}
