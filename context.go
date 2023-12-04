package goweb

import "github.com/google/uuid"

type (
	Context struct {
		id string
	}
)

func (c Context) ID() string {
	return c.id
}

func init() {
	uuid.EnableRandPool()
}
