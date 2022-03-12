package protocol

import (
	"github.com/PapaCharlie/go-restli/protocol/batchkeyset"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

const (
	elementsField = "elements"
	valueField    = "value"
	statusField   = "status"
	statusesField = "statuses"
	resultsField  = "results"
	errorField    = "error"
	errorsField   = "errors"
	idField       = "id"
	locationField = "location"
	pagingField   = "paging"
	metadataField = "metadata"
	entityField   = "entity"
	entitiesField = "entities"
)

var (
	elementsRequiredResponseFields              = restlicodec.RequiredFields{elementsField}
	entitiesRequiredResponseFields              = restlicodec.RequiredFields{entitiesField}
	entityIdsRequiredResponseFields             = restlicodec.RequiredFields{batchkeyset.EntityIDsField}
	actionRequiredResponseFields                = restlicodec.RequiredFields{valueField}
	batchEntityUpdateResponseRequiredFields     = restlicodec.RequiredFields{statusField}
	batchResponseRequiredFields                 = restlicodec.RequiredFields{resultsField}
	batchCreateResponseRequiredFields           = restlicodec.RequiredFields{statusField}
	batchCreateWithReturnResponseRequiredFields = restlicodec.RequiredFields{statusField, entityField}
)
