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
	IsUnknown                  = "IsUnknown"

	RootPackage        = "github.com/PapaCharlie/go-restli"
	HashPackage        = RootPackage + "/fnv1a"
	ProtocolPackage    = RootPackage + "/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
	BatchKeySetPackage = ProtocolPackage + "/batchkeyset"
	EqualsPackage      = ProtocolPackage + "/equals"
	StdStructsPackage  = ProtocolPackage + "/stdstructs"
)

var (
	NewJsonReader = Code(Qual(RestLiCodecPackage, "NewJsonReader"))
	NewRor2Reader = Code(Qual(RestLiCodecPackage, "NewRor2Reader"))

	BatchKeySet           = Code(Id("set"))
	EntityIDsEncoder      = Code(Id("entityIDsEncoder"))
	BatchEntityIDsEncoder = Code(Qual(ProtocolPackage, "BatchEntityIDsEncoder"))

	Hash             = Code(Qual(HashPackage, "Hash"))
	NewHash          = Code(Qual(HashPackage, "NewHash").Call())
	ZeroHash         = Code(Qual(HashPackage, "ZeroHash").Call())
	AddArray         = Code(Qual(HashPackage, "AddArray"))
	AddHashableArray = Code(Qual(HashPackage, "AddHashableArray"))
	AddMap           = Code(Qual(HashPackage, "AddMap"))
	AddHashableMap   = Code(Qual(HashPackage, "AddHashableMap"))

	IllegalEnumConstant = Code(Qual(ProtocolPackage, "IllegalEnumConstant"))
	UnknownEnumValue    = Code(Qual(ProtocolPackage, "UnknownEnumValue"))

	Enum = Code(Qual(ProtocolPackage, "Enum"))

	RequiredFields  = Code(Qual(RestLiCodecPackage, "RequiredFields"))
	ReadMap         = Code(Qual(RestLiCodecPackage, "ReadMap"))
	ReadArray       = Code(Qual(RestLiCodecPackage, "ReadArray"))
	UnmarshalerFunc = Code(Qual(RestLiCodecPackage, "UnmarshalerFunc"))

	WriteMap   = Code(Qual(RestLiCodecPackage, "WriteMap"))
	WriteArray = Code(Qual(RestLiCodecPackage, "WriteArray"))
)
