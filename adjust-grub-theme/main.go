/*
 * Copyright (C) 2016 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"

	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/log"

	"github.com/disintegration/imaging"

	"pkg.deepin.io/dde/api/grub_theme/font"
	tt "pkg.deepin.io/dde/api/grub_theme/themetxt"
	"pkg.deepin.io/lib/locale"
)

const (
	defaultThemeOutputDir = "/boot/grub/themes/"
	defaultThemeInputDir  = "/usr/share/dde-api/data/grub-themes/"

	themeNameNormal   = "deepin"
	themeNameFallback = "deepin-fallback"
)

var optScreenHeight int
var optScreenWidth int
var optThemeInputDir string
var optThemeOutputDir string
var optLang string
var optVersion bool
var optSetBackground string
var optLogSys bool
var optTerminalFontSize int
var optTerminalFontName string
var optFallbackOnly bool
var logger *log.Logger

func init() {
	logger = log.NewLogger("adjust-grub-theme")

	flag.IntVar(&optScreenWidth, "width", 0, "screen width")
	flag.IntVar(&optScreenHeight, "height", 0, "screen height")
	flag.StringVar(&optThemeInputDir, "theme-input", defaultThemeInputDir,
		"theme input directory")
	flag.StringVar(&optThemeOutputDir, "theme-output", defaultThemeOutputDir,
		"theme output directory")
	flag.StringVar(&optLang, "lang", "", "language")
	flag.BoolVar(&optVersion, "version", false, "show version")
	flag.StringVar(&optSetBackground, "set-background", "", "")
	flag.BoolVar(&optLogSys, "log-sys", false, "")
	flag.IntVar(&optTerminalFontSize, "tf-size", -1, "terminal font size")
	flag.StringVar(&optTerminalFontName, "tf-name", "Unifont:style=Medium",
		"terminal font name")
	flag.BoolVar(&optFallbackOnly, "fallback-only", false, "fallback only")
}

func loadBackgroundImage() (image.Image, error) {
	img, err := loadImage(filepath.Join(optThemeOutputDir, themeNameNormal, "background_source"))
	if err != nil {
		logger.Warning("failed to load image background_source:", err)
		originDesktopImageFile := filepath.Join(optThemeInputDir, themeNameNormal, "background.origin.jpg")
		img, err = loadImage(originDesktopImageFile)
		if err != nil {
			logger.Warning(err)
			return nil, err
		}
	}
	return img, nil
}

func adjustBackground(themeOutputDir string, img image.Image) (image.Image, error) {
	logger.Debug("adjustBackground")
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	x, y, w, h, err := graphic.GetPreferScaleClipRect(optScreenWidth, optScreenHeight, imgWidth, imgHeight)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}
	img0 := imaging.Crop(img, image.Rect(x, y, x+w, y+h))
	img0 = imaging.Resize(img0, optScreenWidth, optScreenHeight, imaging.Lanczos)
	// save img
	err = saveJpeg(img0, filepath.Join(themeOutputDir, "background.jpg"))
	if err != nil {
		return nil, err
	}
	return img0, nil
}

func adjustResourcesOsLogos(themeInputDir, themeOutputDir string, width, height int) error {
	srcDir := filepath.Join(themeInputDir, "resources/os-logos")
	fileInfoList, err := ioutil.ReadDir(srcDir)
	if err != nil {
		logger.Warning(err)
		return err
	}

	outDir := filepath.Join(themeOutputDir, "icons")
	err = os.Mkdir(outDir, 0755)
	if err != nil {
		logger.Warning(err)
		return err
	}

	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			continue
		}

		file := filepath.Join(srcDir, fileInfo.Name())
		ext := filepath.Ext(fileInfo.Name())
		if ext != ".svg" {
			continue
		}
		outFileName := strings.TrimSuffix(fileInfo.Name(), ext) + ".png"
		outFile := filepath.Join(outDir, outFileName)
		err = convertSvg(file, outFile, width, height)
		if err != nil {
			logger.Warning(err)
			return err
		}
	}
	return nil
}

const (
	minFontSize         = 12
	minTerminalFontSize = 14
)

// min 12px
func getFontSize(screenWidth int, screenHeight int) int {
	var x1 float64 = 768
	var y1 float64 = minFontSize
	var x2 float64 = 2000
	var y2 float64 = 32
	y := (float64(screenHeight)-x1)/(x2-x1)*(y2-y1) + y1

	if y < minFontSize {
		y = minFontSize
	}

	return round(y)
}

func getScreenSizeFromGrubParams(grubParamsFilePath string) (w, h int, err error) {
	params, err := loadGrubParams(grubParamsFilePath)
	if err != nil {
		return
	}

	w, h, err = parseResolution(getGfxMode(params))
	if err != nil {
		return
	}
	return
}

func cropSaveStyleBox(img image.Image, filenamePrefix string, r int) {
	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()

	// center width
	cw := imgW - 2*r
	// center height
	ch := imgH - 2*r

	items := []struct {
		rect image.Rectangle
		name string
	}{
		{
			rect: image.Rect(0, 0, r, r),
			name: "nw",
		},
		{
			rect: image.Rect(r, 0, r+cw, r),
			name: "n",
		},
		{
			rect: image.Rect(r+cw, 0, imgW, r),
			name: "ne",
		},
		{
			rect: image.Rect(0, r, r, r+ch),
			name: "w",
		},
		{
			rect: image.Rect(r, r, r+cw, r+ch),
			name: "c",
		},
		{
			rect: image.Rect(r+cw, r, imgW, r+ch),
			name: "e",
		},

		{
			rect: image.Rect(0, r+ch, r, imgH),
			name: "sw",
		},
		{
			rect: image.Rect(r, r+ch, r+cw, imgH),
			name: "s",
		},
		{
			rect: image.Rect(r+cw, r+ch, imgW, imgH),
			name: "se",
		},
	}

	for _, item := range items {
		img0 := imaging.Crop(img, item.rect)
		err := savePng(img0, filenamePrefix+"_"+item.name+".png")
		if err != nil {
			logger.Warning(err)
		}
	}
}

func getFallbackDir() string {
	return filepath.Clean(filepath.Join(optThemeOutputDir, themeNameFallback))
}

func setBackground(bgFile string) {
	err := copyBgSource(bgFile)
	if err != nil {
		logger.Fatal(err)
	}

	bgImg, err := loadBackgroundImage()
	if err != nil {
		logger.Fatal(err)
	}

	fallbackDir := getFallbackDir()
	_, err = os.Stat(fallbackDir)
	if err == nil {
		err = saveJpeg(bgImg, filepath.Join(fallbackDir, "background.jpg"))
		if err != nil {
			logger.Fatal(err)
		}
	}

	themeOutputDir := filepath.Join(optThemeOutputDir, themeNameNormal)
	bgImg, err = adjustBackground(themeOutputDir, bgImg)
	if err != nil {
		logger.Fatal(err)
	}

	themeTxtFile := filepath.Join(themeOutputDir, "theme.txt")
	theme, err := tt.ParseThemeFile(themeTxtFile)
	if err != nil {
		logger.Warning(err)
		return
	}

	var bmComp *tt.Component
	for _, comp := range theme.Components {
		if comp.Type == tt.ComponentTypeBootMenu {
			bmComp = comp
			break
		}
	}
	if bmComp == nil {
		logger.Warning("not found boot_menu component")
		return
	}

	convertPropRel2Abs(bmComp, "left", orientationHorizontal)
	convertPropRel2Abs(bmComp, "top", orientationVertical)
	adjustBootMenuPixmapStyle(themeOutputDir, bmComp, bgImg)
}

func adjustTheme() {
	err := adjustThemeFallback()
	if err != nil {
		logger.Fatal(err)
	}
	if optFallbackOnly {
		return
	}
	err = adjustThemeNormal()
	if err != nil {
		logger.Fatal(err)
	}
}

func adjustThemeNormal() error {
	themeInputDir := filepath.Join(optThemeInputDir, themeNameNormal)
	themeOutputDir := filepath.Join(optThemeOutputDir, themeNameNormal)

	vars := map[string]float64{}
	themeFile := filepath.Join(themeInputDir, "theme.txt.tpl")
	theme, err := tt.ParseThemeFile(themeFile)
	if err != nil {
		return err
	}

	cleanupThemeOutputDir(themeOutputDir)
	err = os.MkdirAll(themeOutputDir, 0755)
	if err != nil {
		return err
	}
	copyThemeFiles(themeInputDir, themeOutputDir)

	stdFontSize := getFontSize(optScreenWidth, optScreenHeight)
	vars["std_font_size"] = float64(stdFontSize)
	vars["screen_width"] = float64(optScreenWidth)
	vars["screen_height"] = float64(optScreenHeight)

	err = adjustTerminalFont(themeOutputDir, theme, vars)
	if err != nil {
		logger.Fatal(err)
	}

	for _, comp := range theme.Components {
		if comp.Type == tt.ComponentTypeBootMenu {
			adjustBootMenu(themeOutputDir, comp, vars)

			iconWidth, _ := comp.GetPropInt("icon_width")
			iconHeight, _ := comp.GetPropInt("icon_height")
			err = adjustResourcesOsLogos(themeInputDir, themeOutputDir, iconWidth, iconHeight)
			if err != nil {
				logger.Warning(err)
			}

		} else if comp.Type == tt.ComponentTypeLabel {
			adjustLabel(themeOutputDir, comp, vars)
		}
	}

	themeOutput := filepath.Join(themeOutputDir, "theme.txt")
	themeOutputFh, err := os.Create(themeOutput)
	if err != nil {
		return err
	}
	defer func() {
		_ = themeOutputFh.Close()
	}()
	bw := bufio.NewWriter(themeOutputFh)
	// write head
	_, err = fmt.Fprintf(bw, "#version:%d\n", VERSION)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(bw, "#lang:%s\n", optLang)
	if err != nil {
		return err
	}

	var inputDir string
	inputDir, err = filepath.Abs(themeInputDir)
	if err != nil {
		logger.Warning(err)
		inputDir = themeInputDir
	}

	_, err = fmt.Fprintf(bw, "#themeInputDir:%s\n", inputDir)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(bw, "#screenWidth:%d\n", optScreenWidth)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(bw, "#screenHeight:%d\n", optScreenHeight)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(bw, "#head end")
	if err != nil {
		return err
	}

	_, err = theme.WriteTo(bw)
	if err != nil {
		return err
	}
	err = bw.Flush()
	return err
}

func adjustThemeFallback() error {
	themeInputDir := filepath.Join(optThemeInputDir, themeNameFallback)
	themeOutputDir := filepath.Join(optThemeOutputDir, themeNameFallback)

	cleanupThemeOutputDir(themeOutputDir)
	err := os.MkdirAll(themeOutputDir, 0755)
	if err != nil {
		return err
	}
	copyThemeFiles(themeInputDir, themeOutputDir)

	bgImg, err := loadBackgroundImage()
	if err != nil {
		return err
	}

	err = saveJpeg(bgImg, filepath.Join(themeOutputDir, "background.jpg"))
	if err != nil {
		return err
	}

	themeFile := filepath.Join(themeInputDir, "theme.txt.tpl")
	theme, err := tt.ParseThemeFile(themeFile)
	if err != nil {
		return err
	}

	for _, comp := range theme.Components {
		if comp.Type == tt.ComponentTypeLabel {
			adjustLabelText(comp)
		}
	}

	themeOutput := filepath.Join(themeOutputDir, "theme.txt")
	themeOutputFh, err := os.Create(themeOutput)
	if err != nil {
		return err
	}
	defer func() {
		_ = themeOutputFh.Close()
	}()
	bw := bufio.NewWriter(themeOutputFh)
	// write head
	_, err = fmt.Fprintf(bw, "#version:%d\n", VERSION)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(bw, "#lang:%s\n", optLang)
	if err != nil {
		return err
	}

	var inputDir string
	inputDir, err = filepath.Abs(themeInputDir)
	if err != nil {
		logger.Warning(err)
		inputDir = themeInputDir
	}

	_, err = fmt.Fprintf(bw, "#themeInputDir:%s\n", inputDir)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(bw, "#head end")
	if err != nil {
		return err
	}

	_, err = theme.WriteTo(bw)
	if err != nil {
		return err
	}
	err = bw.Flush()
	return err
}

func main() {
	flag.Parse()
	if optVersion {
		fmt.Printf("%d\n", VERSION)
		return
	}

	if optLogSys {
		logger.RemoveBackendConsole()
	}

	if optScreenWidth == 0 || optScreenHeight == 0 {
		var err error
		optScreenWidth, optScreenHeight, err = getScreenSizeFromGrubParams(grubParamsFile)
		if err != nil {
			logger.Debug(err)
			optScreenWidth = 1024
			optScreenHeight = 768
		}
		logger.Debug("screen width:", optScreenWidth)
		logger.Debug("screen height:", optScreenHeight)
	}

	if optSetBackground != "" {
		setBackground(optSetBackground)
		return
	}

	if optLang == "" {
		optLang = getCurrentLocale()
	}
	logger.Debug("lang:", optLang)

	adjustTheme()
}

func copyBgSource(filename string) error {
	dstFile := filepath.Join(optThemeOutputDir, themeNameNormal, "background_source")
	err := os.Remove(dstFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(filepath.Dir(dstFile), 0755)
	if err != nil {
		return err
	}
	_, err = copyFile(filename, dstFile)
	return err
}

func copyThemeFiles(themeInputDir, themeOutputDir string) {
	fileInfoList, err := ioutil.ReadDir(themeInputDir)
	if err != nil {
		logger.Warning(err)
		return
	}
	for _, fileInfo := range fileInfoList {
		name := fileInfo.Name()
		ext := filepath.Ext(name)
		if ext == ".png" || ext == ".pf2" {
			srcFile := filepath.Join(themeInputDir, name)
			dstFile := filepath.Join(themeOutputDir, name)
			logger.Debug("copyFile", srcFile, dstFile)
			_, err := copyFile(srcFile, dstFile)
			if err != nil {
				logger.Warning("failed to copy file:", err)
			}
		}
	}
}

func cleanupThemeOutputDir(dir string) {
	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Warning(err)
		}
	}

	for _, fileInfo := range fileInfoList {
		filename := filepath.Join(dir, fileInfo.Name())
		if fileInfo.IsDir() {
			_ = os.RemoveAll(filename)
		} else {
			if fileInfo.Name() == "background_source" {
				// do not delete it
			} else {
				_ = os.Remove(filename)
			}
		}
	}
}

var genFontCache map[genFontCacheKey]*font.Face

type genFontCacheKey struct {
	fontFile  string
	faceIndex int
	size      int
}

func genPF2Font(themeOutputDir, fontFile string, faceIndex, size int) (*font.Face, error) {
	cacheKey := genFontCacheKey{
		fontFile:  fontFile,
		faceIndex: faceIndex,
		size:      size,
	}
	face, ok := genFontCache[cacheKey]
	if ok {
		logger.Debug("genPF2Font use cache")
		return face, nil
	}

	sizeStr := strconv.Itoa(size)
	faceIndexStr := strconv.Itoa(faceIndex)

	fontBaseName := filepath.Base(fontFile)
	// trim ext
	fontBaseName = strings.TrimSuffix(fontBaseName, filepath.Ext(fontBaseName))
	fontBaseName = fmt.Sprintf("ag-%s-%d-%d.pf2", fontBaseName, faceIndex, size)
	output := filepath.Join(themeOutputDir, fontBaseName)

	// #nosec G204
	cmd := exec.Command("grub-mkfont", "-i", faceIndexStr,
		"-s", sizeStr, "-o", output, fontFile)
	logger.Debugf("$ grub-mkfont -i %d -s %d -o %s %s",
		faceIndex, size, output, fontFile)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	face, err = font.LoadFont(output)
	if err != nil {
		return nil, err
	}

	if genFontCache == nil {
		genFontCache = make(map[genFontCacheKey]*font.Face)
	}
	genFontCache[cacheKey] = face
	return face, nil
}

func parseTplFont(str string) (fontName string, sizeScale float64, err error) {
	fields := strings.SplitN(str, ";", 2)
	if len(fields) != 2 {
		return "", 0, errors.New("invalid font format")
	}
	fontName = fields[0]
	sizeScale, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return "", 0, err
	}
	return fontName, sizeScale, nil
}

func adjustFont(themeOutputDir string, comp *tt.Component, propName string,
	vars map[string]float64) (*font.Face, error) {

	propFont, _ := comp.GetPropString(propName)
	fontName, sizeScale, err := parseTplFont(propFont)
	if err != nil {
		return nil, err
	}
	logger.Debugf("adjustFont fontName: %q, sizeScale: %g", fontName, sizeScale)

	fontFile, faceIndex, err := findFont(fontName)
	if err != nil {
		return nil, err
	}

	fontSize := round(vars["std_font_size"] * sizeScale)
	if fontSize < minFontSize {
		fontSize = minFontSize
	}

	face, err := genPF2Font(themeOutputDir, fontFile, faceIndex, fontSize)
	if err != nil {
		return nil, err
	}

	comp.SetProp(propName, face.Name)
	return face, nil
}

func adjustTerminalFont(themeOutputDir string, theme *tt.Theme, vars map[string]float64) error {
	var fontName string
	var sizeScale float64
	var err error

	const propName = "terminal-font"

	if optTerminalFontSize > 0 {
		fontName = optTerminalFontName
	} else {
		propFont, _ := theme.GetPropString(propName)
		fontName, sizeScale, err = parseTplFont(propFont)
		if err != nil {
			return err
		}
	}

	fontFile, faceIndex, err := findFont(fontName)
	if err != nil {
		return err
	}

	var fontSize int
	if optTerminalFontSize > 0 {
		fontSize = optTerminalFontSize
	} else {
		fontSize = round(vars["std_font_size"] * sizeScale)
		if fontSize < minTerminalFontSize {
			fontSize = minTerminalFontSize
		}

		// NOTE: grub gfxterm 的终端在使用 unifont 字体时，字体大小为 16 时，会出现光标残留问题。
		if fontSize == 16 {
			fontSize = 17
		}
	}

	face, err := genPF2Font(themeOutputDir, fontFile, faceIndex, fontSize)
	if err != nil {
		return err
	}

	theme.SetProp(propName, face.Name)
	return nil
}

func adjustProp(comp *tt.Component, propName string, vars map[string]float64) int {
	propVal, ok := comp.GetProp(propName)
	if !ok {
		return 0
	}
	propValStr, ok := propVal.(string)
	if !ok {
		return 0
	}
	evalResult, err := eval(vars, propValStr)
	if err != nil {
		logger.Fatal(err)
	}
	evalRet := round(evalResult)
	if evalRet < 0 {
		evalRet = 0
	}
	vars[propName] = float64(evalRet)
	comp.SetProp(propName, evalRet)
	logger.Debug("adjustProp", propName, evalRet)
	return evalRet
}

func adjustSelectedItemPixmapStyle(r int) {
	width := 2*r + 1
	dc := gg.NewContext(width, width)
	dc.SetRGBA(1, 1, 1, 0.2)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(width), float64(r))
	dc.Fill()
	prefix := filepath.Join(optThemeOutputDir, themeNameNormal, "selected_item")
	cropSaveStyleBox(dc.Image(), prefix, r)
}

func adjustItemPixmapStyle(r int) {
	width := 2*r + 1
	img := image.NewAlpha(image.Rect(0, 0, width, width))
	prefix := filepath.Join(optThemeOutputDir, themeNameNormal, "item")
	cropSaveStyleBox(img, prefix, r)
}

func adjustScrollbarThumbPixmapStyle(r int) {
	width := 2 * r
	height := 2*r + 1
	dc := gg.NewContext(width, height)
	dc.SetRGBA(1, 1, 1, 0.2)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), float64(r))
	dc.Fill()
	cropSaveScrollbarThumbStyleBox(dc.Image(), r, "scrollbar_thumb")
}

func cropSaveScrollbarThumbStyleBox(img image.Image, r int, name string) {
	w := 2 * r
	imgN := imaging.CropAnchor(img, w, r, imaging.Top)
	imgS := imaging.CropAnchor(img, w, r, imaging.Bottom)
	imgC := imaging.CropAnchor(img, w, 1, imaging.Center)
	namePrefix := filepath.Join(optThemeOutputDir, themeNameNormal, name)
	err := imaging.Save(imgN, namePrefix+"_n.png")
	if err != nil {
		logger.Warning(err)
		return
	}
	err = imaging.Save(imgS, namePrefix+"_s.png")
	if err != nil {
		logger.Warning(err)
		return
	}
	err = imaging.Save(imgC, namePrefix+"_c.png")
	if err != nil {
		logger.Warning(err)
		return
	}
}

func getBootMenuR(itemHeight int) int {
	return round(float64(itemHeight) * 0.5)
}

func getItemR(itemHeight int) int {
	return round(float64(itemHeight) * 0.16)
}

func getScrollbarThumbR(menuR int) int {
	return round(float64(menuR) * 0.167)
}

func adjustBootMenuPixmapStyle(themeOutputDir string, comp *tt.Component, bgImg image.Image) {
	logger.Debug("adjustBootMenuPixmapStyle")
	itemHeight, _ := comp.GetPropInt("item_height")
	bmLeft, _ := comp.GetPropInt("left")
	bmTop, _ := comp.GetPropInt("top")
	bmWidth, _ := comp.GetPropInt("width")
	bmHeight, _ := comp.GetPropInt("height")

	menuR := getBootMenuR(itemHeight)
	logger.Debug("menuR:", menuR)
	// boot menu rect
	rect := image.Rect(bmLeft, bmTop,
		bmLeft+bmWidth, bmTop+bmHeight)

	x := menuR * 2
	y := x
	w := bmWidth - x*2
	h := bmHeight - y*2

	shadowDc := gg.NewContext(rect.Dx(), rect.Dy())
	shadowDc.SetRGBA(0, 0, 0, 0.2) // black
	shadowYOffset := menuR
	shadowDc.DrawRoundedRectangle(float64(x), float64(y+shadowYOffset), float64(w), float64(h), float64(menuR))
	shadowDc.Fill()
	// shadow blur sigma : 10
	shadowImg := imaging.Blur(shadowDc.Image(), 10)

	img1 := imaging.Crop(bgImg, rect)

	img1b := imaging.Blur(img1, 15)
	imgWhite := imaging.New(bmWidth, bmHeight, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	img1b = imaging.Overlay(img1b, imgWhite, image.Pt(0, 0), 0.1)
	dc := gg.NewContext(bmWidth, bmHeight)
	dc.DrawRoundedRectangle(float64(x), float64(y), float64(w), float64(h),
		float64(menuR))
	dc.Clip()
	dc.DrawImage(img1b, 0, 0)
	// img2 是模糊过的圆角的
	img2 := dc.Image()

	img3 := imaging.Overlay(shadowImg, img2, image.Pt(0, 0), 1)

	prefix := filepath.Join(themeOutputDir, "menu")
	cropSaveStyleBox(img3, prefix, x+menuR)
}

func adjustBootMenu(themeOutputDir string, comp *tt.Component, vars map[string]float64) {
	vars = copyVars(vars)
	face, err := adjustFont(themeOutputDir, comp, "item_font", vars)
	if err != nil {
		logger.Fatal(err)
	}

	fontHeight := face.Height()
	vars["font_height"] = float64(fontHeight)

	itemHeight := adjustProp(comp, "item_height", vars)
	itemR := getItemR(itemHeight)
	vars["item_r"] = float64(itemR)
	menuR := getBootMenuR(itemHeight)
	vars["menu_r"] = float64(menuR)

	for _, propName := range []string{
		"item_spacing",
		"item_padding", "icon_width",
		"icon_height", "item_icon_space",
		"height", "width",
		"left", "top",
	} {
		adjustProp(comp, propName, vars)
	}

	scrollbarThumbR := getScrollbarThumbR(menuR)
	comp.SetProp("scrollbar_width", scrollbarThumbR*2)

	bgImg, err := loadBackgroundImage()
	if err != nil {
		logger.Fatal(err)
	}

	bgImg, err = adjustBackground(themeOutputDir, bgImg)
	if err != nil {
		logger.Fatal(err)
	}
	adjustBootMenuPixmapStyle(themeOutputDir, comp, bgImg)

	convertPropAbs2Rel(comp, "left", orientationHorizontal)
	convertPropAbs2Rel(comp, "top", orientationVertical)

	adjustSelectedItemPixmapStyle(itemR)
	adjustItemPixmapStyle(itemR)
	adjustScrollbarThumbPixmapStyle(scrollbarThumbR)
}

const (
	orientationHorizontal = 0
	orientationVertical   = 1
)

func convertPropAbs2Rel(comp *tt.Component, propName string, orientation int) {
	var ref int
	switch orientation {
	case orientationHorizontal:
		ref = optScreenWidth
	case orientationVertical:
		ref = optScreenHeight
	default:
		panic("invalid orientation")
	}

	abs, _ := comp.GetPropInt(propName)
	rel := tt.RelNum(round(float64(abs) / float64(ref) * 100.0))
	comp.SetProp(propName, rel)
}

func convertPropRel2Abs(comp *tt.Component, propName string, orientation int) {
	var ref int
	switch orientation {
	case orientationHorizontal:
		ref = optScreenWidth
	case orientationVertical:
		ref = optScreenHeight
	default:
		panic("invalid orientation")
	}

	p, _ := comp.GetProp(propName)

	switch pp := p.(type) {
	case tt.AbsNum:
		return
	case int:
		return
	case tt.RelNum:
		abs := round(float64(pp) / 100.0 * float64(ref))
		comp.SetProp(propName, abs)
	}
}

func adjustLabel(themeOutputDir string, comp *tt.Component, vars map[string]float64) {
	vars = copyVars(vars)
	face, err := adjustFont(themeOutputDir, comp, "font", vars)
	if err != nil {
		logger.Fatal(err)
	}

	fontHeight := face.Height()
	vars["font_height"] = float64(fontHeight)

	for _, propName := range []string{"left", "top", "width", "height"} {
		adjustProp(comp, propName, vars)
	}
	convertPropAbs2Rel(comp, "top", orientationVertical)

	adjustLabelText(comp)
}

func adjustLabelText(comp *tt.Component) {
	localeVars := locale.GetLocaleVariants(optLang)
	var text string
	var textFound bool
	for _, localeVar := range localeVars {
		var ok bool
		text, ok = comp.GetPropString("_text_" + localeVar)
		if ok {
			textFound = true
			break
		}
	}
	if !textFound {
		text, _ = comp.GetPropString("_text_en")
	}
	comp.SetProp("text", text)
}
