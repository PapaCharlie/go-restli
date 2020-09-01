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
	PartialUpdate              = "_PartialUpdate"

	ProtocolPackage    = "github.com/PapaCharlie/go-restli/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
)

var (
	NewJsonReader = Code(Qual(RestLiCodecPackage, "NewJsonReader"))

	Other = Code(Id("other"))
)
