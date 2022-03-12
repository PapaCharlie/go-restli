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
	ClientReceiver      = "c"
	ClientType          = "client"
	ClientInterfaceType = "Client"

	ResourceInterfaceType = "Resource"

	CreatedEntity            = "CreatedEntity"
	CreatedAndReturnedEntity = "CreatedAndReturnedEntity"
	Elements                 = "Elements"
	BatchEntities            = "BatchEntities"
	BatchResponse            = "BatchResponse"
)

var (
	RestLiClientQual     = Code(Qual(utils.ProtocolPackage, RestLiClient))
	RestLiClientReceiver = Code(Id(ClientReceiver).Dot(RestLiClient))
	Context              = Code(Qual("context", "Context"))
	RequestContext       = Code(Op("*").Qual(utils.ProtocolPackage, "RequestContext"))
	RequestContextParam  = Code(Add(Ctx).Add(RequestContext))
	ElementsWithMetadata = Code(Qual(utils.ProtocolPackage, "ElementsWithMetadata"))

	Rp              = Code(Id("rp"))
	Ctx             = Code(Id("ctx"))
	Entity          = Code(Id("entity"))
	Entities        = Code(Id("entities"))
	Results         = Code(Id("results"))
	CreatedEntities = Code(Id("createdEntities"))
	Keys            = Code(Id("keys"))
	QueryParams     = Code(Id("queryParams"))
	ActionParams    = Code(Id("actionParams"))
	EmptyRecord     = Code(Qual(utils.StdTypesPackage, "EmptyRecord"))

	NoExcludedFields        = Code(Qual(utils.RestLiCodecPackage, "NoExcludedFields"))
	ReadOnlyFields          = Code(Id("ReadOnlyFields"))
	CreateOnlyFields        = Code(Id("CreateOnlyFields"))
	CreateAndReadOnlyFields = Code(Id("CreateAndReadOnlyFields"))

	BatchEntityUpdateResponse = Code(Op("*").Qual(utils.ProtocolPackage, "BatchEntityUpdateResponse"))
)
