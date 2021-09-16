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
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"pkg.deepin.io/lib/encoding/kv"
	"pkg.deepin.io/lib/imgutil"
)

// copy from go source
func round(f float64) int {
	i := int(f)
	if f-float64(i) >= 0.5 {
		i += 1
	}
	return i
}

func convertSvg(svgFile, outFile string, width, height int) error {
	// #nosec G204
	cmd := exec.Command("rsvg-convert", "-o", outFile,
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		"-f", "png",
		svgFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Debugf("$ rsvg-convert --output %s -width %d -height %d -f png %s",
		outFile, width, height, svgFile)
	return cmd.Run()
}

func loadImage(filename string) (image.Image, error) {
	return imgutil.Load(filename)
}

func savePng(img image.Image, filename string) error {
	var enc png.Encoder
	enc.CompressionLevel = png.NoCompression
	return imgutil.SavePng(img, filename, &enc)
}

func saveJpeg(img image.Image, filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = fh.Close()
	}()
	bw := bufio.NewWriter(fh)
	err = jpeg.Encode(bw, img, nil)
	if err != nil {
		return err
	}
	err = bw.Flush()
	return err
}

var findFontCache map[string]findFontResult

type findFontResult struct {
	fontFile  string
	faceIndex int
}

func findFont(pattern string) (fontFile string, faceIndex int, err error) {
	cache, ok := findFontCache[pattern]
	if ok {
		return cache.fontFile, cache.faceIndex, nil
	}
	// #nosec G204
	cmd := exec.Command("fc-match", "--format", "%{file}\n%{index}", pattern)
	out, err := cmd.Output()
	if err != nil {
		return
	}
	fields := bytes.SplitN(out, []byte{'\n'}, 2)
	if len(fields) != 2 {
		err = errors.New("findFont: len of fields is not 2")
		return
	}
	fontFile = string(fields[0])
	faceIndex, err = strconv.Atoi(string(fields[1]))
	if err != nil {
		return
	}

	if findFontCache == nil {
		findFontCache = make(map[string]findFontResult)
	}
	findFontCache[pattern] = findFontResult{
		fontFile:  fontFile,
		faceIndex: faceIndex,
	}
	return
}

func copyVars(vars map[string]float64) map[string]float64 {
	varsCopy := make(map[string]float64, len(vars))
	for key, value := range vars {
		varsCopy[key] = value
	}
	return varsCopy
}

func eval(vars map[string]float64, expr string) (float64, error) {
	bc := exec.Command("bc")
	var stdInBuf bytes.Buffer

	for key, value := range vars {
		fmt.Fprintf(&stdInBuf, "%s=%f\n", key, value)
	}

	stdInBuf.WriteString("scale=10\n")
	stdInBuf.WriteString(expr)
	stdInBuf.WriteByte('\n')
	bc.Stdin = &stdInBuf
	out, err := bc.Output()
	if err != nil {
		return 0, err
	}
	out = bytes.TrimSuffix(out, []byte{'\n'})
	v, err := strconv.ParseFloat(string(out), 64)
	return v, err
}

func decodeShellValue(in string) string {
	// #nosec G204
	output, err := exec.Command("/bin/sh", "-c", "echo -n "+in).Output()
	if err != nil {
		// fallback
		return strings.Trim(in, "\"")
	}
	return string(output)
}

const defaultGrubGfxMode = "auto"
const grubGfxMode = "GRUB_GFXMODE"
const grubParamsFile = "/etc/default/grub"

func getGfxMode(params map[string]string) (val string) {
	val = decodeShellValue(params[grubGfxMode])
	if val == "" {
		val = defaultGrubGfxMode
	}
	return
}

func loadGrubParams(grubParamsFilePath string) (map[string]string, error) {
	params := make(map[string]string)
	f, err := os.Open(grubParamsFilePath)
	if err != nil {
		return params, err
	}
	defer func() {
		_ = f.Close()
	}()

	r := kv.NewReader(f)
	r.TrimSpace = kv.TrimLeadingTailingSpace
	r.Comment = '#'
	for {
		pair, err := r.Read()
		if err != nil {
			break
		}
		if pair.Key == "" {
			continue
		}
		params[pair.Key] = pair.Value
	}

	return params, nil
}

type InvalidResolutionError struct {
	Resolution string
}

func (err InvalidResolutionError) Error() string {
	return fmt.Sprintf("invalid resolution %q", err.Resolution)
}

func parseResolution(v string) (w, h int, err error) {
	if v == "auto" || v == "" {
		err = InvalidResolutionError{v}
		return
	}

	arr := strings.Split(v, "x")
	if len(arr) != 2 {
		err = InvalidResolutionError{v}
		return
	}
	// parse width
	tmpw, err := strconv.ParseUint(arr[0], 10, 32)
	if err != nil {
		err = InvalidResolutionError{v}
		return
	}

	// parse height
	tmph, err := strconv.ParseUint(arr[1], 10, 32)
	if err != nil {
		err = InvalidResolutionError{v}
		return
	}

	w = int(tmpw)
	h = int(tmph)

	if w == 0 || h == 0 {
		err = InvalidResolutionError{v}
		return
	}

	return
}

func getCurrentLocale() (locale string) {
	for _, envVar := range []string{"DEEPIN_LASTORE_LANG", "LC_ALL", "LANG"} {
		locale = os.Getenv(envVar)
		if locale != "" {
			return locale
		}
	}

	return getDefaultLocale()
}

const (
	systemLocaleFile  = "/etc/default/locale"
	systemdLocaleFile = "/etc/locale.conf"
	defaultLocale     = "en_US.UTF-8"
)

func getDefaultLocale() (locale string) {
	files := [...]string{
		systemLocaleFile,
		systemdLocaleFile,
	}
	for _, file := range files {
		locale = getLocaleFromFile(file)
		if locale != "" {
			// get locale success
			break
		}
	}
	if locale == "" {
		return defaultLocale
	}
	return locale
}

func getLocaleFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func() {
		_ = f.Close()
	}()

	r := kv.NewReader(f)
	r.Delim = '='
	r.Comment = '#'
	r.TrimSpace = kv.TrimLeadingTailingSpace
	for {
		pair, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return ""
		}

		if pair.Key == "LANG" {
			return pair.Value
		}
	}
	return ""
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = destination.Close()
	}()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
