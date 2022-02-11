package protocol

import "github.com/PapaCharlie/go-restli/protocol/restlicodec"

const (
	elementsField = "elements"
	valueField    = "value"
	statusField   = "status"
	resultsField  = "results"
	errorsField   = "errors"
	idField       = "id"
	entityField   = "entity"
)

var (
	elementsRequiredResponseFields              = restlicodec.RequiredFields{elementsField}
	actionRequiredResponseFields                = restlicodec.RequiredFields{valueField}
	batchEntityUpdateResponseRequiredFields     = restlicodec.RequiredFields{statusField}
	batchRequestResponseRequiredFields          = restlicodec.RequiredFields{resultsField}
	batchCreateResponseRequiredFields           = restlicodec.RequiredFields{statusField}
	batchCreateWithReturnResponseRequiredFields = restlicodec.RequiredFields{statusField, entityField}
)
