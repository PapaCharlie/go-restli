package schema

import "go-restli/codegen/models"

type Resource struct {
	models.Ns
	Name        string
	Path        string
	Schema      string
	Doc         string
	Simple      *Simple
	Collection  *Collection
	Association *Association
	ActionsSet  *HasActions
}

type HasActions struct {
	Actions []Action
}

type Identifier struct {
	Name string
	Type ResourceModel
}

type Simple struct {
	HasActions
	Supports []string
	Methods  []Method
	Entity   Entity
}

type Collection struct {
	HasActions
	Identifier Identifier
	Supports   []string
	Methods    []Method
	Finders    []Finder
	Entity     Entity
}

type AssocKey struct {
	Name string
	Type ResourceModel
}

type Association struct {
	HasActions
	Identifier string
	AssocKeys []AssocKey
	Supports []string
	Methods  []Method
	Entity   Entity
}

type Entity struct {
	HasActions
	Path         string
	Subresources []Resource
}

type Method struct {
	Method     string
	Doc        string
	Parameters []Parameter
}

type Endpoint struct {
	Name       string
	Doc        string
	Parameters []Parameter
	Returns    *ResourceModel
}

type Parameter struct {
	Name     string
	Doc      string
	Type     ResourceModel
	Optional bool
	Default  *string
}

type Finder struct {
	Endpoint
	PagingSupported bool
}

type Action struct {
	Endpoint
}
