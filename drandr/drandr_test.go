package drandr

import (
	"testing"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/stretchr/testify/assert"
)

func Test_GetScreenInfo(t *testing.T) {
	xConn, err := x.NewConn()
	if err != nil {
		t.Skip(err)
	}

	t.Run("Test_GetScreenInfo", func(t *testing.T) {
		_, err := GetScreenInfo(xConn)
		assert.NoError(t, err)
	})
}

func Test_GetPrimary(t *testing.T) {
	xConn, err := x.NewConn()
	if err != nil {
		t.Skip(err)
	}
	screenInfo, err := GetScreenInfo(xConn)
	if err != nil {
		t.Skip(err)
	}

	_, err = screenInfo.GetPrimary()
	assert.NoError(t, err)
}
