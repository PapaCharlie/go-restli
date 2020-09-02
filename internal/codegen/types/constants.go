package types

import . "github.com/dave/jennifer/jen"

const (
	MarshalRestLi              = "MarshalRestLi"
	UnmarshalRestLi            = "UnmarshalRestLi"
	EncodeQueryParams          = "EncodeQueryParams"
	PopulateLocalDefaultValues = "populateLocalDefaultValues"
	Equals                     = "Equals"
	ValidateUnionFields        = "ValidateUnionFields"
	ComplexKeyParamsField      = "Params"
	FinderNameParam            = "q"
	EntityIDsParam             = "ids"
	PartialUpdate              = "_PartialUpdate"

	ProtocolPackage    = "github.com/PapaCharlie/go-restli/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
)

var (
	NewJsonReader = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	NewRor2Reader = Code(Qual(RestLiCodecPackage, "NewRor2Reader"))

	Other = Code(Id("other"))
)
