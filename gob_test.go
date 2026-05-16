package std

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GobBb(t *testing.T) {
	tCases := []struct {
		value   any
		encoded []byte
	}{
		{
			value:   42,
			encoded: []byte{0x3, 0x4, 0x0, 0x54},
		},
		{
			value:   "hello",
			encoded: []byte{0x8, 0xc, 0x0, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			value: []int{1, 2, 3},
			encoded: []byte{
				0xb, 0x7f, 0x2, 0x1, 0x2, 0xff, 0x80, 0x0, 0x1, 0x4, 0x0, 0x0, 0x7, 0xff, 0x80, 0x0, 0x3,
				0x2, 0x4, 0x6,
			},
		},
		{
			value: map[string]int{"a": 1},
			encoded: []byte{
				0xe, 0xff, 0x81, 0x4, 0x1, 0x2, 0xff, 0x82, 0x0, 0x1, 0xc, 0x1, 0x4, 0x0, 0x0, 0x7, 0xff, 0x82, 0x0,
				0x1, 0x1, 0x61, 0x2,
			},
		},
		{
			value: struct{ X, Y int }{10, 20},
			encoded: []byte{
				0x18, 0xff, 0x83, 0x3, 0x1, 0x2, 0xff, 0x84, 0x0, 0x1, 0x2, 0x1, 0x1, 0x58, 0x1, 0x4, 0x0, 0x1, 0x1,
				0x59, 0x1, 0x4, 0x0, 0x0, 0x0, 0x7, 0xff, 0x84, 0x1, 0x14, 0x1, 0x28, 0x0,
			},
		},
		{
			value: &struct{ Z string }{Z: "z"},
			encoded: []byte{
				0x12, 0xff, 0x85, 0x3, 0x1, 0x2, 0xff, 0x86, 0x0, 0x1, 0x1, 0x1, 0x1, 0x5a, 0x1, 0xc, 0x0, 0x0, 0x0,
				0x6, 0xff, 0x86, 0x1, 0x1, 0x7a, 0x0,
			},
		},
	}

	t.Run("encode", func(t *testing.T) {
		for _, tCase := range tCases {
			tCaseName := fmt.Sprintf("TCase: %+v", tCase.value)

			encoded, err := GobEncodeBytes(tCase.value)

			assert.NoError(t, err)
			assert.Equal(t, tCase.encoded, encoded, tCaseName)
		}
	})

	t.Run("decode", func(t *testing.T) {
		for _, tCase := range tCases {
			tCaseName := fmt.Sprintf("TCase: %+v", tCase.value)

			switch tcv := tCase.value.(type) {
			case int:
				v, err := GobDecodeBytes[int](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case string:
				v, err := GobDecodeBytes[string](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case []int:
				v, err := GobDecodeBytes[[]int](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case map[string]int:
				v, err := GobDecodeBytes[map[string]int](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case struct{ X, Y int }:
				v, err := GobDecodeBytes[struct{ X, Y int }](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case *struct{ Z string }:
				v, err := GobDecodeBytes[*struct{ Z string }](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
				assert.Equal(t, *tcv, *v, tCaseName)
			default:
				panic("unknown case")
			}
		}
	})

	t.Run("round trip", func(t *testing.T) {
		for _, tCase := range tCases {
			tCaseName := fmt.Sprintf("TCase: %+v", tCase.value)

			encoded, err := GobEncodeBytes(tCase.value)
			assert.NoError(t, err, tCaseName)

			switch tcv := tCase.value.(type) {
			case int:
				v, err := GobDecodeBytes[int](encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case string:
				v, err := GobDecodeBytes[string](encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case []int:
				v, err := GobDecodeBytes[[]int](encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case map[string]int:
				v, err := GobDecodeBytes[map[string]int](encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case struct{ X, Y int }:
				v, err := GobDecodeBytes[struct{ X, Y int }](encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
			case *struct{ Z string }:
				v, err := GobDecodeBytes[*struct{ Z string }](tCase.encoded)
				assert.NoError(t, err, tCaseName)
				assert.Equal(t, tcv, v, tCaseName)
				assert.Equal(t, *tcv, *v, tCaseName)
			default:
				panic("unknown case")
			}
		}
	})

	t.Run("nil value encode", func(t *testing.T) {
		b, err := GobEncodeBytes(nil)
		assert.Nil(t, b)
		assert.Error(t, err)
	})

	t.Run("invalid data decode", func(t *testing.T) {
		v, err := GobDecodeBytes[int]([]byte{0x8, 0xc, 0x0, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f})
		assert.Error(t, err)
		assert.Equal(t, v, *new(int))
	})
}
