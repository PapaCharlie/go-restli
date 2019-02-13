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
	NetHttp             = "net/http"
	Client              = "Client"
	HostnameClientField = "Hostname"
)

func (r *Resource) GenerateCode(packagePrefix string, sourceFilename string) *CodeFile {
	c := &CodeFile{
		SourceFilename: sourceFilename,
		PackagePath:    r.PackagePath(packagePrefix),
		Filename:       ExportedIdentifier(r.Name),
		Code:           Empty(),
	}

	// WIP
	r.generateClient(packagePrefix, c.Code)

	for _, s := range generateAllActionStructs(packagePrefix, nil, r) {
		c.Code.Add(s).Line()
	}

	return c
}

func (r *Resource) generateClient(packagePrefix string, code *Statement) {
	AddWordWrappedComment(code, r.Doc).Line()
	code.Type().Id(r.clientType()).Struct(
		Op("*").Qual(NetHttp, Client),
		Id(HostnameClientField).String(),
	)
	code.Line()
}

func (r *Resource) clientType() string {
	return ExportedIdentifier(r.Name) + Client
}

func (r *Resource) getIdentifier() *Identifier {
	if r.Simple != nil || r.ActionsSet != nil {
		return nil
	}

	if r.Collection != nil {
		return &r.Collection.Identifier
	}

	if r.Association != nil {
		str := models.String
		return &Identifier{
			Name: r.Association.Identifier,
			Type: ResourceModel{ Primitive: &str, }, }
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
