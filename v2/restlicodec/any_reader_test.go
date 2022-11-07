package restlicodec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Foo              = "foo"
	FortyTwo         = 42
	FortyTwoPointTwo = 42.2
)

func testRead[T any](t *testing.T, expected T, source any, reader GenericUnmarshaler[T]) {
	run := func(source any) {
		actual, err := reader(NewInterfaceReader(source))
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
	t.Run(fmt.Sprintf("%s(%v)", reflect.TypeOf(source), source), func(t *testing.T) {
		run(source)
	})
	sourcePointer := reflect.New(reflect.TypeOf(source))
	sourcePointer.Elem().Set(reflect.ValueOf(source))

	t.Run(fmt.Sprintf("%s(%v)", reflect.TypeOf(sourcePointer.Interface()), source), func(t *testing.T) {
		run(sourcePointer.Interface())
	})
}

func Test_anyReader_ReadBool(t *testing.T) {
	testRead(t, true, "true", Reader.ReadBool)
	testRead(t, true, true, Reader.ReadBool)
}

func Test_anyReader_ReadBytes(t *testing.T) {
	testRead(t, []byte(Foo), []byte(Foo), Reader.ReadBytes)
	testRead(t, []byte(Foo), Foo, Reader.ReadBytes)
}

func Test_anyReader_ReadString(t *testing.T) {
	testRead(t, Foo, []byte(Foo), Reader.ReadString)
	testRead(t, Foo, Foo, Reader.ReadString)
}

var intTests = []any{
	int(42),
	int8(42),
	int16(42),
	int32(42),
	int64(42),
	float32(42),
	float64(42),
	"42",
	[]byte("42"),
}

func testReadInt[T any](t *testing.T, expected T, unmarshaler GenericUnmarshaler[T]) {
	for _, v := range intTests {
		testRead[T](t, expected, v, unmarshaler)
	}
}

var floatTests = []any{
	float32(42.2),
	float64(42.2),
	"42.2",
	[]byte("42.2"),
}

func testReadFloat[T float64 | float32](t *testing.T, expected T, reader GenericUnmarshaler[T]) {
	for _, v := range floatTests {
		t.Run(fmt.Sprintf("%s(%v)", reflect.TypeOf(v).String(), v), func(t *testing.T) {
			actual, err := reader(NewInterfaceReader(v))
			require.NoError(t, err)

			if expected == actual {
				return
			}

			diff := expected - actual
			if diff < 0 {
				diff = -diff
			}
			require.Less(t, diff, 0.001)
		})
	}
}

func Test_anyReader_ReadFloat32(t *testing.T) {
	testReadInt[float32](t, FortyTwo, Reader.ReadFloat32)
	testReadFloat[float32](t, FortyTwoPointTwo, Reader.ReadFloat32)
}

func Test_anyReader_ReadFloat64(t *testing.T) {
	testReadInt[float64](t, FortyTwo, Reader.ReadFloat64)
	testReadFloat[float64](t, FortyTwoPointTwo, Reader.ReadFloat64)
}

func Test_anyReader_ReadInt(t *testing.T) {
	testReadInt[int](t, FortyTwo, Reader.ReadInt)
}

func Test_anyReader_ReadInt32(t *testing.T) {
	testReadInt[int32](t, FortyTwo, Reader.ReadInt32)
}

func Test_anyReader_ReadInt64(t *testing.T) {
	testReadInt[int64](t, FortyTwo, Reader.ReadInt64)
}

func Test_anyReader_ReadRecord(t *testing.T) {
	testRead[*Object](t, &Object{Status: 500}, map[string]any{"status": 500.5}, UnmarshalRestLi[*Object])
}

func Test_anyReader_ReadArray(t *testing.T) {
	testRead[[]int](t, []int{27, 42}, []int32{27, 42}, func(reader Reader) ([]int, error) {
		return ReadArray[int](reader, Reader.ReadInt)
	})
}

func Test_anyReader_ReadMap(t *testing.T) {
	testRead[map[string]int](
		t,
		map[string]int{"foo": 27, "bar": 42},
		map[string]int32{"foo": 27, "bar": 42},
		func(reader Reader) (map[string]int, error) {
			return ReadMap[int](reader, Reader.ReadInt)
		},
	)
}

func Test_anyReader_noPanic(t *testing.T) {
	requireError := func(_ any, err error) {
		require.Error(t, err)
		_, ok := err.(*InvalidTypeError)
		require.True(t, ok)
	}

	r := NewInterfaceReader(make(chan int))

	requireError(r.ReadInt())
	requireError(r.ReadInt32())
	requireError(r.ReadInt64())
	requireError(r.ReadFloat32())
	requireError(r.ReadFloat64())
	requireError(r.ReadBool())
	requireError(r.ReadString())
	requireError(r.ReadBytes())
	requireError(nil, r.ReadMap(func(reader Reader, field string) (err error) { return nil }))
	requireError(nil, NewInterfaceReader(map[int]int{}).ReadMap(func(reader Reader, field string) (err error) { return nil }))
	requireError(nil, r.ReadArray(func(reader Reader) (err error) { return nil }))
}
