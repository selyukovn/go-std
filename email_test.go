package std

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func testDataProvider_Email_correctValues() []string {
	return []string{
		"te@te.te",
		"te1@te.te",
		"te1@te.te.te",
		"te2te@te.te",
		"te2te@te.te.te",
		"3te@te.te",
		"3te@te.te.te",
		"t.te@te.te",
		"t.te@te.te.te",
		"te.te@te.te",
		"te.te@te.te.te",
		"te.t@te.te",
		"te.t@te.te.te",
		"t+te@te.te",
		"te+te@te.te",
		"te+t@te.te",
		"t.e.s.t+t.e.s.t@te.te",
		"t.e.s.t+t.e.s.t@te.te.te",
		"t_e_s_t@te.te",
		"test-@te.te",
		"test_@te.te",
		"test+te-@te.te",
		"test+te_@te.te",
		"test-+te-@te.te",
		"test-+te_@te.te",
		"test_+te-@te.te",
		"test_+te_@te.te",
		"123456@t.test.te.te",
		"123456@test.t.te.te",
		"123456@test.te.t.te",
		"123456@test.te.te.t",
		strings.Repeat("a", 249) + "@aa.aa",
	}
}

func testDataProvider_Email_incorrectValues() []string {
	return []string{
		"",
		" ",
		"   ",
		"a",
		"aaaaaaaa",
		"a@a.a",
		"a@aa.aa",
		" aa@aa.aa",
		"aa@aa.aa ",
		" aa@aa.aa ",
		"aa @ aa . aa",
		"__@__.__",
		"a_@aa.aa",
		"aa@a_.aa",
		"aa@aa.a_",
		strings.Repeat("a", 250) + "@aa.aa",
		strings.Repeat("a", 1000) + "@aa.aa",
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// EmailFromString
// ---------------------------------------------------------------------------------------------------------------------

func Test_EmailFromString(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		tCases := testDataProvider_Email_correctValues()
		for _, tCase := range tCases {
			email, err := EmailFromString(tCase)
			// no error
			assert.NoError(t, err)
			// email not nil
			assert.NotEqual(t, EmailNil, email)
			assert.False(t, email.IsNil())
			assert.Equal(t, tCase, email.String())
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		tCases := testDataProvider_Email_incorrectValues()
		for _, tCase := range tCases {
			email, err := EmailFromString(tCase)
			// error
			assert.Error(t, err)
			// email nil
			assert.Equal(t, EmailNil, email)
			assert.True(t, email.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// EmailFromStringMust
// ---------------------------------------------------------------------------------------------------------------------

func Test_EmailFromStringMust(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		tCases := testDataProvider_Email_correctValues()
		for _, tCase := range tCases {
			// no panic
			assert.NotPanics(t, func() { EmailFromStringMust(tCase) })
			// email not nil
			email := EmailFromStringMust(tCase)
			assert.NotEqual(t, EmailNil, email)
			assert.False(t, email.IsNil())
			assert.Equal(t, tCase, email.String())
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		tCases := testDataProvider_Email_incorrectValues()
		for _, tCase := range tCases {
			assert.Panics(t, func() { EmailFromStringMust(tCase) })
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// IsNil
// ---------------------------------------------------------------------------------------------------------------------

func Test_Email_IsNil(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, EmailNil.IsNil())
	})

	t.Run("false", func(t *testing.T) {
		tCases := testDataProvider_Email_correctValues()
		for _, tCase := range tCases {
			email, _ := EmailFromString(tCase)
			assert.False(t, email.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------------------------------------------------

func Test_Email_String(t *testing.T) {
	tCases := testDataProvider_Email_correctValues()
	for _, tCase := range tCases {
		email, _ := EmailFromString(tCase)
		assert.Equal(t, tCase, email.String())
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Name
// ---------------------------------------------------------------------------------------------------------------------

func Test_Email_Name(t *testing.T) {
	tCases := map[string]string{
		"te@te.te":              "te",
		"te1@te.te":             "te1",
		"te2te@te.te.te":        "te2te",
		"3te@te.te":             "3te",
		"t+te@te.te":            "t+te",
		"t.e.s.t+t.e.s.t@te.te": "t.e.s.t+t.e.s.t",
		"t_e_s_t@te.te":         "t_e_s_t",
		"test-@te.te":           "test-",
		"test_@te.te":           "test_",
		"test-+te-@te.te":       "test-+te-",
		"123456@t.test.te.te":   "123456",
	}

	for tCase := range tCases {
		email, _ := EmailFromString(tCase)
		assert.Equal(t, tCases[tCase], email.Name())
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Domain
// ---------------------------------------------------------------------------------------------------------------------

func Test_Email_Domain(t *testing.T) {
	tCases := map[string]string{
		"te@te.te":                 "te.te",
		"t.te@te.te.te":            "te.te.te",
		"t+te@te.te":               "te.te",
		"t.e.s.t+t.e.s.t@te.te":    "te.te",
		"t.e.s.t+t.e.s.t@te.te.te": "te.te.te",
		"123456@t.test.te.te":      "t.test.te.te",
		"123456@test.t.te.te":      "test.t.te.te",
	}

	for tCase := range tCases {
		email, _ := EmailFromString(tCase)
		assert.Equal(t, tCases[tCase], email.Domain())
	}
}
