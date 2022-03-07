package tests

import (
	"encoding/json"
	"math"
	"net/url"
	"testing"
	"unicode/utf8"

	"github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
	"github.com/stretchr/testify/require"
)

func TestEncodePrimitives(t *testing.T) {
	expected := &Primitives{
		PrimitiveInteger: 1,
		PrimitiveLong:    23,
		PrimitiveFloat:   52.5,
		PrimitiveDouble:  66.5,
		PrimitiveBytes:   []byte("@ABC🕴" + string([]byte{1})),
		PrimitiveString:  `a string,()'`,
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, UnmarshalRestLiPrimitives, `{
  "primitiveBytes": "@ABC🕴\u0001",
  "primitiveDouble": 66.5,
  "primitiveFloat": 52.5,
  "primitiveInteger": 1,
  "primitiveLong": 23,
  "primitiveString": "a string,()'"
}`)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Encoding(t, expected, UnmarshalRestLiPrimitives,
			`(primitiveBytes:@ABC🕴,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a string%2C%28%29%27)`,
		)
	})
}

func TestEncodeInfinity(t *testing.T) {
	inf := math.Inf(-1)
	expected := &Optionals{
		OptionalDouble: &inf,
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, UnmarshalRestLiOptionals, `{
  "optionalDouble": "-Infinity"
}`)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Encoding(t, expected, UnmarshalRestLiOptionals, `(optionalDouble:-Infinity)`)
	})
}

func TestEncodeComplexTypes(t *testing.T) {
	integer := int32(5)
	hello := "Hello"
	expected := &ComplexTypes{
		ArrayOfMaps: ArrayOfMaps{
			ArrayOfMaps: []map[string]int32{
				{
					"one": 1,
				},
				{
					"two": 2,
				},
			},
		},
		MapOfInts: MapOfInts{
			MapOfInts: map[string]int32{
				"one": 1,
			},
		},
		RecordWithProps: RecordWithProps{
			Integer: &integer,
		},
		UnionOfComplexTypes: UnionOfComplexTypes{
			ComplexTypeUnion: UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: testsuite.Fruits_ORANGE.Pointer(),
			},
		},
		UnionOfPrimitives: UnionOfPrimitives{
			PrimitivesUnion: UnionOfPrimitives_PrimitivesUnion{
				Int: &integer,
			},
		},
		AnotherUnionOfComplexTypes: UnionOfComplexTypes{
			ComplexTypeUnion: UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: testsuite.Fruits_APPLE.Pointer(),
			},
		},
		UnionOfSameTypes: UnionOfSameTypes{
			SameTypesUnion: UnionOfSameTypes_SameTypesUnion{
				Greeting: &hello,
			},
			UnionWithArrayMembers: UnionOfSameTypes_UnionWithArrayMembers{
				FruitArray: &[]testsuite.Fruits{
					testsuite.Fruits_ORANGE,
					testsuite.Fruits_APPLE,
				},
			},
			UnionWithMapMembers: UnionOfSameTypes_UnionWithMapMembers{
				IntMap: &map[string]int32{
					"one": 1,
				},
			},
		},
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, UnmarshalRestLiComplexTypes, `{
  "anotherUnionOfComplexTypes": {
    "complexTypeUnion": {
      "testsuite.Fruits": "APPLE"
    }
  },
  "arrayOfMaps": {
    "arrayOfMaps": [
      {
        "one": 1
      },
      {
        "two": 2
      }
    ]
  },
  "mapOfInts": {
    "mapOfInts": {
      "one": 1
    }
  },
  "recordWithProps": {
    "integer": 5
  },
  "unionOfComplexTypes": {
    "complexTypeUnion": {
      "testsuite.Fruits": "ORANGE"
    }
  },
  "unionOfPrimitives": {
    "primitivesUnion": {
      "int": 5
    }
  },
  "unionOfSameTypes": {
    "sameTypesUnion": {
      "greeting": "Hello"
    },
    "unionWithArrayMembers": {
      "fruitArray": [
        "ORANGE",
        "APPLE"
      ]
    },
    "unionWithMapMembers": {
      "intMap": {
        "one": 1
      }
    }
  }
}`)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Encoding(t, expected, UnmarshalRestLiComplexTypes,
			`(anotherUnionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:APPLE)),`+
				`arrayOfMaps:(arrayOfMaps:List((one:1),(two:2))),`+
				`mapOfInts:(mapOfInts:(one:1)),`+
				`recordWithProps:(integer:5),`+
				`unionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:ORANGE)),`+
				`unionOfPrimitives:(primitivesUnion:(int:5)),`+
				`unionOfSameTypes:(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1))))`,
		)
	})
}

func TestUnknownFieldReads(t *testing.T) {
	id := int64(1)
	expected := conflictresolution.Message{
		Id:      &id,
		Message: "test",
	}

	tests := []struct {
		Name   string
		Json   string
		RestLi string
		Actual conflictresolution.Message
	}{
		{
			Name:   "Extra primitive field before",
			Json:   `{"foo":false,"id":1,"message":"test"}`,
			RestLi: `(foo:false,id:1,message:test)`,
		},
		{
			Name:   "Extra primitive field in the middle",
			Json:   `{"id":1,"foo":false,"message":"test"}`,
			RestLi: `(id:1,foo:false,message:test)`,
		},
		{
			Name:   "Extra primitive field at the end",
			Json:   `{"id":1,"message":"test","foo":false}`,
			RestLi: `(id:1,message:test,foo:false)`,
		},
		{
			Name:   "Extra object field before",
			Json:   `{"foo":{"bar": 1},"id":1,"message":"test"}`,
			RestLi: `(foo:(bar:1),id:1,message:test)`,
		},
		{
			Name:   "Extra object field in the middle",
			Json:   `{"id":1,"foo":{"bar":1},"message":"test"}`,
			RestLi: `(id:1,foo:(bar:1),message:test)`,
		},
		{
			Name:   "Extra object field at the end",
			Json:   `{"id":1,"message":"test","foo":{"bar":1}}`,
			RestLi: `(id:1,message:test,foo:(bar:1))`,
		},
		{
			Name:   "Extra array field before",
			Json:   `{"foo":[42],"id":1,"message":"test"}`,
			RestLi: `(foo:List(42),id:1,message:test)`,
		},
		{
			Name:   "Extra array field in the middle",
			Json:   `{"id":1,"foo":[42],"message":"test"}`,
			RestLi: `(id:1,foo:List(42),message:test)`,
		},
		{
			Name:   "Extra array field at the end",
			Json:   `{"id":1,"message":"test","foo":[42]}`,
			RestLi: `(id:1,message:test,foo:List(42))`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				reader := restlicodec.NewJsonReader([]byte(test.Json))
				require.NoError(t, test.Actual.UnmarshalRestLi(reader))
				requireEqual(t, &expected, &test.Actual)
			})

			t.Run("ror2", func(t *testing.T) {
				reader := newRor2Reader(t, test.RestLi)
				require.NoError(t, test.Actual.UnmarshalRestLi(reader))
				requireEqual(t, &expected, &test.Actual)
			})
		})
	}
}

func TestRawRecord(t *testing.T) {
	expected := &stdtypes.RawRecord{
		"arrayOfInts": []int{1, 2, 3},
		"arrayOfMaps": []map[string]int{
			{"foo": 1},
			{"bar": 2},
		},
		"mapOfInts": map[string]int{"foo": 1, "bar": 2},
		"mapOfArrays": map[string][]string{
			"one": {"a"},
			"two": {"b", "c"},
		},
		"record":    &extras.SinglePrimitiveField{String: "abc"},
		"int":       42,
		"fixed":     [5]byte{'a', 'b', 'c', 'd', 'e'},
		"realFixed": [5]byte{'f', 'g', 'h', 'i', 'j'},
		"bytes":     []byte{'a', 'b'},
	}

	t.Run("json", func(t *testing.T) {
		expectedJson := `{
  "arrayOfInts": [
    1,
    2,
    3
  ],
  "arrayOfMaps": [
    {
      "foo": 1
    },
    {
      "bar": 2
    }
  ],
  "bytes": "ab",
  "fixed": "abcde",
  "int": 42,
  "mapOfArrays": {
    "one": [
      "a"
    ],
    "two": [
      "b",
      "c"
    ]
  },
  "mapOfInts": {
    "bar": 2,
    "foo": 1
  },
  "realFixed": "fghij",
  "record": {
    "string": "abc"
  }
}`
		w := restlicodec.NewPrettyJsonWriter()
		require.NoError(t, expected.MarshalRestLi(w))
		require.Equal(t, expectedJson, w.Finalize())

		raw := new(stdtypes.RawRecord)
		require.NoError(t, raw.UnmarshalRestLi(restlicodec.NewJsonReader([]byte(expectedJson))))

		w = restlicodec.NewPrettyJsonWriter()
		require.NoError(t, raw.MarshalRestLi(w))
		require.Equal(t, expectedJson, w.Finalize())
	})

	t.Run("ror2", func(t *testing.T) {
		expectedRor2 := `(` +
			`arrayOfInts:List(1,2,3),` +
			`arrayOfMaps:List((foo:1),(bar:2)),` +
			`bytes:ab,` +
			`fixed:abcde,` +
			`int:42,` +
			`mapOfArrays:(one:List(a),two:List(b,c)),` +
			`mapOfInts:(bar:2,foo:1),` +
			`realFixed:fghij,` +
			`record:(string:abc)` +
			`)`

		w := restlicodec.NewRor2HeaderWriter()
		require.NoError(t, expected.MarshalRestLi(w))
		require.Equal(t, expectedRor2, w.Finalize())

		raw := new(stdtypes.RawRecord)
		r, err := restlicodec.NewRor2Reader(expectedRor2)
		require.NoError(t, err)
		require.NoError(t, raw.UnmarshalRestLi(r))

		w = restlicodec.NewRor2HeaderWriter()
		require.NoError(t, raw.MarshalRestLi(w))
		require.Equal(t, expectedRor2, w.Finalize())
	})

}

func TestRawRecordUnmarshalTo(t *testing.T) {
	raw := &stdtypes.RawRecord{
		"string": "abc",
	}
	expected := &extras.SinglePrimitiveField{String: "abc"}
	actual := new(extras.SinglePrimitiveField)
	require.NoError(t, raw.UnmarshalTo(actual))
	require.Equal(t, expected, actual)
	require.True(t, expected.Equals(actual))
}

func TestDeserializationErrorHandling(t *testing.T) {
	checkDeserializationError := func(t *testing.T, err error, scope string) {
		require.Error(t, err)
		require.Equal(t, scope, err.(*restlicodec.DeserializationError).Scope)
	}

	t.Run("array", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			var actual conflictresolution.Message
			reader := restlicodec.NewJsonReader([]byte(`{"message": [1]}`))
			checkDeserializationError(t, actual.UnmarshalRestLi(reader), "message")
		})

		t.Run("ror2", func(t *testing.T) {
			var actual conflictresolution.Message
			reader := newRor2Reader(t, `(message:List(1))`)
			checkDeserializationError(t, actual.UnmarshalRestLi(reader), "message")
		})
	})

	t.Run("illegal primitive", func(t *testing.T) {
		t.Run("json", func(t *testing.T) {
			var actual Optionals
			reader := restlicodec.NewJsonReader([]byte(`{"optionalInteger": "asd"}`))
			checkDeserializationError(t, actual.UnmarshalRestLi(reader), "optionalInteger")
		})

		t.Run("ror2", func(t *testing.T) {
			var actual Optionals
			reader := newRor2Reader(t, `(optionalInteger:asd)`)
			checkDeserializationError(t, actual.UnmarshalRestLi(reader), "optionalInteger")
		})
	})
}

func TestMissingRequiredFields(t *testing.T) {
	checkRequiredFieldsError := func(t *testing.T, err error, fields ...string) {
		require.Equal(t, &restlicodec.MissingRequiredFieldsError{Fields: fields}, err)
	}

	t.Run("simple", func(t *testing.T) {
		expected := conflictresolution.Message{Id: new(int64)}
		*expected.Id = 1

		t.Run("json", func(t *testing.T) {
			var actual conflictresolution.Message
			reader := restlicodec.NewJsonReader([]byte(`{"id":1}`))
			checkRequiredFieldsError(t, actual.UnmarshalRestLi(reader), "message")
			require.Equal(t, expected, actual)
		})

		t.Run("ror2", func(t *testing.T) {
			var actual conflictresolution.Message
			reader := newRor2Reader(t, `(id:1)`)
			checkRequiredFieldsError(t, actual.UnmarshalRestLi(reader), "message")
			require.Equal(t, expected, actual)
		})
	})

	t.Run("complex", func(t *testing.T) {
		expected := extras.RecordArray{
			Records: []*extras.TopLevel{
				{
					Foo: "",
					Bar: "bar",
				},
				{
					Foo: "foo",
					Bar: "bar",
				},
			},
		}

		t.Run("json", func(t *testing.T) {
			var actual extras.RecordArray
			reader := restlicodec.NewJsonReader([]byte(
				`{ "records": [ { "bar": "bar" }, { "foo": "foo", "bar": "bar" } ] }`,
			))
			checkRequiredFieldsError(t, actual.UnmarshalRestLi(reader), "records[0].foo")
			require.Equal(t, expected, actual)
		})

		t.Run("ror2", func(t *testing.T) {
			var actual extras.RecordArray
			reader := newRor2Reader(t, `(records:List((bar:bar),(foo:foo,bar:bar)))`)
			checkRequiredFieldsError(t, actual.UnmarshalRestLi(reader), "records[0].foo")
			require.Equal(t, expected, actual)
		})
	})
}

func TestIllegalRor2Strings(t *testing.T) {
	t.Run("primitives", func(t *testing.T) {
		tests := []struct {
			Name string
			Data string
		}{
			{
				Name: "(",
				Data: `(`,
			},
			{
				Name: ",",
				Data: `,`,
			},
			{
				Name: ")",
				Data: `)`,
			},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				reader, err := restlicodec.NewRor2Reader(test.Data)
				if err == nil {
					_, err = reader.ReadString()
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("complex", func(t *testing.T) {
		tests := []struct {
			Name string
			Data string
		}{
			{
				Name: "Unbalanced end parens",
				Data: `(message:`,
			},
			{
				Name: "Unbalanced start parens",
				Data: `message:)`,
			},
			{
				Name: "Unescaped ','",
				Data: `(message:,)`,
			},
			{
				Name: "Garbage",
				Data: `(message:foo,bar)`,
			},
			{
				Name: "Too many parens",
				Data: `(message:))`,
			},
		}

		var m conflictresolution.Message
		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				reader, err := restlicodec.NewRor2Reader(test.Data)
				if err == nil {
					require.Error(t, m.UnmarshalRestLi(reader))
				}
			})
		}
	})
}

func TestMapEncoding(t *testing.T) {
	expected := &Optionals{
		OptionalMap: &map[string]int32{
			"one": 1,
			"two": 2,
		},
	}

	t.Run("multipleElements", func(t *testing.T) {
		writer := restlicodec.NewRor2HeaderWriter()
		require.NoError(t, expected.MarshalRestLi(writer))

		serialized := writer.Finalize()
		if serialized != `(optionalMap:(one:1,two:2))` && serialized != `(optionalMap:(two:2,one:1))` {
			t.Fail()
		}
	})

	expected = &Optionals{OptionalMap: &map[string]int32{}}
	writer := restlicodec.NewRor2HeaderWriter()
	require.NoError(t, expected.MarshalRestLi(writer))
	require.Equal(t, `(optionalMap:())`, writer.Finalize())
}

func TestArrayEncoding(t *testing.T) {
	expected := &Optionals{OptionalArray: &[]int32{1, 2}}

	testRor2Encoding(t, expected, UnmarshalRestLiOptionals, `(optionalArray:List(1,2))`)

	expected = &Optionals{OptionalArray: &[]int32{}}
	writer := restlicodec.NewRor2HeaderWriter()
	require.NoError(t, expected.MarshalRestLi(writer))
	require.Equal(t, `(optionalArray:List())`, writer.Finalize())
}

func TestEmptyStringAndBytes(t *testing.T) {
	expected := &Optionals{
		OptionalBytes:  new([]byte),
		OptionalString: new(string),
	}

	testRor2Encoding(t, expected, UnmarshalRestLiOptionals, `(optionalBytes:'',optionalString:'')`)
}

func TestRaw(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		reader := restlicodec.NewJsonReader([]byte(`{
		  "map": { "foo": 1, "bar": 42 },
		  "array": [1,2],
		  "primitive": "test"
		}`))

		require.NoError(t, reader.ReadMap(func(reader restlicodec.Reader, field string) error {
			switch field {
			case "map":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				var actual map[string]int
				require.NoError(t, json.Unmarshal(raw, &actual))
				require.Equal(t, map[string]int{"foo": 1, "bar": 42}, actual)
			case "array":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				var actual []int
				require.NoError(t, json.Unmarshal(raw, &actual))
				require.Equal(t, []int{1, 2}, actual)
			case "primitive":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				var actual string
				require.NoError(t, json.Unmarshal(raw, &actual))
				require.Equal(t, "test", actual)
			}
			return nil
		}))
	})
	t.Run("ror2", func(t *testing.T) {
		reader := newRor2Reader(t, `(map:(foo:1,bar:42),array:List(1,2),primitive:test)`)

		require.NoError(t, reader.ReadMap(func(reader restlicodec.Reader, field string) error {
			switch field {
			case "map":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				require.Equal(t, "(foo:1,bar:42)", string(raw))
			case "array":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				require.Equal(t, "List(1,2)", string(raw))
			case "primitive":
				raw, err := reader.ReadRawBytes()
				require.NoError(t, err)
				require.Equal(t, "test", string(raw))
			}
			return nil
		}))
	})
}

func testJsonEncoding[T protocol.RestLiObject[T]](t *testing.T, expected T, unmarshaler restlicodec.GenericUnmarshaler[T], expectedRawJson string) {
	t.Run("encode", func(t *testing.T) {
		testJsonEquality(t, expected, expectedRawJson, nil, true)
	})

	t.Run("decode", func(t *testing.T) {
		actual, err := unmarshaler(restlicodec.NewJsonReader([]byte(expectedRawJson)))
		require.NoError(t, err)
		requireEqual(t, expected, actual)
	})
}

func testJsonEquality[T protocol.RestLiObject[T]](t *testing.T, obj T, expectedRawJson string, excludedFields restlicodec.PathSpec, equal bool) {
	writer := restlicodec.NewPrettyJsonWriterWithExcludedFields(excludedFields)
	require.NoError(t, obj.MarshalRestLi(writer))

	actualRawJson := writer.Finalize()
	if equal {
		require.Equal(t, expectedRawJson, actualRawJson)
	} else {
		require.NotEqual(t, expectedRawJson, actualRawJson)
	}
}

func testRor2Encoding[T protocol.RestLiObject[T]](t *testing.T, expected T, unmarshaler restlicodec.GenericUnmarshaler[T], expectedRawRor2 string) {
	t.Run("encode", func(t *testing.T) {
		testRor2Equality(t, expected, expectedRawRor2, nil, true)
	})

	t.Run("decode", func(t *testing.T) {
		actual, err := unmarshaler(newRor2Reader(t, expectedRawRor2))
		require.NoError(t, err)
		requireEqual(t, expected, actual)
	})
}

func testRor2Equality[T protocol.RestLiObject[T]](t *testing.T, obj T, expectedRawRor2 string, excludedFields restlicodec.PathSpec, equal bool) {
	writer := restlicodec.NewRor2HeaderWriterWithExcludedFields(excludedFields)
	require.NoError(t, obj.MarshalRestLi(writer))

	actualRawRor2 := writer.Finalize()

	if equal {
		require.Equal(t, expectedRawRor2, actualRawRor2)
	} else {
		require.NotEqual(t, expectedRawRor2, actualRawRor2)
	}
}

func requireEqual[T protocol.RestLiObject[T]](t *testing.T, expected, actual T) {
	require.Equal(t, expected, actual)
	require.True(t, expected.Equals(actual), "%+v and %+v should be equals", expected, actual)
}

func newRor2Reader(t *testing.T, data string) restlicodec.Reader {
	reader, err := restlicodec.NewRor2Reader(data)
	require.NoError(t, err)
	return reader
}

var pathEscapeSpecialCases = map[rune]string{
	'!': "!",
	'*': "*",
	':': "%3A",
}

func TestPathEscape(t *testing.T) {
	// Enumerate all valid UTF-8 runes
	for c := rune(0); c < 0x110000; c++ {
		// Skip surrogate half runes
		if !utf8.ValidRune(c) {
			continue
		}
		// Ensure that each rune is escaped correctly according to the normal URL path encoding spec, except for the
		// known special cases
		if expected, ok := pathEscapeSpecialCases[c]; ok {
			require.Equal(t, expected, restlicodec.Ror2PathEscape(string(c)))
		} else {
			require.Equal(t, url.PathEscape(string(c)), restlicodec.Ror2PathEscape(string(c)))
		}
	}
}

func TestReadArray(t *testing.T) {
	r := restlicodec.NewJsonReader([]byte(`["APPLE"]`))
	fruits, err := restlicodec.ReadArray[testsuite.Fruits](r, testsuite.UnmarshalRestLiFruits)
	require.NoError(t, err)
	require.Equal(t, []testsuite.Fruits{testsuite.Fruits_APPLE}, fruits)
}
