// Code generated by "stringer -type=Method -trimprefix Method_"; DO NOT EDIT.

package restli

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Method_Unknown-0]
	_ = x[Method_get-1]
	_ = x[Method_create-2]
	_ = x[Method_delete-3]
	_ = x[Method_update-4]
	_ = x[Method_partial_update-5]
	_ = x[Method_batch_get-6]
	_ = x[Method_batch_create-7]
	_ = x[Method_batch_delete-8]
	_ = x[Method_batch_update-9]
	_ = x[Method_batch_partial_update-10]
	_ = x[Method_get_all-11]
	_ = x[Method_action-12]
	_ = x[Method_finder-13]
}

const _RestLiMethod_name = "Unknowngetcreatedeleteupdatepartial_updatebatch_getbatch_createbatch_deletebatch_updatebatch_partial_updateget_allactionfinder"

var _RestLiMethod_index = [...]uint8{0, 7, 10, 16, 22, 28, 42, 51, 63, 75, 87, 107, 114, 120, 126}

func (i Method) String() string {
	if i < 0 || i >= Method(len(_RestLiMethod_index)-1) {
		return "Method(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RestLiMethod_name[_RestLiMethod_index[i]:_RestLiMethod_index[i+1]]
}
