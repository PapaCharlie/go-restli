package utils

import . "github.com/dave/jennifer/jen"

const (
	MarshalRestLi              = "MarshalRestLi"
	UnmarshalRestLi            = "UnmarshalRestLi"
	Marshal                    = "Marshal"
	Unmarshal                  = "Unmarshal"
	EncodeQueryParams          = "EncodeQueryParams"
	PopulateLocalDefaultValues = "populateLocalDefaultValues"
	Equals                     = "Equals"
	ComputeHash                = "ComputeHash"
	ValidateUnionFields        = "ValidateUnionFields"
	ComplexKeyParamsField      = "Params"
	ComplexKeyParams           = "$params"
	FinderNameParam            = "q"

	RootPackage         = "github.com/PapaCharlie/go-restli/v2"
	HashPackage         = RootPackage + "/fnv1a"
	RestLiPackage       = RootPackage + "/restli"
	RestLiPatchPackage  = RestLiPackage + "/patch"
	RestLiCodecPackage  = RootPackage + "/restlicodec"
	RestLiDataPackage   = RootPackage + "/restlidata"
	RestLiCommonPackage = RestLiDataPackage + "/generated/com/linkedin/restli/common"
	BatchKeySetPackage  = RestLiPackage + "/batchkeyset"
	EqualsPackage       = RestLiPackage + "/equals"
)

var (
	NewJsonReader     = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	NewRequiredFields = Code(Qual(RestLiCodecPackage, "NewRequiredFields"))

	WriteCustomTyperef = Code(Qual(RestLiCodecPackage, "WriteCustomTyperef"))
	ReadCustomTyperef  = Code(Qual(RestLiCodecPackage, "ReadCustomTyperef"))

	BatchKeySet = Code(Id("set"))

	Hash     = Code(Qual(HashPackage, "Hash"))
	NewHash  = Code(Qual(HashPackage, "NewHash").Call())
	ZeroHash = Code(Qual(HashPackage, "ZeroHash").Call())

	IllegalEnumConstant = Code(Qual(RestLiPackage, "IllegalEnumConstant"))
	UnknownEnumValue    = Code(Qual(RestLiPackage, "UnknownEnumValue"))

	MultiLineCall = Options{
		Open:      "(",
		Close:     ")",
		Separator: ",",
		Multi:     true,
	}
	MultiLineValues = Options{
		Open:      "{",
		Close:     "}",
		Separator: ",",
		Multi:     true,
	}
)
