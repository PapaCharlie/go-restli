package tests

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/stretchr/testify/require"
)

func TestEncodePrimitives(t *testing.T) {
	expected := &testsuite.Primitives{
		PrimitiveInteger: 1,
		PrimitiveLong:    23,
		PrimitiveFloat:   52.5,
		PrimitiveDouble:  66.5,
		PrimitiveBytes:   []byte("@ABC🕴" + string([]byte{1})),
		PrimitiveString:  `a string,()'`,
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, `{
  "primitiveInteger": 1,
  "primitiveLong": 23,
  "primitiveFloat": 52.5,
  "primitiveDouble": 66.5,
  "primitiveBytes": "@ABC🕴\u0001",
  "primitiveString": "a string,()'"
}`)
	})

	t.Run("urlEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(primitiveBytes:%40ABC%F0%9F%95%B4%01,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a+string%2C%28%29%27)`, false)
	})

	t.Run("reducedEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(primitiveBytes:@ABC🕴,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a string%2C%28%29%27)`, true)
	})
}

func TestEncodeComplexTypes(t *testing.T) {
	orange := conflictresolution.Fruits_ORANGE
	apple := conflictresolution.Fruits_APPLE
	integer := int32(5)
	hello := "Hello"
	expected := &testsuite.ComplexTypes{
		ArrayOfMaps: testsuite.ArrayOfMaps{
			ArrayOfMaps: []map[string]int32{
				{
					"one": 1,
				},
				{
					"two": 2,
				},
			},
		},
		MapOfInts: testsuite.MapOfInts{
			MapOfInts: map[string]int32{
				"one": 1,
			},
		},
		RecordWithProps: testsuite.RecordWithProps{
			Integer: &integer,
		},
		UnionOfComplexTypes: testsuite.UnionOfComplexTypes{
			ComplexTypeUnion: testsuite.UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: &orange,
			},
		},
		UnionOfPrimitives: testsuite.UnionOfPrimitives{
			PrimitivesUnion: testsuite.UnionOfPrimitives_PrimitivesUnion{
				Int: &integer,
			},
		},
		AnotherUnionOfComplexTypes: testsuite.UnionOfComplexTypes{
			ComplexTypeUnion: testsuite.UnionOfComplexTypes_ComplexTypeUnion{
				Fruits: &apple,
			},
		},
		UnionOfSameTypes: testsuite.UnionOfSameTypes{
			SameTypesUnion: testsuite.UnionOfSameTypes_SameTypesUnion{
				Greeting: &hello,
			},
			UnionWithArrayMembers: testsuite.UnionOfSameTypes_UnionWithArrayMembers{
				FruitArray: &[]conflictresolution.Fruits{
					conflictresolution.Fruits_ORANGE,
					conflictresolution.Fruits_APPLE,
				},
			},
			UnionWithMapMembers: testsuite.UnionOfSameTypes_UnionWithMapMembers{
				IntMap: &map[string]int32{
					"one": 1,
				},
			},
		},
	}

	t.Run("json", func(t *testing.T) {
		testJsonEncoding(t, expected, `{
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

	t.Run("urlEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(anotherUnionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:APPLE)),arrayOfMaps:(arrayOfMaps:List((one:1),(two:2))),mapOfInts:(mapOfInts:(one:1)),recordWithProps:(integer:5),unionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:ORANGE)),unionOfPrimitives:(primitivesUnion:(int:5)),unionOfSameTypes:(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1))))`, false)
	})

	t.Run("reducedEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(anotherUnionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:APPLE)),arrayOfMaps:(arrayOfMaps:List((one:1),(two:2))),mapOfInts:(mapOfInts:(one:1)),recordWithProps:(integer:5),unionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:ORANGE)),unionOfPrimitives:(primitivesUnion:(int:5)),unionOfSameTypes:(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1))))`, true)
	})
}

func TestMapEncoding(t *testing.T) {
	expected := &testsuite.Optionals{
		OptionalMap: &map[string]int32{
			"one": 1,
			"two": 2,
		},
	}

	t.Run("multipleElements", func(t *testing.T) {
		var serialized strings.Builder
		require.NoError(t, expected.RestLiEncode(protocol.RestLiQueryEncoder, &serialized))

		if serialized.String() != `(optionalMap:(one:1,two:2))` && serialized.String() != `(optionalMap:(two:2,one:1))` {
			t.Fail()
		}
	})

	expected = &testsuite.Optionals{OptionalMap: &map[string]int32{}}
	t.Run("empty", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalMap:())`, false)
	})
}

func TestArrayEncoding(t *testing.T) {
	expected := &testsuite.Optionals{OptionalArray: &[]int32{1, 2}}

	t.Run("urlEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalArray:List(1,2))`, false)
	})
	t.Run("reducedEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalArray:List(1,2))`, true)
	})

	expected = &testsuite.Optionals{OptionalArray: &[]int32{}}
	t.Run("empty", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalArray:List())`, false)
	})
}

func TestEmpty(t *testing.T) {
	expected := &testsuite.Optionals{OptionalArray: &[]int32{1, 2}}

	t.Run("urlEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalArray:List(1,2))`, false)
	})
	t.Run("reducedEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalArray:List(1,2))`, true)
	})
}

func TestEmptyStringAndBytes(t *testing.T) {
	expected := &testsuite.Optionals{
		OptionalBytes:  new(protocol.Bytes),
		OptionalString: new(string),
	}

	t.Run("urlEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalBytes:'',optionalString:'')`, false)
	})
	t.Run("reducedEncode", func(t *testing.T) {
		testRestliEncoding(t, expected, `(optionalBytes:'',optionalString:'')`, true)
	})
}

func testJsonEncoding(t *testing.T, expected interface{}, expectedRawJson string) {
	serialized, err := json.Marshal(expected)
	require.NoError(t, err)

	var expectedRaw map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(expectedRawJson), &expectedRaw))
	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(serialized, &raw))
	require.Equal(t, expectedRaw, raw)

	v := reflect.New(reflect.TypeOf(expected)).Interface()
	require.NoError(t, json.Unmarshal(serialized, v))
	require.Equal(t, expected, reflect.ValueOf(v).Elem().Interface())
}

func testRestliEncoding(t *testing.T, expected protocol.RestLiEncodable, expectedRawEncoded string, reducedEncoding bool) {
	var encoder *protocol.RestLiCodec
	if reducedEncoding {
		encoder = protocol.RestLiReducedEncoder
	} else {
		encoder = protocol.RestLiQueryEncoder
	}

	var serialized strings.Builder
	require.NoError(t, expected.RestLiEncode(encoder, &serialized))

	require.Equal(t, expectedRawEncoded, serialized.String())
}
