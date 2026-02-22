package std

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MaskStr(t *testing.T) {
	t.Run("basic cases", func(t *testing.T) {
		tCases := []struct {
			input    string
			expected string
		}{
			{"", ""},
			{"s", "*"},
			{"Secret-codE", "***********"},
			{"Секретный коД", "*************"},
		}

		for _, tCase := range tCases {
			result := MaskStr(tCase.input)
			assert.Equal(t, tCase.expected, result)
		}
	})
}

func Test_MaskStrNotFirst(t *testing.T) {
	t.Run("basic cases", func(t *testing.T) {
		tCases := []struct {
			input    string
			expected string
		}{
			{"", ""},
			{"s", "s"},
			{"Secret-codE", "S**********"},
			{"Секретный коД", "С************"},
		}

		for _, tCase := range tCases {
			result := MaskStrNotFirst(tCase.input)
			assert.Equal(t, tCase.expected, result)
		}
	})
}

func Test_MaskStrNotFirstLast(t *testing.T) {
	t.Run("basic cases", func(t *testing.T) {
		tCases := []struct {
			input    string
			expected string
		}{
			{"", ""},
			{"s", "s"},
			{"se", "se"},
			{"Secret-codE", "S*********E"},
			{"Секретный коД", "С***********Д"},
		}

		for _, tCase := range tCases {
			result := MaskStrNotFirstLast(tCase.input)
			assert.Equal(t, tCase.expected, result)
		}
	})
}
