package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"pkg.deepin.io/dde/api/grub_theme/font"
	tt "pkg.deepin.io/dde/api/grub_theme/themetxt"
)

type su struct {
	suite.Suite
}

func (s *su) SetupTest() {
	fmt.Println("before SetupTest")
	optThemeOutputDir = "testdata"
	optThemeInputDir = "testdata"
	optScreenWidth = 720
	optScreenHeight = 480
}

func (s *su) TearDownTest() {
	fmt.Println("after TearDownTest")
	optThemeOutputDir = defaultThemeOutputDir
	optThemeInputDir = defaultThemeInputDir
	optScreenHeight = 0
	optScreenWidth = 0

}

func TestAdjustGrubTheme(t *testing.T) {
	suite.Run(t, new(su))
}

func (s *su) TestRound() {
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
		assert.Equal(s.T(), data.Expected, round(data.Input))
	}
}

func (s *su) TestLoadBackgroundImage() {
	img, err := loadBackgroundImage()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), img)
}

func (s *su) TestAdjustBackground() {
	themeOutputDir := "testdata"
	img, err := loadBackgroundImage()
	require.NoError(s.T(), err)
	image, err := adjustBackground(themeOutputDir, img)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), image)

	optScreenWidth = 0
	optScreenHeight = 0
	img, err = loadBackgroundImage()
	require.NoError(s.T(), err)
	image, err = adjustBackground(themeOutputDir, img)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), image)
	_ = os.RemoveAll(filepath.Join(themeOutputDir, "background.jpg"))
}

func (s *su) TestAdjustResourcesOsLogos() {
	_, err := exec.LookPath("rsvg-convert")
	if err != nil {
		s.T().Skip(err)
	}
	themeInputDir := "testdata"
	themeOutputDir := "testdata"
	err = adjustResourcesOsLogos(themeInputDir, themeOutputDir, 720, 480)
	assert.NoError(s.T(), err)
	_ = os.RemoveAll(filepath.Join(themeOutputDir, "icons"))

}

func (s *su) TestGetFontSize() {
	tests := []struct {
		Input    int
		Expected int
	}{
		{
			360,
			12,
		},
		{
			720,
			12,
		},
		{
			1080,
			17,
		},
		{
			1440,
			23,
		},
	}
	for i, test := range tests {
		s.T().Run("TestGetFontSize"+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.Expected, getFontSize(0, test.Input))
		})
	}
}

func (s *su) TestGetScreenSizeFromGrubParams() {
	grubParamsFilePath := "testdata/grub"
	require.FileExists(s.T(), grubParamsFilePath)
	w, h, err := getScreenSizeFromGrubParams(grubParamsFilePath)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1024, w)
	assert.Equal(s.T(), 768, h)
}

var items = []string{
	"nw",
	"n",
	"ne",
	"w",
	"c",
	"e",
	"sw",
	"s",
	"se",
}

func (s *su) TestCropSaveStyleBox() {
	img, err := loadBackgroundImage()

	require.NoError(s.T(), err)
	require.NotNil(s.T(), img)
	filenamePrefix := "testdata/crop"

	cropSaveStyleBox(img, filenamePrefix, 3)
	for _, name := range items {
		fileName := strings.Join([]string{filenamePrefix, "_", name, ".png"}, "")
		assert.FileExists(s.T(), fileName)
		_ = os.RemoveAll(fileName)
	}
}

func (s *su) TestGetFallbackDir() {
	fallbackDir := getFallbackDir()
	assert.Equal(s.T(), "testdata/deepin-fallback", fallbackDir)
}

func (s *su) TestSetBackground() {
	defer func() {
		_ = os.RemoveAll(filepath.Join("testdata/deepin", "background.jpg"))
		_ = os.RemoveAll(filepath.Join("testdata/deepin", "background_source"))
		_ = os.RemoveAll(filepath.Join("testdata/deepin-fallback", "background.jpg"))
	}()
	filenamePrefix := "menu"
	setBackground("testdata/deepin/background.origin.jpg")
	for _, name := range items {
		fileName := strings.Join([]string{filenamePrefix, "_", name, ".png"}, "")
		fileNamePath := filepath.Join("testdata/deepin", fileName)
		assert.FileExists(s.T(), fileNamePath)
		_ = os.RemoveAll(fileNamePath)
	}

}

func (s *su) TestAdjustThemeNormal() {
	_, err := exec.LookPath("grub-mkfont")
	if err != nil {
		s.T().Skip(err)
	}
	optThemeOutputDir = "testdata/tmp"
	defer func() {
		_ = os.RemoveAll(optThemeOutputDir)
	}()
	err = adjustThemeNormal()
	assert.Equal(s.T(), nil, err)

}

func (s *su) TestAdjustThemeFallback() {
	optThemeOutputDir = "testdata/tmp"
	err := adjustThemeFallback()
	assert.Equal(s.T(), nil, err)
	_ = os.RemoveAll(optThemeOutputDir)
}

func (s *su) TestCopyBgSource() {
	err := copyBgSource("testdata/deepin/background.origin.jpg")
	assert.NoError(s.T(), err)
	assert.FileExists(s.T(), "testdata/deepin/background_source")

	_ = os.RemoveAll("testdata/deepin/background_source")
}

func (s *su) TestCopyThemeFiles() {
	copyThemeFiles("testdata/deepin", "testdata")
	assert.FileExists(s.T(), "testdata/terminal_box_c.png")
	_ = os.RemoveAll("testdata/terminal_box_c.png")
}

func (s *su) TestGenPF2Font() {
	_, err := exec.LookPath("grub-mkfont")
	if err != nil {
		s.T().Skip(err)
	}
	fontBaseName := "tmpFront"
	faceIndex := 1
	size := 20
	frontName := fmt.Sprintf("ag-%s-%d-%d.pf2", fontBaseName, faceIndex, size)
	pf2Font, err := genPF2Font("testdata", fontBaseName, faceIndex, size)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pf2Font)
	assert.FileExists(s.T(), filepath.Join("testdata", frontName))
	_ = os.RemoveAll(filepath.Join("testdata", frontName))
	genFontCache = make(map[genFontCacheKey]*font.Face)
}

func (s *su) TestAdjustFont() {
	_, err := exec.LookPath("grub-mkfont")
	if err != nil {
		s.T().Skip(err)
	}
	theme, err := tt.ParseThemeFile("testdata/deepin/theme.txt.tpl")
	vars := map[string]float64{
		"std_font_size": float64(18),
	}
	require.NoError(s.T(), err)
	require.NotNil(s.T(), theme)
	tmpDir := "testdata/tmp"
	_ = os.Mkdir(tmpDir, 0755)
	_, err = adjustFont(tmpDir, theme.Components[0], "item_font", vars)

	assert.NoError(s.T(), err)
	_ = os.RemoveAll(tmpDir)
}

func (s *su) TestAdjustTerminalFont() {
	_, err := exec.LookPath("grub-mkfont")
	if err != nil {
		s.T().Skip(err)
	}
	theme, err := tt.ParseThemeFile("testdata/deepin/theme.txt.tpl")
	vars := map[string]float64{
		"std_font_size": float64(18),
	}
	require.NoError(s.T(), err)
	require.NotNil(s.T(), theme)
	tmpDir := "testdata/tmp"
	_ = os.Mkdir(tmpDir, 0755)
	err = adjustTerminalFont(tmpDir, theme, vars)

	assert.NoError(s.T(), err)
	_ = os.RemoveAll(tmpDir)
}

func (s *su) TestAdjustProp() {
	_, err := exec.LookPath("bc")
	if err != nil {
		s.T().Skip(err)
	}
	theme, err := tt.ParseThemeFile("testdata/deepin/theme.txt.tpl")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), theme)
	itemHeight := adjustProp(theme.Components[0], "item_height", map[string]float64{
		"font_height": float64(18),
	})
	assert.Equal(s.T(), 28, itemHeight)

	theme, err = tt.ParseThemeFile("testdata/deepin/theme.txt.tpl")
	assert.NoError(s.T(), err)
	itemHeight = adjustProp(theme.Components[0], "item_height", map[string]float64{
		"font_height": float64(5),
	})
	assert.Equal(s.T(), 8, itemHeight)

	theme, err = tt.ParseThemeFile("testdata/deepin/theme.txt.tpl")
	assert.NoError(s.T(), err)
	itemHeight = adjustProp(theme.Components[0], "item_height", map[string]float64{
		"font_height": float64(10),
	})
	assert.Equal(s.T(), 16, itemHeight)
}

func (s *su) TestGetCurrentLocale() {
	envVar := "DEEPIN_LASTORE_LANG"
	currentEnv := os.Getenv(envVar)
	err := os.Setenv(envVar, "zh_CN.UTF-8")
	require.NoError(s.T(), err)
	locale := getCurrentLocale()
	assert.Equal(s.T(), "zh_CN.UTF-8", locale)
	err = os.Setenv(envVar, currentEnv)
}

func (s *su) TestGetLocaleFromFile() {
	localeFilePath := "testdata/locale"
	locale := getLocaleFromFile(localeFilePath)
	assert.Equal(s.T(), "zh_CN.UTF-8", locale)

	locale = getLocaleFromFile("testdata/tmp")
	assert.Equal(s.T(), "", locale)
}

func (s *su) TestParseResolution() {
	tests := []struct {
		Input         string
		ExpectedErr   error
		ExpectedWidth int
		ExpectedHigh  int
	}{
		{
			"auto",
			InvalidResolutionError{"auto"},
			0,
			0,
		},
		{
			"1920x1080x1",
			InvalidResolutionError{"1920x1080x1"},
			0,
			0,
		},
		{
			"0x0",
			InvalidResolutionError{"0x0"},
			0,
			0,
		},
		{
			"1920x1080",
			nil,
			1920,
			1080,
		},
	}
	for i, test := range tests {
		s.T().Run("TestParseResolution"+strconv.Itoa(i), func(t *testing.T) {
			w, h, err := parseResolution(test.Input)
			assert.Equal(s.T(), test.ExpectedErr, err)
			assert.Equal(s.T(), test.ExpectedWidth, w)
			assert.Equal(s.T(), test.ExpectedWidth, w)
			assert.Equal(s.T(), test.ExpectedHigh, h)
		})

	}
}

func (s *su) TestError() {
	err := InvalidResolutionError{"This is test error"}
	assert.Equal(s.T(), fmt.Sprintf("invalid resolution %q", "This is test error"), err.Error())
}
