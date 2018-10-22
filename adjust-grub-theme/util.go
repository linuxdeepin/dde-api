package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"pkg.deepin.io/lib/encoding/kv"
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
	cmd := exec.Command("rsvg-convert", "-o", outFile,
		"-w", strconv.Itoa(width),
		"-h", strconv.Itoa(height),
		"-f", "png",
		svgFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("$ rsvg-convert --output %s -width %d -height %d -f png %s\n",
		outFile, width, height, svgFile)
	return cmd.Run()
}

func loadImage(filename string) (image.Image, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	br := bufio.NewReader(fh)
	img, _, err := image.Decode(br)
	return img, err
}

func savePng(img image.Image, filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	bw := bufio.NewWriter(fh)

	var enc png.Encoder
	enc.CompressionLevel = png.NoCompression
	err = enc.Encode(bw, img)
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

func loadGrubParams() (map[string]string, error) {
	params := make(map[string]string)
	f, err := os.Open(grubParamsFile)
	if err != nil {
		return params, err
	}
	defer f.Close()

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

func getCurrentLocale() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en"
	}
	return lang
}

func loadThemeHeadInfo(filename string) (map[string]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	headInfo := make(map[string]string)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.Equal(line, []byte("#head end")) ||
			!bytes.HasPrefix(line, []byte{'#'}) {
			break
		}
		fields := bytes.SplitN(line, []byte{':'}, 2)
		if len(fields) != 2 {
			continue
		}

		key := string(fields[0])
		value := string(fields[1])
		headInfo[key] = value
	}
	return headInfo, nil
}
