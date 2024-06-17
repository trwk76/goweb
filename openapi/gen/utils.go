package gen

import (
	"fmt"
	"strings"
)

func refKey(path string) string {
	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		return path[idx+1:]
	}

	return path
}

func uniqueName[T any](items map[string]T, base string) string {
	if _, fnd := items[base]; !fnd {
		return base
	}

	i := 1
	name := fmt.Sprintf("%s%d", base, i)

	for _, fnd := items[name]; fnd; {
		i++
		name = fmt.Sprintf("%s%d", base, i)
	}

	return name
}
