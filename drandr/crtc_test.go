package drandr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseCrtcRotation(t *testing.T) {
	tests := []struct {
		origin   uint16
		rotation uint16
		reflect  uint16
	}{
		{
			1,
			1,
			0,
		},
		{
			8,
			8,
			0,
		},
		{
			16,
			1,
			16,
		},
		{
			48,
			1,
			48,
		},
	}
	for _, data := range tests {
		rotation, reflect := parseCrtcRotation(data.origin)
		assert.Equal(t, data.rotation, rotation)
		assert.Equal(t, data.reflect, reflect)
	}
}

func sliceEq(a, b []uint16) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func Test_getRotations(t *testing.T) {
	tests := []struct {
		origin   uint16
		expected []uint16
	}{
		{
			1 + 2 + 4 + 8 + 16,
			[]uint16{
				1, 2, 4, 8,
			},
		},
		{
			1 + 2 + 4 + 8,
			[]uint16{
				1, 2, 4, 8,
			},
		},
		{
			1 + 2 + 4,
			[]uint16{
				1, 2, 4,
			},
		},
	}
	for _, data := range tests {
		rotations := getRotations(data.origin)
		assert.True(t, sliceEq(rotations, data.expected))
	}
}

func Test_getReflects(t *testing.T) {
	tests := []struct {
		origin   uint16
		expected []uint16
	}{
		{
			16 + 32 + 64,
			[]uint16{
				0, 16, 32, 48,
			},
		},
		{
			16 + 32,
			[]uint16{
				0, 16, 32, 48,
			},
		},
		{
			1 + 2 + 4 + 16,
			[]uint16{
				0, 16,
			},
		},
		{
			1 + 2 + 4,
			[]uint16{
				0,
			},
		},
	}
	for _, data := range tests {
		rotations := getReflects(data.origin)
		assert.True(t, sliceEq(rotations, data.expected))
	}
}
