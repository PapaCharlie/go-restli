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

	RootPackage        = "github.com/PapaCharlie/go-restli"
	HashPackage        = RootPackage + "/fnv1a"
	RestLiPackage      = RootPackage + "/restli"
	RestLiCodecPackage = RootPackage + "/restlicodec"
	RestLiDataPackage  = RootPackage + "/restlidata"
	BatchKeySetPackage = RestLiPackage + "/batchkeyset"
	EqualsPackage      = RestLiPackage + "/equals"
)

var (
	NewJsonReader  = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	RequiredFields = Code(Qual(RestLiCodecPackage, "RequiredFields"))

	BatchKeySet = Code(Id("set"))

	Hash     = Code(Qual(HashPackage, "Hash"))
	NewHash  = Code(Qual(HashPackage, "NewHash").Call())
	ZeroHash = Code(Qual(HashPackage, "ZeroHash").Call())

	IllegalEnumConstant = Code(Qual(RestLiPackage, "IllegalEnumConstant"))
	UnknownEnumValue    = Code(Qual(RestLiPackage, "UnknownEnumValue"))
)
