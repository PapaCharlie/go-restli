package restlidata

import (
	"github.com/PapaCharlie/go-restli/restlicodec"
)

const (
	ElementsField = "elements"
	ValueField    = "value"
	StatusField   = "status"
	StatusesField = "statuses"
	ResultsField  = "results"
	ErrorField    = "error"
	ErrorsField   = "errors"
	IdField       = "id"
	LocationField = "location"
	PagingField   = "paging"
	MetadataField = "metadata"
	EntityField   = "entity"
	EntitiesField = "entities"
)

var (
	elementsRequiredResponseFields              = restlicodec.RequiredFields{ElementsField}
	batchEntityUpdateResponseRequiredFields     = restlicodec.RequiredFields{StatusField}
	batchResponseRequiredFields                 = restlicodec.RequiredFields{ResultsField}
	batchCreateResponseRequiredFields           = restlicodec.RequiredFields{StatusField}
	batchCreateWithReturnResponseRequiredFields = restlicodec.RequiredFields{StatusField, EntityField}
)
