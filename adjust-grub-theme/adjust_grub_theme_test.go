package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_round(t *testing.T) {
	tests := []struct {
		Input    float64
		Expected int
	}{
		{
			1.1,
			1,
		},
		{
			1.6,
			2,
		},
		{
			2.1,
			2,
		},
	}
	for _, data := range tests {
		assert.Equal(t, data.Expected, round(data.Input))
	}
}
