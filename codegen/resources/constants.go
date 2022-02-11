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

	RestLiClient        = "RestLiClient"
	ClientReceiver      = "c"
	ClientType          = "client"
	ClientInterfaceType = "Client"
)

var (
	RestLiClientQual     = Code(Qual(utils.ProtocolPackage, RestLiClient))
	RestLiClientReceiver = Code(Id(ClientReceiver).Dot(RestLiClient))
	Context              = Code(Qual("context", "Context"))

	Url             = Code(Id("url"))
	Path            = Code(Id("path"))
	Ctx             = Code(Id("ctx"))
	Key             = Code(Id("key"))
	Entity          = Code(Id("entity"))
	Entities        = Code(Id("entities"))
	CreatedEntities = Code(Id("createdEntities"))
	Statuses        = Code(Id("statuses"))
	EntityKey       = Code(Id("entityKey"))
	ReturnedEntity  = Code(Id("returnedEntity"))
	Keys            = Code(Id("keys"))
	QueryParams     = Code(Id("queryParams"))
	ActionParams    = Code(Id("actionParams"))

	NoExcludedFields        = Code(Qual(utils.RestLiCodecPackage, "NoExcludedFields"))
	ReadOnlyFields          = Code(Id("ReadOnlyFields"))
	CreateOnlyFields        = Code(Id("CreateOnlyFields"))
	CreateAndReadOnlyFields = Code(Id("CreateAndReadOnlyFields"))

	BatchEntities             = Code(Qual(utils.ProtocolPackage, "BatchEntities"))
	BatchEntityUpdateResponse = Code(Qual(utils.ProtocolPackage, "BatchEntityUpdateResponse"))
)
