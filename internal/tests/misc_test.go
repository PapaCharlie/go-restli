package tests

import (
	"crypto/md5"
	"testing"

	"github.com/PapaCharlie/go-restli/fnv1a"
	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	collectionwithtyperefkey "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras/collectionWithTyperefKey"
	"github.com/PapaCharlie/go-restli/protocol"
	"github.com/PapaCharlie/go-restli/protocol/equals"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/PapaCharlie/go-restli/protocol/stdtypes"
	"github.com/stretchr/testify/require"
)

func TestInclude(t *testing.T) {
	expected := &testsuite.Include{
		PrimitiveField: testsuite.PrimitiveField{Integer: int32(1)},
		F1:             4.27,
	}
	testEncoding(t, expected, `{
  "f1": 4.27,
  "integer": 1
}`, `(f1:4.27,integer:1)`)
}

// TestDefaults tests that default values are loaded correctly (see
// rest.li-test-suite/client-testsuite/schemas/testsuite/Defaults.pdsc) for the default values used here
func TestDefaults(t *testing.T) {
	expected := &testsuite.Defaults{
		DefaultInteger: protocol.Int32Pointer(1),
		DefaultLong:    protocol.Int64Pointer(23),
		DefaultFloat:   protocol.Float32Pointer(52.5),
		DefaultDouble:  protocol.Float64Pointer(66.5),
		DefaultBytes:   protocol.BytesPointer([]byte("@ABC")),
		DefaultString:  protocol.StringPointer("default string"),
		DefaultEnum:    conflictresolution.Fruits_APPLE.Pointer(),
		DefaultFixed:   testsuite.Fixed5{1, 2, 3, 4, 5}.Pointer(),
		DefaultRecord:  &testsuite.PrimitiveField{Integer: 10},
		DefaultArray:   &[]int32{1, 3, 5},
		DefaultMap:     &map[string]int32{"a": 1, "b": 2},
		DefaultUnion:   &testsuite.Defaults_DefaultUnion{Int: protocol.Int32Pointer(5)},
	}
	d := testsuite.NewDefaultsWithDefaultValues()
	require.Equal(t, expected, d)
	// flex the Equals code a little
	require.True(t, expected.Equals(d))

	moreExpected := &extras.MoreDefaults{
		DefaultRecord:  extras.DefaultTyperef{Foo: extras.Temperature(42).Pointer()},
		EmptyArray:     new([]string),
		EmptyMap:       &map[string]string{},
		DefaultBoolean: protocol.BoolPointer(true),
	}
	moreD := extras.NewMoreDefaultsWithDefaultValues()
	require.Equal(t, moreExpected, moreD)
	// flex the Equals code some more
	require.True(t, moreExpected.Equals(moreD))
}

func TestEnum(t *testing.T) {
	_, err := conflictresolution.GetFruitsFromString("BANANA")
	require.IsType(t, new(protocol.UnknownEnumValue), err)
	const illegal = conflictresolution.Fruits(42)
	err = illegal.MarshalRestLi(restlicodec.NoopWriter)
	require.IsType(t, new(protocol.IllegalEnumConstant), err)
	require.True(t, illegal.ComputeHash().Equals(fnv1a.ZeroHash()))
	require.False(t, illegal.Equals(illegal))
}

func testEquality[T equals.Equatable[T]](t *testing.T, tests [][]bool, data []T) {
	for i, row := range tests {
		for j, expected := range row {
			a, b := data[i], data[j]
			require.Equal(t, expected, a.Equals(b), "Equals(%d, %d)", a, b)
			if expected {
				require.Equal(t, a, b)
			} else {
				require.NotEqual(t, a, b)
			}
		}
	}
}

func TestEquals(t *testing.T) {
	t.Run("enum", func(t *testing.T) {
		testEquality(t, [][]bool{
			{true, true, false},
			{true, true, false},
			{false, false, true},
		}, []conflictresolution.Fruits{
			conflictresolution.Fruits_APPLE,
			conflictresolution.Fruits_APPLE,
			conflictresolution.Fruits_ORANGE,
		})
	})

	t.Run("fixed", func(t *testing.T) {
		testEquality(t, [][]bool{
			{true, true, false},
			{true, true, false},
			{false, false, true},
		}, []*testsuite.Fixed5{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 3, 4},
			{1, 2, 3, 4, 5},
		})
	})

	t.Run("union", func(t *testing.T) {
		i1, i2, i3, l := int32(1), int32(1), int32(2), int64(2)
		testEquality(t, [][]bool{
			{true, true, true, false, false},
			{true, true, true, false, false},
			{true, true, true, false, false},
			{false, false, false, true, false},
			{false, false, false, false, true},
		}, []*testsuite.UnionOfPrimitives_PrimitivesUnion{
			{Int: &i1},
			{Int: &i1},
			{Int: &i2},
			{Int: &i3},
			{Long: &l},
		})
	})

	t.Run("record", func(t *testing.T) {
		t1, t2, t3 := extras.Temperature(1), extras.Temperature(1), extras.Temperature(2)
		testEquality(t, [][]bool{
			{true, false, false, false, false},
			{false, true, true, true, false},
			{false, true, true, true, false},
			{false, true, true, true, false},
			{false, false, false, false, true},
		}, []*extras.DefaultTyperef{
			{Foo: nil},
			{Foo: &t1},
			{Foo: &t1},
			{Foo: &t2},
			{Foo: &t3},
		})
	})
}

func TestReadInterface(t *testing.T) {
	t.Run("ror2", func(t *testing.T) {
		read := func(t *testing.T, s string) interface{} {
			reader := newRor2Reader(t, s)
			i, err := reader.ReadInterface()
			require.NoError(t, err)

			return i
		}

		t.Run("string", func(t *testing.T) {
			require.Equal(t, "asd", read(t, "asd"))
			require.Equal(t, "11", read(t, "11"))
			require.Equal(t, "43.9", read(t, "43.9"))
			require.Equal(t, "false", read(t, "false"))
		})

		t.Run("map", func(t *testing.T) {
			require.Equal(t, map[string]interface{}{
				"primitive": "1",
				"map": map[string]interface{}{
					"one": "1",
					"two": "2",
				},
				"array": []interface{}{
					map[string]interface{}{"foo": "bar"},
				},
			}, read(t, "(primitive:1,map:(one:1,two:2),array:List((foo:bar)))"))
			require.Equal(t, map[string]interface{}{}, read(t, "()"))
		})

		t.Run("array", func(t *testing.T) {
			require.Equal(t, []interface{}{"1", "2", "3"}, read(t, "List(1,2,3)"))
			require.Equal(t, []interface{}(nil), read(t, "List()"))
		})
	})

	t.Run("json", func(t *testing.T) {
		read := func(t *testing.T, s string) interface{} {
			reader := newJsonReader(t, s)
			i, err := reader.ReadInterface()
			require.NoError(t, err)

			return i
		}

		t.Run("primitives", func(t *testing.T) {
			require.Equal(t, "asd", read(t, `"asd"`))
			require.Equal(t, 43.9, read(t, `43.9`))
			require.Equal(t, false, read(t, `false`))
		})

		t.Run("map", func(t *testing.T) {
			require.Equal(t, map[string]interface{}{
				"primitive": 1.,
				"map": map[string]interface{}{
					"one": 1.,
					"two": 2.,
				},
				"array": []interface{}{
					map[string]interface{}{"foo": "bar"},
				},
			}, read(t, `{
                   "primitive":1,
                   "map": {
                     "one": 1,
                     "two": 2
                   },
                   "array": [{
                     "foo": "bar"
                   }]
                 }`),
			)
			require.Equal(t, map[string]interface{}{}, read(t, `{}`))
		})

		t.Run("array", func(t *testing.T) {
			require.Equal(t, []interface{}{false, true, false}, read(t, "[false,true,false]"))
			require.Equal(t, []interface{}{}, read(t, "[]"))
		})
	})
}

// ensures the Hex method on MD5 is never deleted by the code generator
func TestMD5Hex(t *testing.T) {
	m := extras.MD5(md5.Sum([]byte("abc")))
	require.Equal(t, "900150983cd24fb0d6963f7d28e17f72", m.Hex())
}

func TestQueryParamsReader(t *testing.T) {
	q, err := restlicodec.ParseQueryParams("b=()&c=(string:foo)")
	require.NoError(t, err)

	c := new(extras.SinglePrimitiveField)
	err = q.ReadRecord(restlicodec.RequiredFields{"a", "b", "c"}, func(reader restlicodec.Reader, field string) (err error) {
		if field == "c" {
			return c.UnmarshalRestLi(reader)
		} else {
			return new(extras.SinglePrimitiveField).UnmarshalRestLi(reader)
		}
	})
	require.Equal(t, &restlicodec.MissingRequiredFieldsError{Fields: []string{"a", "b.string"}}, err)
	require.Equal(t, &extras.SinglePrimitiveField{String: "foo"}, c)
}

func TestEmbeddedPagingContext(t *testing.T) {
	var start, count int32
	start = 10
	count = 20
	tests := []struct {
		name     string
		params   collectionwithtyperefkey.FindBySearchParams
		expected string
	}{
		{
			name: "empty context",
			params: collectionwithtyperefkey.FindBySearchParams{
				Keyword: "foo",
			},
			expected: "keyword=foo&q=search",
		},
		{
			name: "start only",
			params: collectionwithtyperefkey.FindBySearchParams{
				PagingContext: stdtypes.PagingContext{
					Start: &start,
				},
				Keyword: "foo",
			},
			expected: "keyword=foo&q=search&start=10",
		},
		{
			name: "count only",
			params: collectionwithtyperefkey.FindBySearchParams{
				PagingContext: stdtypes.PagingContext{
					Count: &count,
				},
				Keyword: "foo",
			},
			expected: "count=20&keyword=foo&q=search",
		},
		{
			name: "full context",
			params: collectionwithtyperefkey.FindBySearchParams{
				PagingContext: stdtypes.NewPagingContext(start, count),
				Keyword:       "foo",
			},
			expected: "count=20&keyword=foo&q=search&start=10",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.params.EncodeQueryParams()
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}
