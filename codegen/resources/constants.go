package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
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
	RestLiClient = Code(Qual(utils.ProtocolPackage, "RestLiClient"))
	Context      = Code(Qual("context", "Context"))

	Url            = Code(Id("url"))
	Path           = Code(Id("path"))
	Ctx            = Code(Id("ctx"))
	Key            = Code(Id("key"))
	Entity         = Code(Id("entity"))
	Entities       = Code(Id("entities"))
	Statuses       = Code(Id("statuses"))
	EntityKey      = Code(Id("entityKey"))
	ReturnedEntity = Code(Id("returnedEntity"))
	CreateParam    = Code(Id("create"))
	UpdateParam    = Code(Id("update"))
	Keys           = Code(Id("keys"))
	QueryParams    = Code(Id("queryParams"))
	ActionParams   = Code(Id("actionParams"))

	NoExcludedFields        = Code(Qual(utils.RestLiCodecPackage, "NoExcludedFields"))
	ReadOnlyFields          = Code(Id("ReadOnlyFields"))
	CreateOnlyFields        = Code(Id("CreateOnlyFields"))
	CreateAndReadOnlyFields = Code(Id("CreateAndReadOnlyFields"))

	BatchEntities             = Code(Qual(utils.ProtocolPackage, "BatchEntities"))
	BatchEntityUpdateResponse = Code(Qual(utils.ProtocolPackage, "BatchEntityUpdateResponse"))
)
