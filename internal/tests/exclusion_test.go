package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/stretchr/testify/require"
)

func TestPrimitiveFieldExclusion(t *testing.T) {
	expected := new(testsuite.PrimitiveField)
	spec := restlicodec.NewPathSpec("integer")

	t.Run("json", func(t *testing.T) {
		testJsonEquality(t, expected, `{"integer": 0}`, nil, true)
		testJsonEquality(t, expected, `{}`, spec, true)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Equality(t, expected, `(integer:0)`, nil, true)
		testRor2Equality(t, expected, `()`, spec, true)
	})
}

func TestComplexFieldExclusion(t *testing.T) {
	record := &extras.TopLevel{
		Foo: "foo",
		Bar: "bar",
	}
	expected := &extras.EvenMoreComplexTypes{
		MapOfInts:      map[string]int32{"one": 1, "two": 2},
		TopLevelRecord: *record,
		ArrayOfRecords: []*extras.TopLevel{record},
		MapOfRecords: map[string]*extras.TopLevel{
			"record1": record,
			"record2": record,
		},
		TopLevelUnion: extras.TopLevelUnion{TopLevel: record},
	}
	excludedFields := restlicodec.NewPathSpec(
		"mapOfInts/one",
		"topLevelRecord/foo",
		"arrayOfRecords/*/bar",
		"mapOfRecords/*/bar",
		"mapOfRecords/record1/foo",
		"topLevelUnion/extras.TopLevel/foo",
	)

	t.Run("json", func(t *testing.T) {
		testJsonEquality(t, expected, `{
  "mapOfInts": {
    "two": 2
  },
  "topLevelRecord": {
    "bar": "bar"
  },
  "arrayOfRecords": [
    {"foo": "foo"}
  ],
  "mapOfRecords": {
    "record1": {},
    "record2": {"foo": "foo"}
  },
  "topLevelUnion": {"extras.TopLevel": {"bar": "bar"}}
}`, excludedFields, true)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Equality(t, expected, `(`+
			`arrayOfRecords:List((foo:foo)),`+
			`mapOfInts:(two:2),`+
			`mapOfRecords:(record1:(),record2:(foo:foo)),`+
			`topLevelRecord:(bar:bar),`+
			`topLevelUnion:(extras.TopLevel:(bar:bar))`+
			`)`, excludedFields, true)
	})
}

func TestExcludeOnPartialUpdate(t *testing.T) {
	readOnlyFields := restlicodec.NewPathSpec(
		"optionalArray",
		"key/part1",
	)

	excluded := new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Delete_Fields.OptionalArray = true
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Update_Fields.OptionalArray = new([]int32)
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Key = new(conflictresolution.ComplexKey_PartialUpdate)
	require.NoError(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Key = new(conflictresolution.ComplexKey_PartialUpdate)
	excluded.Key.Update_Fields.Part1 = new(string)
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))
}
