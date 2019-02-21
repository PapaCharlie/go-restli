package schema

var RestliMethodNameMapping = map[string]string{
	"get":            "Get",
	"create":         "Create",
	"delete":         "Delete",
	"update":         "Update",
	"partial_update": "PartialUpdate",

	"batch_get":            "BatchGet",
	"batch_create":         "BatchCreate",
	"batch_delete":         "BatchDelete",
	"batch_update":         "BatchUpdate",
	"batch_partial_update": "BatchPartialUpdate",

	"get_all": "GetAll",
}

var RestliMethodToHttpMethod = map[string]string{
	"get":            "GET",
	"create":         "POST",
	"delete":         "DELETE",
	"update":         "PUT",
	"partial_update": "POST",

	"batch_get":            "GET",
	"batch_create":         "POST",
	"batch_delete":         "DELETE",
	"batch_update":         "PUT",
	"batch_partial_update": "POST",

	"get_all": "GET",
}
