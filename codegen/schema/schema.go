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

type HasMethods struct {
	Methods []Method
}

type Identifier struct {
	Name string
	Type ResourceModel
}

type Simple struct {
	HasActions
	HasMethods
	Supports []string
	Entity   Entity
}

type Collection struct {
	HasActions
	HasMethods
	Identifier Identifier
	Supports   []string
	Finders    []Finder
	Entity     Entity
}

type AssocKey struct {
	Name string
	Type ResourceModel
}

type Association struct {
	HasActions
	HasMethods
	Identifier string
	AssocKeys  []AssocKey
	Supports   []string
	Entity     Entity
}

type Entity struct {
	HasActions
	Path         string
	Subresources []Resource
}

type Method struct {
	models.RecordModel
	Method          string
	PagingSupported bool
}

type Endpoint struct {
	models.RecordModel
	Returns *ResourceModel
}

type Finder struct {
	Endpoint
	PagingSupported bool
}

type Action struct {
	Endpoint
	ActionName string
	StructName   string
}
