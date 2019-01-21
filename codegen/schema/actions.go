package schema

import (
	. "github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"strings"
)

func (r *Resource) generateAllActionStructs(packagePrefix string, parentResources ...*Resource) (defs []*Statement) {
	fullName := prefixNameWithParentResources(r.Name, parentResources...)
	defs = append(defs, Const().Id(fullName + "Path").Op("=").Lit(r.Path))

	parentResources = append(parentResources, r)

	defs = append(defs, r.Simple.generateActionParamStructs(packagePrefix, parentResources...)...)
	defs = append(defs, r.Collection.generateActionParamStructs(packagePrefix, parentResources...)...)
	defs = append(defs, r.Association.generateActionParamStructs(packagePrefix, parentResources...)...)
	defs = append(defs, r.ActionsSet.generateActionParamStructs(packagePrefix, parentResources...)...)

	defs = append(defs, r.Simple.Entity.generateActionParamStructs(packagePrefix, parentResources...)...)
	defs = append(defs, r.Collection.Entity.generateActionParamStructs(packagePrefix, parentResources...)...)
	defs = append(defs, r.Association.Entity.generateActionParamStructs(packagePrefix, parentResources...)...)

	for _, r := range r.Simple.Entity.Subresources {
		defs = append(defs, r.generateAllActionStructs(packagePrefix, parentResources...)...)
	}
	for _, r := range r.Collection.Entity.Subresources {
		defs = append(defs, r.generateAllActionStructs(packagePrefix, parentResources...)...)
	}
	for _, r := range r.Association.Entity.Subresources {
		defs = append(defs, r.generateAllActionStructs(packagePrefix, parentResources...)...)
	}

	return
}

func (h *HasActions) generateActionParamStructs(packagePrefix string, parentResources ...*Resource) (defs []*Statement) {
	for _, a := range h.Actions {
		defs = append(defs, a.generateActionParamStructs(packagePrefix, parentResources...))
	}
	return defs
}

func (a *Action) generateActionParamStructs(packagePrefix string, parentResources ...*Resource) (def *Statement) {
	fullName := prefixNameWithParentResources(a.Name, parentResources...)

	def = Empty()
	def.Const().Id(fullName + "Action").Op("=").Lit(a.Name).Line()

	var params []Code
	for _, p := range a.Parameters {
		paramDef := Empty()
		AddWordWrappedComment(paramDef, p.Doc).Line()
		paramDef.Id(ExportedIdentifier(p.Name))
		paramDef.Add(p.Type.GoType(packagePrefix)).Tag(JsonTag(p.Name))
		params = append(params, paramDef)
	}

	def.Type().Id(fullName + "ActionParams").Struct(params...)
	return
}

func prefixNameWithParentResources(name string, parentResources ...*Resource) string {
	var names []string
	for _, r := range parentResources {
		names = append(names, ExportedIdentifier(r.Name))
	}
	names = append(names, ExportedIdentifier(name))
	return strings.Join(names, "_")
}
