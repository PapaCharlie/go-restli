package types

const (
	EncodingJson  = "encoding/json"
	Unmarshal     = "Unmarshal"
	UnmarshalJSON = "UnmarshalJSON"
	Marshal       = "Marshal"
	MarshalJSON   = "MarshalJSON"

	Codec                = "codec"
	RestLiHeaderID       = "RestLiHeader_ID"
	MarshalRestLi        = "MarshalRestLi"
	UnmarshalRestLi      = "UnmarshalRestLi"
	EncodeQueryParams    = "EncodeQueryParams"
	RestLiDecode         = "RestLiDecode"
	RestLiCodec          = "RestLiCodec"
	RestLiUrlEncoder     = "RestLiQueryEncoder"
	RestLiUrlPathEncoder = "RestLiUrlPathEncoder"
	RestLiReducedEncoder = "RestLiReducedEncoder"

	PopulateLocalDefaultValues = "populateLocalDefaultValues"
	ValidateUnionFields        = "ValidateUnionFields"
	ComplexKeyParams           = "Params"

	PartialUpdate = "_PartialUpdate"

	NetHttp = "net/http"

	ProtocolPackage    = "github.com/PapaCharlie/go-restli/protocol"
	RestLiCodecPackage = ProtocolPackage + "/restlicodec"
)
