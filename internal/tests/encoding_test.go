package tests

import (
	"encoding/json"
	"reflect"
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/generated/testsuite"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/restliencoding"
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

	testRestliEncoding(t, expected,
		`(primitiveBytes:@ABC🕴,primitiveDouble:66.5,primitiveFloat:52.5,primitiveInteger:1,primitiveLong:23,primitiveString:a string%2C%28%29%27)`,
		`primitiveBytes=%40ABC%F0%9F%95%B4%01&primitiveDouble=66.5&primitiveFloat=52.5&primitiveInteger=1&primitiveLong=23&primitiveString=a+string%2C%28%29%27`,
	)
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

	testRestliEncoding(t, expected,
		`(anotherUnionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:APPLE)),arrayOfMaps:(arrayOfMaps:List((one:1),(two:2))),mapOfInts:(mapOfInts:(one:1)),recordWithProps:(integer:5),unionOfComplexTypes:(complexTypeUnion:(testsuite.Fruits:ORANGE)),unionOfPrimitives:(primitivesUnion:(int:5)),unionOfSameTypes:(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1))))`,
		`anotherUnionOfComplexTypes=(complexTypeUnion:(testsuite.Fruits:APPLE))&arrayOfMaps=(arrayOfMaps:List((one:1),(two:2)))&mapOfInts=(mapOfInts:(one:1))&recordWithProps=(integer:5)&unionOfComplexTypes=(complexTypeUnion:(testsuite.Fruits:ORANGE))&unionOfPrimitives=(primitivesUnion:(int:5))&unionOfSameTypes=(sameTypesUnion:(greeting:Hello),unionWithArrayMembers:(fruitArray:List(ORANGE,APPLE)),unionWithMapMembers:(intMap:(one:1)))`,
	)
}

func TestMapEncoding(t *testing.T) {
	expected := &testsuite.Optionals{
		OptionalMap: &map[string]int32{
			"one": 1,
			"two": 2,
		},
	}

	t.Run("multipleElements", func(t *testing.T) {
		encoder := restliencoding.NewQueryParamsEncoder()
		require.NoError(t, expected.RestLiEncode(encoder))

		serialized := encoder.Finalize()
		if serialized != `optionalMap=(one:1,two:2)` && serialized != `optionalMap=(two:2,one:1)` {
			t.Fail()
		}
	})

	expected = &testsuite.Optionals{OptionalMap: &map[string]int32{}}
	encoder := restliencoding.NewPathEncoder().Encoder
	require.NoError(t, expected.RestLiEncode(encoder))
	require.Equal(t, `(optionalMap:())`, encoder.Finalize())
}

func TestArrayEncoding(t *testing.T) {
	expected := &testsuite.Optionals{OptionalArray: &[]int32{1, 2}}

	testRestliEncoding(t, expected,
		`(optionalArray:List(1,2))`,
		`optionalArray=List(1,2)`,
	)

	expected = &testsuite.Optionals{OptionalArray: &[]int32{}}
	encoder := restliencoding.NewPathEncoder().Encoder
	require.NoError(t, expected.RestLiEncode(encoder))
	require.Equal(t, `(optionalArray:List())`, encoder.Finalize())
}

func TestEmptyStringAndBytes(t *testing.T) {
	expected := &testsuite.Optionals{
		OptionalBytes:  new(protocol.Bytes),
		OptionalString: new(string),
	}

	testRestliEncoding(t, expected,
		`(optionalBytes:'',optionalString:'')`,
		`optionalBytes=''&optionalString=''`,
	)
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

func testRestliEncoding(t *testing.T, expected restliencoding.Encodable, expectedHeaderEncoded, expectedQueryEncoded string) {
	tests := []struct {
		Name     string
		Expected string
		Encoder  func() *restliencoding.Encoder
	}{
		{
			Name:     "headerEncoded",
			Expected: expectedHeaderEncoded,
			Encoder:  restliencoding.NewHeaderEncoder,
		},
		{
			Name:     "queryEncoded",
			Expected: expectedQueryEncoded,
			Encoder:  restliencoding.NewQueryParamsEncoder,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			encoder := test.Encoder()
			require.NoError(t, expected.RestLiEncode(encoder))
			require.Equal(t, test.Expected, encoder.Finalize())
		})
	}
}
