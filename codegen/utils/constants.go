package utils

import . "github.com/dave/jennifer/jen"

const (
	MarshalRestLi              = "MarshalRestLi"
	UnmarshalRestLi            = "UnmarshalRestLi"
	EncodeQueryParams          = "EncodeQueryParams"
	PopulateLocalDefaultValues = "populateLocalDefaultValues"
	Equals                     = "Equals"
	EqualsInterface            = "EqualsInterface"
	ComputeHash                = "ComputeHash"
	ValidateUnionFields        = "ValidateUnionFields"
	ComplexKeyParamsField      = "Params"
	ComplexKeyParams           = "$params"
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

	EntityIDsEncoder      = Code(Id("entityIDsEncoder"))
	BatchEntityIDsEncoder = Code(Qual(ProtocolPackage, "BatchEntityIDsEncoder"))

	NewHash = Code(Qual(HashPackage, "NewHash").Call())
	Hash    = Code(Qual(HashPackage, "Hash"))

	IllegalEnumConstant = Code(Qual(ProtocolPackage, "IllegalEnumConstant"))
	UnknownEnumValue    = Code(Qual(ProtocolPackage, "UnknownEnumValue"))
)
