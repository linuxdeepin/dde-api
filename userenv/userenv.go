package userenv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"pkg.deepin.io/lib/xdg/basedir"
)

var defaultFile = filepath.Join(basedir.GetUserHomeDir(), ".dde_env")

func DefaultFile() string {
	return defaultFile
}

func isSpecialChar(char byte) bool {
	switch char {
	case '`', '$', '\\', '"':
		return true
	default:
		return false
	}
}

// ex. $ => "\$"
func bashQuote(str string) (string, error) {
	var buf bytes.Buffer
	buf.WriteByte('"')
	r := strings.NewReader(str)
	for {
		char, err := r.ReadByte()
		if err != nil {
			break
		}
		if isSpecialChar(char) {
			buf.WriteByte('\\')
		} else if char == '\n' {
			return "", errors.New("invalid char newline")
		}
		buf.WriteByte(char)
	}

	buf.WriteByte('"')
	return buf.String(), nil
}

// ex. \$ => $
func bashUnquote(str string) (string, error) {
	var buf bytes.Buffer

	r := strings.NewReader(str)
	var escapeNext bool
	for {
		char, err := r.ReadByte()
		if err != nil {
			break
		}

		if escapeNext {
			if isSpecialChar(char) {
				buf.WriteByte(char)
			} else {
				buf.WriteByte('\\')
				buf.WriteByte(char)
			}
			escapeNext = false
		} else {
			if char == '\\' {
				escapeNext = true
			} else {
				buf.WriteByte(char)
			}
		}
	}

	if escapeNext {
		return "", errors.New("escape at end")
	}
	return buf.String(), nil
}

var regLine = regexp.MustCompile(`^export\s([^\s=]+)="(.*)";$`)

func LoadFromFile(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "#") {
			// skip comment
			continue
		}
		match := regLine.FindStringSubmatch(text)
		if len(match) == 3 {
			value := match[2]
			value, err = bashUnquote(value)
			if err == nil {
				result[match[1]] = value
			}
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return result, nil
}

func Load() (map[string]string, error) {
	return LoadFromFile(defaultFile)
}

func GetFromFile(filename, key string) (string, error) {
	m, err := LoadFromFile(filename)
	if err != nil {
		return "", err
	}
	return m[key], nil
}

func Get(key string) (string, error) {
	return GetFromFile(defaultFile, key)
}

func Save(m map[string]string) error {
	return SaveToFile(defaultFile, m)
}

func SaveToFile(filename string, m map[string]string) error {
	temp := filename + ".tmp"
	f, err := os.Create(temp)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	bw := bufio.NewWriter(f)
	err = writeTo(bw, m)
	if err != nil {
		return err
	}

	err = bw.Flush()
	if err != nil {
		return err
	}

	return os.Rename(temp, filename)
}

func writeTo(w io.Writer, m map[string]string) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// write header
	_, err := fmt.Fprintln(w, "# DDE user env file, bash script")
	if err != nil {
		return err
	}

	for _, key := range keys {
		value := m[key]
		value, err := bashQuote(value)
		if err != nil {
			continue
		}

		_, err = fmt.Fprintf(w, "export %s=%s;\n", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Delete(key string) error {
	return DeleteFromFile(defaultFile, key)
}

func DeleteFromFile(filename, key string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	m, err := LoadFromFile(filename)
	if err != nil {
		return err
	}

	if _, ok := m[key]; !ok {
		return nil
	}

	delete(m, key)
	return SaveToFile(filename, m)
}

func Set(key, value string) error {
	return SetAndSaveToFile(defaultFile, key, value)
}

func Modify(fn func(map[string]string)) error {
	return ModifyAndSaveToFile(defaultFile, fn)
}

func ModifyAndSaveToFile(filename string, fn func(map[string]string)) error {
	_, err := os.Stat(filename)
	var m map[string]string
	if os.IsNotExist(err) {
		// ignore not exist
		m = make(map[string]string)
	} else if err != nil {
		return err
	} else {
		m, err = LoadFromFile(filename)
		if err != nil {
			return err
		}
	}

	fn(m)
	return SaveToFile(filename, m)
}

func SetAndSaveToFile(filename, key, value string) error {
	return ModifyAndSaveToFile(filename, func(m map[string]string) {
		m[key] = value
	})
}
