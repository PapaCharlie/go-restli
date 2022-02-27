package utils

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
	ComplexKeyParams           = "$params"
	FinderNameParam            = "q"
	IsUnknown                  = "IsUnknown"

	RootPackage        = "github.com/PapaCharlie/go-restli"
	HashPackage        = RootPackage + "/fnv1a"
	ProtocolPackage    = RootPackage + "/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
	BatchKeySetPackage = ProtocolPackage + "/batchkeyset"
	EqualsPackage      = ProtocolPackage + "/equals"
	StdTypesPackage    = ProtocolPackage + "/stdtypes"
)

var (
	NewJsonReader  = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	RequiredFields = Code(Qual(RestLiCodecPackage, "RequiredFields"))

	BatchKeySet = Code(Id("set"))

	Hash     = Code(Qual(HashPackage, "Hash"))
	NewHash  = Code(Qual(HashPackage, "NewHash").Call())
	ZeroHash = Code(Qual(HashPackage, "ZeroHash").Call())

	Enum                = Code(Qual(StdTypesPackage, "Enum"))
	IllegalEnumConstant = Code(Qual(StdTypesPackage, "IllegalEnumConstant"))
	UnknownEnumValue    = Code(Qual(StdTypesPackage, "UnknownEnumValue"))
)
