package schema

import (
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/protocol"
)

type Resource struct {
	models.Identifier
	Schema      *ResourceModel
	Path        string
	Doc         string
	Simple      *Simple
	Collection  *Collection
	Association *Association
	ActionsSet  *ActionsSet
}

type HasFinders struct {
	Finders []Finder
}

type HasActions struct {
	Actions []Action
}

type HasMethods struct {
	Methods []Method
}

type Identifier struct {
	Name string
	Type *ResourceModel
}

type ActionsSet struct {
	HasActions
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
	HasFinders
	Identifier Identifier
	Supports   []string
	Entity     Entity
}

type AssocKey struct {
	Name string
	Type *ResourceModel
}

type Association struct {
	HasActions
	HasMethods
	HasFinders
	Namespace  string
	Identifier string
	AssocKeys  []AssocKey
	Supports   []string
	Entity     Entity
}

type Entity struct {
	HasActions
	Path         string
	Subresources []*Resource
}

type Method struct {
	models.RecordModel
	Method          protocol.RestLiMethod
	PagingSupported bool
}

type Endpoint struct {
	models.RecordModel
	Returns *ResourceModel
}

type Finder struct {
	Endpoint
	FinderName      string
	StructName      string
	PagingSupported bool
}

type Action struct {
	Endpoint
	ActionName string
	StructName string
}
