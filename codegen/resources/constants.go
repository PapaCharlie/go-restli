package resources

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const (
	ResourcePath       = "resourcePath"
	ResourceEntityPath = "resourceEntityPath"

	WithContext = "WithContext"
	FindBy      = "FindBy"

	RestLiClient        = "RestLiClient"
	SimpleClient        = "SimpleClient"
	CollectionClient    = "CollectionClient"
	ClientReceiver      = "c"
	ClientType          = "client"
	ClientInterfaceType = "Client"
)

var (
	RestLiClientQual     = Code(Qual(utils.ProtocolPackage, RestLiClient))
	SimpleClientQual     = Code(Qual(utils.ProtocolPackage, SimpleClient))
	CollectionClientQual = Code(Qual(utils.ProtocolPackage, CollectionClient))
	RestLiClientReceiver = Code(Id(ClientReceiver).Dot(RestLiClient))
	Context              = Code(Qual("context", "Context"))

	Rp              = Code(Id("rp"))
	Ctx             = Code(Id("ctx"))
	Entity          = Code(Id("entity"))
	Entities        = Code(Id("entities"))
	CreatedEntity   = Code(Id("createdEntity"))
	CreatedEntities = Code(Id("createdEntities"))
	Statuses        = Code(Id("statuses"))
	Keys            = Code(Id("keys"))
	QueryParams     = Code(Id("queryParams"))
	ActionParams    = Code(Id("actionParams"))

	NoExcludedFields        = Code(Qual(utils.RestLiCodecPackage, "NoExcludedFields"))
	ReadOnlyFields          = Code(Id("ReadOnlyFields"))
	CreateOnlyFields        = Code(Id("CreateOnlyFields"))
	CreateAndReadOnlyFields = Code(Id("CreateAndReadOnlyFields"))

	BatchEntityUpdateResponse = Code(Qual(utils.ProtocolPackage, "BatchEntityUpdateResponse"))
)
