package themetxt

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseThemeFile(t *testing.T) {
	theme, err := ParseThemeFile("testdata/theme.txt.tpl")
	require.NoError(t, err)
	require.NotNil(t, theme)
	assert.Equal(t, 10, len(theme.Props))
	assert.Equal(t, 3, len(theme.Components))
}

func Test_GetProp(t *testing.T) {
	theme, err := ParseThemeFile("testdata/theme.txt.tpl")
	require.NoError(t, err)
	require.NotNil(t, theme)
	var component *Component
	for _, c := range theme.Components {
		if c.Type == "boot_menu" {
			component = c
		}
	}
	require.NotNil(t, component)
	prop, b := component.GetProp("left")
	assert.True(t, b)
	assert.Equal(t, "15%", strconv.Itoa(int(prop.(RelNum)))+"%")

	prop, b = component.GetProp("item_font")
	assert.True(t, b)
	assert.Equal(t, "Unifont Regular 16", prop.(string))

	prop, b = component.GetProp("menu_pixmap_style")
	assert.True(t, b)
	assert.Equal(t, "menu_*.png", prop.(string))
}

func Test_GetPropString(t *testing.T) {
	theme, err := ParseThemeFile("testdata/theme.txt.tpl")
	require.NoError(t, err)
	require.NotNil(t, theme)
	var component *Component
	for _, c := range theme.Components {
		if c.Type == "boot_menu" {
			component = c
		}
	}
	require.NotNil(t, component)
	prop, b := component.GetPropString("item_font")
	assert.True(t, b)
	assert.Equal(t, "Unifont Regular 16", prop)

}

func Test_GetPropInt(t *testing.T) {
	theme, err := ParseThemeFile("testdata/theme.txt.tpl")
	require.NoError(t, err)
	require.NotNil(t, theme)
	var component *Component
	for _, c := range theme.Components {
		if c.Type == "boot_menu" {
			component = c
		}
	}
	require.NotNil(t, component)
	prop, b := component.GetPropInt("item_height")
	assert.True(t, b)
	assert.Equal(t, 24, prop)
}

func Test_SetProp(t *testing.T) {
	theme, err := ParseThemeFile("testdata/theme.txt.tpl")
	require.NoError(t, err)
	require.NotNil(t, theme)
	theme.SetProp("NewProp", "utProp")
	assert.Equal(t, 11, len(theme.Props))
	theme.SetProp("terminal-height", "98%")
	propString, b := theme.GetPropString("terminal-height")
	assert.True(t, b)
	assert.Equal(t, "98%", propString)
}
