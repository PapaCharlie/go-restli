package tests

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/stretchr/testify/require"
)

func TestInclude(t *testing.T) {
	expected := &testsuite.Include{
		Integer: int32(1),
		F1:      4.27,
	}
	testJsonEncoding(t, expected, new(testsuite.Include), `{ "integer": 1, "f1": 4.27 }`)
}

// TestDefaults tests that default values are loaded correctly (see
// rest.li-test-suite/client-testsuite/schemas/testsuite/Defaults.pdsc) for the default values used here
func TestDefaults(t *testing.T) {
	five := int32(5)
	d := testsuite.NewDefaultsWithDefaultValues()
	require.Equal(t, int32(1), *d.DefaultInteger)
	require.Equal(t, int64(23), *d.DefaultLong)
	require.Equal(t, float32(52.5), *d.DefaultFloat)
	require.Equal(t, float64(66.5), *d.DefaultDouble)
	require.Equal(t, []byte("@ABC"), *d.DefaultBytes)
	require.Equal(t, string("default string"), *d.DefaultString)
	require.Equal(t, conflictresolution.Fruits_APPLE, *d.DefaultEnum)
	require.Equal(t, testsuite.Fixed5{1, 2, 3, 4, 5}, *d.DefaultFixed)
	require.Equal(t, testsuite.PrimitiveField{Integer: 10}, *d.DefaultRecord)
	require.Equal(t, []int32{1, 3, 5}, *d.DefaultArray)
	require.Equal(t, map[string]int32{"a": 1, "b": 2}, *d.DefaultMap)
	require.Equal(t, testsuite.Defaults_DefaultUnion{Int: &five}, *d.DefaultUnion)
}

func TestEquals(t *testing.T) {
	testEquality := func(t *testing.T, tests [][]bool, supplier func(index int) restliObject) {
		for i, row := range tests {
			for j, expected := range row {
				a, b := supplier(i), supplier(j)
				require.Equal(t, expected, a.Equals(b), "Equals(%d, %d)", i, j)
				if expected {
					require.Equal(t, a, b)
				} else {
					require.NotEqual(t, a, b)
				}
			}
		}
	}

	t.Run("enum", func(t *testing.T) {
		data := []conflictresolution.Fruits{
			conflictresolution.Fruits_APPLE,
			conflictresolution.Fruits_APPLE,
			conflictresolution.Fruits_ORANGE,
		}
		testEquality(t, [][]bool{
			{true, true, false},
			{true, true, false},
			{false, false, true},
		}, func(i int) restliObject {
			return &data[i]
		})
	})

	t.Run("fixed", func(t *testing.T) {
		data := []*testsuite.Fixed5{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 3, 4},
			{1, 2, 3, 4, 5},
		}
		testEquality(t, [][]bool{
			{true, true, false},
			{true, true, false},
			{false, false, true},
		}, func(i int) restliObject {
			return data[i]
		})
	})

	t.Run("union", func(t *testing.T) {
		i1, i2, i3, l := int32(1), int32(1), int32(2), int64(2)
		data := []*testsuite.UnionOfPrimitives_PrimitivesUnion{
			{Int: &i1},
			{Int: &i1},
			{Int: &i2},
			{Int: &i3},
			{Long: &l},
		}
		testEquality(t, [][]bool{
			{true, true, true, false, false},
			{true, true, true, false, false},
			{true, true, true, false, false},
			{false, false, false, true, false},
			{false, false, false, false, true},
		}, func(i int) restliObject {
			return data[i]
		})
	})

	t.Run("record", func(t *testing.T) {
		t1, t2, t3 := extras.Temperature(1), extras.Temperature(1), extras.Temperature(2)
		data := []*extras.DefaultTyperef{
			{Foo: nil},
			{Foo: &t1},
			{Foo: &t1},
			{Foo: &t2},
			{Foo: &t3},
		}
		testEquality(t, [][]bool{
			{true, false, false, false, false},
			{false, true, true, true, false},
			{false, true, true, true, false},
			{false, true, true, true, false},
			{false, false, false, false, true},
		}, func(i int) restliObject {
			return data[i]
		})
	})
}

func TestReadInterface(t *testing.T) {
	t.Run("ror2", func(t *testing.T) {
		read := func(t *testing.T, s string) interface{} {
			reader, err := restlicodec.NewRor2Reader(s)
			require.NoError(t, err)

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
			reader := restlicodec.NewJsonReader([]byte(s))

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
