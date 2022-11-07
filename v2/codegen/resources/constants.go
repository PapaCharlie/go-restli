package resources

import (
	"github.com/PapaCharlie/go-restli/v2/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const (
	ResourcePath       = "ResourcePath"
	ResourceEntityPath = "ResourceEntityPath"

	WithContext = "WithContext"
	FindBy      = "FindBy"

	RestLiClient        = "Client"
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
	RestLiClientQual          = Code(Qual(utils.RestLiPackage, RestLiClient))
	RestLiClientReceiver      = Code(Id(ClientReceiver).Dot(RestLiClient))
	Context                   = Code(Qual("context", "Context"))
	RequestContext            = Code(Op("*").Qual(utils.RestLiPackage, "RequestContext"))
	RequestContextParam       = Code(Add(Ctx).Add(RequestContext))
	ElementsWithMetadata      = Code(Qual(utils.RestLiCommonPackage, "ElementsWithMetadata"))
	BatchEntityUpdateResponse = Code(Op("*").Qual(utils.RestLiCommonPackage, "BatchEntityUpdateResponse"))
	EmptyRecord               = Code(Qual(utils.RestLiCommonPackage, "EmptyRecord"))

	Rp           = Code(Id("rp"))
	Ctx          = Code(Id("ctx"))
	Entity       = Code(Id("entity"))
	Entities     = Code(Id("entities"))
	Results      = Code(Id("results"))
	Keys         = Code(Id("keys"))
	QueryParams  = Code(Id("queryParams"))
	ActionParams = Code(Id("actionParams"))

	NoExcludedFields        = Code(Qual(utils.RestLiCodecPackage, "NoExcludedFields"))
	ReadOnlyFields          = Code(Id("ReadOnlyFields"))
	CreateAndReadOnlyFields = Code(Id("CreateAndReadOnlyFields"))
)
