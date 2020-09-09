package tests

import (
	"fmt"
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
		testEncodedJson(t, expected, `{"integer": 0}`, nil)
		testEncodedJson(t, expected, `{}`, spec)
	})

	t.Run("ror2", func(t *testing.T) {
		testEncodedRor2(t, expected, []string{`(integer:0)`}, nil)
		testEncodedRor2(t, expected, []string{`()`}, spec)
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
		testEncodedJson(t, expected, `{
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
}`, excludedFields)
	})

	t.Run("ror2", func(t *testing.T) {
		format := `(` +
			`arrayOfRecords:List((foo:foo)),` +
			`mapOfInts:(two:2),` +
			`mapOfRecords:%s,` +
			`topLevelRecord:(bar:bar),` +
			`topLevelUnion:(extras.TopLevel:(bar:bar))` +
			`)`
		for range make([]struct{}, 1000) {
			testEncodedRor2(t, expected, []string{
				fmt.Sprintf(format, "(record1:(),record2:(foo:foo))"),
				fmt.Sprintf(format, "(record2:(foo:foo),record1:())"),
			}, excludedFields)
		}
	})
}

func TestExcludeOnPartialUpdate(t *testing.T) {
	readOnlyFields := restlicodec.NewPathSpec(
		"optionalArray",
		"key/part1",
	)

	excluded := new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Delete.OptionalArray = true
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Update.OptionalArray = new([]int32)
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Key = new(conflictresolution.ComplexKey_PartialUpdate)
	require.NoError(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	excluded = new(conflictresolution.LargeRecord_PartialUpdate)
	excluded.Key = new(conflictresolution.ComplexKey_PartialUpdate)
	excluded.Key.Update.Part1 = new(string)
	require.Error(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))
}
