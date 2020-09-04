package types

import . "github.com/dave/jennifer/jen"

const (
	MarshalRestLi              = "MarshalRestLi"
	UnmarshalRestLi            = "UnmarshalRestLi"
	EncodeQueryParams          = "EncodeQueryParams"
	PopulateLocalDefaultValues = "populateLocalDefaultValues"
	Equals                     = "Equals"
	ComputeHash                = "ComputeHash"
	ValidateUnionFields        = "ValidateUnionFields"
	ComplexKeyParamsField      = "Params"
	FinderNameParam            = "q"
	EntityIDsParam             = "ids"
	PartialUpdate              = "_PartialUpdate"

	RootPackage        = "github.com/PapaCharlie/go-restli"
	HashPackage        = RootPackage + "/fnv1a"
	ProtocolPackage    = RootPackage + "/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
)

var (
	NewJsonReader = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	NewRor2Reader = Code(Qual(RestLiCodecPackage, "NewRor2Reader"))

	NewHash = Code(Qual(HashPackage, "NewHash").Call())
	Hash    = Code(Qual(HashPackage, "Hash"))
)
