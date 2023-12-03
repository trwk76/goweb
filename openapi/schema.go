package openapi

import (
	"reflect"

	"github.com/trwk76/goweb/openapi/spec"
)

type (
	Schemas struct {
		n []*Schema
	}

	Schema struct {
		n string
		t reflect.Type
		m *MediaType
		s spec.Schema
	}
)

func (s Schemas) Name(name string) *Schema {
	for _, itm := range s.n {
		if itm.n == name {
			return itm
		}
	}

	return nil
}

func (s Schemas) Type(t reflect.Type, m *MediaType) *Schema {
	for _, itm := range s.n {
		if itm.t == t && itm.m == m {
			return itm
		}
	}

	return nil
}

func (s Schema) Type() reflect.Type {
	return s.t
}

func (s Schema) MediaType() *MediaType {
	return s.m
}

func (s Schema) Spec() spec.Schema {
	return s.s
}
