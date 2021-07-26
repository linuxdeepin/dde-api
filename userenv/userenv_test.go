package userenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashQuote(t *testing.T) {
	data := []struct {
		in       string
		expected string
		hasErr   bool
	}{
		{"`", "\"\\`\"", false}, // expected: "\`"
		{"$", `"\$"`, false},
		{"\\", `"\\"`, false},
		{"\"", `"\""`, false},
		{" ", `" "`, false},
		{"$abc", `"\$abc"`, false},
		{"a\nb", "", true},
	}

	for _, value := range data {
		result, err := bashQuote(value.in)
		assert.Equal(t, value.expected, result)
		if value.hasErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestBashUnquote(t *testing.T) {
	data := []struct {
		in       string
		expected string
		hasErr   bool
	}{
		// `
		{"\\`", "`", false},
		// $
		{"\\$", "$", false},
		// \
		{`\\`, `\`, false},
		// "
		{`\"`, `"`, false},
		{"abc", "abc", false},
		{"abc\\", "", true},
	}

	for _, value := range data {
		result, err := bashUnquote(value.in)
		assert.Equal(t, value.expected, result)
		if value.hasErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestLoadFromFile(t *testing.T) {
	m, err := LoadFromFile("testdata/t1.txt")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"ENV1": "abc",
		"ENV2": "abc def",
		"ENV3": "abc`def",
		"ENV4": "$abc",
		"ENV5": "\"abc",
		"ENV6": "\\abc",
		"ENV7": "'abc",
	}, m)
}
