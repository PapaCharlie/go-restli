package schema

import (
	"log"
	"strings"

	. "github.com/PapaCharlie/go-restli/codegen"
	"github.com/PapaCharlie/go-restli/codegen/models"
	. "github.com/dave/jennifer/jen"
)

const (
	ClientReceiver = "c"
	Req            = "req"
	Res            = "res"
	Url            = "url"
	ActionResult   = "actionResult"
	Client         = "Client"
	FormatQueryUrl = "FormatQueryUrl"
)

func (r *Resource) GenerateCode() (code []*CodeFile) {
	return generateResourceBindings(nil, r)
}

func (r *Resource) generateClient() (c *CodeFile) {
	c = NewCodeFile("client", r.PackagePath(), r.Name)
	log.Println(r.Namespace, r.Name)

	c.Code.Const().Id(ExportedIdentifier(r.Name + "Path")).Op("=").Lit(r.Path).Line()
	AddWordWrappedComment(c.Code, r.Doc).Line()
	c.Code.Type().Id(Client).Struct(Qual(ProtocolPackage, "RestLiClient")).Line().Line()

	return c
}

func generateResourceBindings(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	clientCodeFile := thisResource.generateClient()
	code = append(code, clientCodeFile)

	var newParentResources []*Resource
	newParentResources = append(newParentResources, parentResources...)
	newParentResources = append(newParentResources, thisResource)

	if thisResource.Simple != nil {
		code = append(code, thisResource.Simple.generateResourceBindings(parentResources, thisResource)...)
		code = append(code, thisResource.Simple.Entity.generateResourceBindings(parentResources, thisResource)...)
		for _, r := range thisResource.Simple.Entity.Subresources {
			code = append(code, generateResourceBindings(newParentResources, r)...)
		}
		return code
	}

	if thisResource.Collection != nil {
		code = append(code, thisResource.Collection.generateResourceBindings(clientCodeFile, parentResources, thisResource)...)
		code = append(code, thisResource.Collection.Entity.generateResourceBindings(parentResources, thisResource)...)
		for _, r := range thisResource.Collection.Entity.Subresources {
			code = append(code, generateResourceBindings(newParentResources, r)...)
		}
		return code
	}

	if thisResource.Association != nil {
		code = append(code, thisResource.Association.generateResourceBindings(parentResources, thisResource)...)
		code = append(code, thisResource.Association.Entity.generateResourceBindings(parentResources, thisResource)...)
		for _, r := range thisResource.Association.Entity.Subresources {
			code = append(code, generateResourceBindings(newParentResources, r)...)
		}
		return code
	}

	if thisResource.ActionsSet != nil {
		code = append(code, thisResource.ActionsSet.generateResourceBindings(parentResources, thisResource)...)
		return code
	}

	log.Panicln(thisResource, "does not define any resources")
	return
}

func (s *Simple) generateResourceBindings(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, action := range s.Actions {
		code = append(code, action.generateActionParamStructs(parentResources, thisResource, false))
	}
	return code
}

func (c *Collection) generateResourceBindings(clientCodeFile *CodeFile, parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, m := range c.Methods {
		if strings.ToLower(m.Method) == "get" {
			c.generateGet(clientCodeFile, parentResources, thisResource, m)
		}
	}
	for _, action := range c.Actions {
		code = append(code, action.generateActionParamStructs(parentResources, thisResource, false))
	}

	return code
}

func (a *Association) generateResourceBindings(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	if identifierCode := a.identifier(thisResource).Type.GenerateModelCode(); identifierCode != nil {
		code = append(code, identifierCode)
	}
	for _, action := range a.Actions {
		code = append(code, action.generateActionParamStructs(parentResources, thisResource, false))
	}
	return code
}

func (a *ActionsSet) generateResourceBindings(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, action := range a.Actions {
		code = append(code, action.generateActionParamStructs(parentResources, thisResource, false))
	}
	return code
}

func (e *Entity) generateResourceBindings(parentResources []*Resource, thisResource *Resource) (code []*CodeFile) {
	for _, action := range e.Actions {
		code = append(code, action.generateActionParamStructs(parentResources, thisResource, true))
	}
	return code
}

func (r *Resource) getIdentifier() *Identifier {
	if r.Simple != nil || r.ActionsSet != nil {
		return nil
	}

	if r.Collection != nil {
		return &r.Collection.Identifier
	}

	if r.Association != nil {
		return r.Association.identifier(r)
	}

	log.Panicln(r, "does not define any resources")
	return nil
}

func (r *Resource) getEntity() *Entity {
	if r.Simple != nil {
		return &r.Simple.Entity
	}

	if r.Collection != nil {
		return &r.Collection.Entity
	}

	if r.Association != nil {
		return &r.Association.Entity
	}

	return nil
}

func (a *Association) identifier(res *Resource) *Identifier {
	r := new(models.RecordModel)
	id := &Identifier{
		Name: a.Identifier,
		Type: &ResourceModel{models.Model{
			Ns:     models.Ns{Namespace: res.Namespace + "." + res.Name},
			Record: r,
		}},
	}

	r.Name = "AssocKey"
	id.Type.Name = "AssocKey"

	for i := range a.AssocKeys {
		k := a.AssocKeys[i]
		r.Fields = append(r.Fields, models.Field{
			NameAndDoc: models.NameAndDoc{Name: k.Name},
			Type:       &k.Type.Model,
		})
	}

	return id
}
