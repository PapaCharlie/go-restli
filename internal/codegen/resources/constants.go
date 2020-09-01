package resources

import (
	"github.com/PapaCharlie/go-restli/internal/codegen/types"
	. "github.com/dave/jennifer/jen"
)

const (
	ResourcePath       = "ResourcePath"
	ResourceEntityPath = "ResourceEntityPath"

	WithContext = "WithContext"
	FindBy      = "FindBy"

	ClientReceiver      = "c"
	ClientType          = "client"
	ClientInterfaceType = "Client"
)

var (
	RestLiClient = Code(Qual(types.ProtocolPackage, "RestLiClient"))
	Context      = Code(Qual("context", "Context"))

	Url            = Code(Id("url"))
	Path           = Code(Id("path"))
	Ctx            = Code(Id("ctx"))
	Entity         = Code(Id("entity"))
	EntityKey      = Code(Id("entityKey"))
	ReturnedEntity = Code(Id("returnedEntity"))
	CreateParam    = Code(Id("create"))
	UpdateParam    = Code(Id("update"))
	QueryParams    = Code(Id("queryParams"))
	ActionParams   = Code(Id("actionParams"))
)
