package tests

import (
	"encoding/json"
	"log"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
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
		testJsonEncoding(t, expected, new(Primitives), `{
  "primitiveInteger": 1,
  "primitiveLong": 23,
  "primitiveFloat": 52.5,
  "primitiveDouble": 66.5,
  "primitiveBytes": "@ABC🕴\u0001",
  "primitiveString": "a string,()'"
}`)
	})

	t.Run("ror2", func(t *testing.T) {
		testRor2Encoding(t, expected, new(Primitives),
			`(primitiveBytes:@ABC🕴,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a string%2C%28%29%27)`,
		)
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
				Fruits: conflictresolution.Fruits_ORANGE.Pointer(),
			},
		},
		UnionOfPrimitives: UnionOfPrimitives{
			PrimitivesUnion: UnionOfPrimitives_PrimitivesUnion{
				Int: &integer,
			},
		},
		AnotherUnionOfComplexTypes: UnionOfComplexTypes{
			ComplexTypeUnion: UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: conflictresolution.Fruits_APPLE.Pointer(),
			},
		},
		UnionOfSameTypes: UnionOfSameTypes{
			SameTypesUnion: UnionOfSameTypes_SameTypesUnion{
				Greeting: &hello,
			},
			UnionWithArrayMembers: UnionOfSameTypes_UnionWithArrayMembers{
				FruitArray: &[]conflictresolution.Fruits{
					conflictresolution.Fruits_ORANGE,
					conflictresolution.Fruits_APPLE,
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
		testJsonEncoding(t, expected, new(ComplexTypes), `{
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
  "anotherUnionOfComplexTypes": {
    "complexTypeUnion": {
      "testsuite.Fruits": "APPLE"
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
		testRor2Encoding(t, expected, new(ComplexTypes),
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

			t.Run("url", func(t *testing.T) {
				reader, err := restlicodec.NewRor2Reader(test.RestLi)
				require.NoError(t, err)
				require.NoError(t, test.Actual.UnmarshalRestLi(reader))
				requireEqual(t, &expected, &test.Actual)
			})
		})
	}
}

func TestReadIncorrectType(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewJsonReader([]byte(`{"message": [1]}`))
		require.Error(t, actual.UnmarshalRestLi(reader))
	})

	t.Run("url", func(t *testing.T) {
		var actual conflictresolution.Message
		reader, err := restlicodec.NewRor2Reader(`(message:List(1))`)
		require.NoError(t, err)
		require.Error(t, actual.UnmarshalRestLi(reader))
	})
}

func TestMissingRequiredFields(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var actual conflictresolution.Message
		reader := restlicodec.NewJsonReader([]byte(`{"id":1}`))
		require.Error(t, actual.UnmarshalRestLi(reader))
	})

	t.Run("url", func(t *testing.T) {
		var actual conflictresolution.Message
		reader, err := restlicodec.NewRor2Reader(`(id:1)`)
		require.NoError(t, err)
		require.Error(t, actual.UnmarshalRestLi(reader))
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

		s := restlicodec.NewStringPrimitiveUnmarshaler(new(string))
		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				reader, err := restlicodec.NewRor2Reader(test.Data)
				if err == nil {
					require.Error(t, s.UnmarshalRestLi(reader))
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

	testRor2Encoding(t, expected, new(Optionals), `(optionalArray:List(1,2))`)

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

	testRor2Encoding(t, expected, new(Optionals), `(optionalBytes:'',optionalString:'')`)
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
		reader, err := restlicodec.NewRor2Reader(`(map:(foo:1,bar:42),array:List(1,2),primitive:test)`)
		require.NoError(t, err)

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

type restliObject interface {
	restlicodec.Marshaler
	restlicodec.Unmarshaler
	Equals(interface{}) bool
}

func testJsonEncoding(t *testing.T, expected, actual restliObject, expectedRawJson string) {
	t.Run("encode", func(t *testing.T) {
		testJsonEquality(t, expected, expectedRawJson, nil, true)
	})

	t.Run("decode", func(t *testing.T) {
		decoder := restlicodec.NewJsonReader([]byte(expectedRawJson))
		require.NoError(t, actual.UnmarshalRestLi(decoder))
		requireEqual(t, expected, actual)
	})
}

func testJsonEquality(t *testing.T, obj restliObject, expectedRawJson string, excludedFields restlicodec.PathSpec, equal bool) {
	writer := restlicodec.NewCompactJsonWriterWithExcludedFields(excludedFields)
	require.NoError(t, obj.MarshalRestLi(writer))

	var expectedRaw map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(expectedRawJson), &expectedRaw))
	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(writer.Finalize()), &raw))
	if equal {
		require.Equal(t, expectedRaw, raw)
	} else {
		require.NotEqual(t, expectedRaw, raw)
	}
}

func testRor2Encoding(t *testing.T, expected, actual restliObject, expectedRawRor2 string) {
	t.Run("encode", func(t *testing.T) {
		testRor2Equality(t, expected, expectedRawRor2, nil, true)
	})

	t.Run("decode", func(t *testing.T) {
		reader, err := restlicodec.NewRor2Reader(expectedRawRor2)
		require.NoError(t, err)
		require.NoError(t, actual.UnmarshalRestLi(reader))
		requireEqual(t, expected, actual)
	})
}

func testRor2Equality(t *testing.T, obj restliObject, expectedRawRor2 string, excludedFields restlicodec.PathSpec, equal bool) {
	writer := restlicodec.NewRor2HeaderWriterWithExcludedFields(excludedFields)
	require.NoError(t, obj.MarshalRestLi(writer))
	log.Println(expectedRawRor2)

	unmarhsal := func(s string) interface{} {
		reader, err := restlicodec.NewRor2Reader(s)
		require.NoError(t, err)

		i, err := reader.ReadInterface()
		require.NoError(t, err)

		return i
	}

	expected := unmarhsal(expectedRawRor2)
	actual := unmarhsal(writer.Finalize())

	if equal {
		require.Equal(t, expected, actual)
	} else {
		require.NotEqual(t, expected, actual)
	}
}

func requireEqual(t *testing.T, expected, actual restliObject) {
	require.Equal(t, expected, actual)
	require.True(t, expected.Equals(actual))
}
