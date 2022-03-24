package restlicodec

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryParamsReader map[string]*ror2QueryReader

type QueryParamsDecoder[T any] interface {
	NewInstance() T
	DecodeQueryParams(reader QueryParamsReader) error
}

func UnmarshalQueryParamsDecoder[T QueryParamsDecoder[T]](query string) (t T, err error) {
	reader, err := ParseQueryParams(query)
	if err != nil {
		return t, fmt.Errorf("go-restli: Illegal query params: %w", err)
	}

	t = t.NewInstance()
	return t, t.DecodeQueryParams(reader)
}

type ror2QueryReader struct{ *ror2Reader }

func (q *ror2QueryReader) ReadRecord(requiredFields RequiredFields, recordReader MapReader) error {
	return readRecord(q, requiredFields, recordReader)
}

func (q *ror2QueryReader) atInputStart() bool {
	return false
}

func (q QueryParamsReader) ReadRecord(requiredFields RequiredFields, recordReader MapReader) (err error) {
	tracker := missingFieldsTracker{}
	requiredFieldsRemaining := requiredFields.toMap()

	for k, v := range q {
		err = recordReader(v, k)
		if err != nil {
			return err
		}
		delete(requiredFieldsRemaining, k)

		tracker.missingFields = append(tracker.missingFields, v.missingFields...)
	}

	tracker.recordMissingRequiredFields(requiredFieldsRemaining)
	return tracker.checkMissingFields()
}

func ParseQueryParams(query string) (QueryParamsReader, error) {
	m := make(QueryParamsReader)
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if key == "" {
			continue
		}

		key, value, _ := strings.Cut(key, "=")

		err := validateRor2Input(value)
		if err != nil {
			return m, err
		}

		m[key] = &ror2QueryReader{&ror2Reader{
			missingFieldsTracker: missingFieldsTracker{
				currentScope: []deserializationScopeSegment{{segment: key}},
			},
			decoder: url.QueryUnescape,
			data:    []byte(value),
		}}
	}

	return m, nil
}
