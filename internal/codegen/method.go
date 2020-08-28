package codegen

import (
	. "github.com/dave/jennifer/jen"
)

type MethodType string

const (
	REST_METHOD MethodType = "REST_METHOD"
	ACTION      MethodType = "ACTION"
	FINDER      MethodType = "FINDER"
)

type Method struct {
	MethodType    MethodType
	Name          string
	Doc           string
	Path          string
	OnEntity      bool
	EntityPathKey *PathKey
	PathKeys      []PathKey
	Params        []Field
	Return        *RestliType
	ReturnEntity  bool
}

type PathKey struct {
	Name string
	Type RestliType
}

func (m *Method) addEntityTypes(def *Group) {
	addEntityTypes(def, m.PathKeys)
}

func addEntityTypes(def *Group, pathKeys []PathKey) {
	for _, pk := range pathKeys {
		def.Id(pk.Name).Add(pk.Type.ReferencedType())
	}
}

func (m *Method) entityParams() (params []Code) {
	for _, p := range m.PathKeys {
		params = append(params, Id(p.Name))
	}
	return params
}

func (r *Resource) callFormatQueryUrl(def *Group) {
	def.List(UrlVar, Err()).
		Op(":=").
		Id(ClientReceiver).Dot("FormatQueryUrl").
		Call(Lit(r.RootResourceName), PathVar)
}

func (r *Resource) callEncodeQueryParams(def *Group) {
	def.List(UrlVar, Err()).
		Op("=").
		Add(QueryParams).Dot(EncodeQueryParams).Call()
}
