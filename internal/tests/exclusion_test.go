package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/stretchr/testify/require"
)

func TestPrimitiveFieldExclusion(t *testing.T) {
	expected := new(testsuite.PrimitiveField)
	spec := restlicodec.NewPathSpec("integer")

	t.Run("json", func(t *testing.T) {
		testJsonEquality(t, expected, `{
  "integer": 0
}`, nil, true)
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
  "arrayOfRecords": [
    {
      "foo": "foo"
    }
  ],
  "mapOfInts": {
    "two": 2
  },
  "mapOfRecords": {
    "record1": {},
    "record2": {
      "foo": "foo"
    }
  },
  "topLevelRecord": {
    "bar": "bar"
  },
  "topLevelUnion": {
    "extras.TopLevel": {
      "bar": "bar"
    }
  }
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

	roundTrip := func(name string, obj *conflictresolution.LargeRecord_PartialUpdate) {
		t.Run(name, func(t *testing.T) {
			require.Error(t, obj.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))
			w := restlicodec.NewCompactJsonWriter()
			require.NoError(t, obj.MarshalRestLi(w))

			r, err := restlicodec.NewJsonReaderWithExcludedFields([]byte(w.Finalize()), readOnlyFields, 1)
			require.NoError(t, err)
			require.Error(t, new(conflictresolution.LargeRecord_PartialUpdate).UnmarshalRestLi(r))
		})
	}

	roundTrip("Illegal delete", &conflictresolution.LargeRecord_PartialUpdate{
		Delete_Fields: conflictresolution.LargeRecord_PartialUpdate_Delete_Fields{
			OptionalArray: true,
		},
	})

	roundTrip("Illegal set", &conflictresolution.LargeRecord_PartialUpdate{
		Set_Fields: conflictresolution.LargeRecord_PartialUpdate_Set_Fields{
			OptionalArray: new([]int32),
		},
	})

	excluded := &conflictresolution.LargeRecord_PartialUpdate{
		Key: &conflictresolution.ComplexKey_PartialUpdate{},
	}
	require.NoError(t, excluded.MarshalRestLi(restlicodec.NewCompactJsonWriterWithExcludedFields(readOnlyFields)))

	roundTrip("Illegal update", &conflictresolution.LargeRecord_PartialUpdate{
		Key: &conflictresolution.ComplexKey_PartialUpdate{
			Set_Fields: conflictresolution.ComplexKey_PartialUpdate_Set_Fields{
				Part1: protocol.StringPointer(""),
			},
		},
	})
}
