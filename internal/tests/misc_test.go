package tests

import (
	"crypto/md5"
	"testing"
	"time"

	nativetestsuite "github.com/PapaCharlie/go-restli/internal/tests/native/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated_extras/extras"
	"github.com/PapaCharlie/go-restli/protocol/equals"
	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
	"github.com/stretchr/testify/require"
)

func TestInclude(t *testing.T) {
	expected := &testsuite.Include{
		PrimitiveField: testsuite.PrimitiveField{Integer: int32(1)},
		F1:             4.27,
	}
	testJsonEncoding(t, expected, testsuite.UnmarshalRestLiInclude, `{
  "f1": 4.27,
  "integer": 1
}`)
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
	require.Equal(t, nativetestsuite.Fruits_APPLE, *d.DefaultEnum)
	require.Equal(t, testsuite.Fixed5{1, 2, 3, 4, 5}, *d.DefaultFixed)
	require.Equal(t, testsuite.PrimitiveField{Integer: 10}, *d.DefaultRecord)
	require.Equal(t, []int32{1, 3, 5}, *d.DefaultArray)
	require.Equal(t, map[string]int32{"a": 1, "b": 2}, *d.DefaultMap)
	require.Equal(t, testsuite.Defaults_DefaultUnion{Int: &five}, *d.DefaultUnion)
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
		}, []nativetestsuite.Fruits{
			nativetestsuite.Fruits_APPLE,
			nativetestsuite.Fruits_APPLE,
			nativetestsuite.Fruits_ORANGE,
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
		var t1, t2, t3 nativetestsuite.Temperature = 1, 1, 2
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

	t.Run("time", func(t *testing.T) {
		t1, t2 := extras.Time(time.Now()), extras.Time(time.Now())
		t3 := t1
		testEquality(t, [][]bool{
			{true, false, true},
			{false, true, false},
			{true, false, true},
		}, []extras.Time{t1, t2, t3})
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

// ensures the Hex method on MD5 is never deleted by the code generator
func TestMD5Hex(t *testing.T) {
	m := extras.MD5(md5.Sum([]byte("abc")))
	require.Equal(t, "900150983cd24fb0d6963f7d28e17f72", m.Hex())
}
