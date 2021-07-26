package font

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadFont(t *testing.T) {
	font, err := LoadFont("testdata/unifont-regular-16.pf2")
	assert.NoError(t, err)
	assert.NotNil(t, font)
}

func Test_getString(t *testing.T) {
	s := &section{
		name: "section1",
		data: []byte("457"),
	}
	str := s.getString()
	assert.Equal(t, "457", str)
}

func Test_getUint16BE(t *testing.T) {
	s := &section{
		name: "section1",
		data: []byte("457"),
	}
	assert.Equal(t, 13365, int(s.getUint16BE()))
}

func Test_Close(t *testing.T) {
	font, err := LoadFont("testdata/unifont-regular-16.pf2")
	require.NoError(t, err)
	assert.NotNil(t, font)
	assert.Nil(t, font.Close())
}

func Test_Height(t *testing.T) {
	font, err := LoadFont("testdata/unifont-regular-16.pf2")
	require.NoError(t, err)
	assert.NotNil(t, font)
	assert.Equal(t, 16, font.Height())
}
