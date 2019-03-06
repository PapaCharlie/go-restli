package schema

import (
	. "github.com/dave/jennifer/jen"
	. "go-restli/codegen"
	"go-restli/codegen/models"
	"log"
)

const (
	ClientReceiver      = "c"
	Req                 = "req"
	Res                 = "res"
	Url                 = "url"
	ActionResult        = "actionResult"
	Client              = "Client"
	HostnameClientField = "Hostname"
)

func (r *Resource) GenerateCode(sourceFilename string) (code []*CodeFile) {
	code = append(code, r.generateClient())
	code = append(code, generateAllActionStructs(nil, r)...)
	for _, c := range code {
		c.SourceFilename = sourceFilename
	}

	return
}

func (r *Resource) generateClient() (c *CodeFile) {
	c = &CodeFile{
		PackagePath: r.PackagePath() + "/" + r.Name,
		Filename:    "client",
		Code:        Empty(),
	}

	Const().Id(r.Name + "Path").Op("=").Lit(r.Path).Line()
	AddWordWrappedComment(c.Code, r.Doc).Line()
	c.Code.Type().Id(Client).Struct(
		Op("*").Qual(NetHttp, Client),
		Id(HostnameClientField).String(),
	)

	return
}

func (r *Resource) getIdentifier() *Identifier {
	if r.Simple != nil || r.ActionsSet != nil {
		return nil
	}

	if r.Collection != nil {
		return &r.Collection.Identifier
	}

	if r.Association != nil {
		str := models.StringPrimitive
		return &Identifier{
			Name: r.Association.Identifier,
			Type: ResourceModel{models.Model{Primitive: &str,}},
		}
		//return &r.Association.Identifier
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

	log.Panicln("actionsSets do not have entities")
	return nil
}
