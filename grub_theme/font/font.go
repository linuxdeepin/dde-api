package font

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
)

type Face struct {
	Name      string
	Family    string
	Weight    string
	Slant     string
	PointSize int

	MaxWidth  int
	MaxHeight int
	Ascent    int
	Descent   int

	CharIndexes []charIndex
	br          *bytes.Reader
}

func LoadFont(filename string) (*Face, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(data)

	var sections = make(map[string]*section)
	for {
		section, err := parseSection(br)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		sections[section.name] = section
	}

	// section FILE
	section, ok := sections["FILE"]
	if !ok {
		return nil, errors.New("not found section FILE")
	}
	file0 := section.getString()
	if file0 != "PFF2" {
		return nil, errors.New("FILE is not PFF2")
	}

	var face Face
	// section NAME
	section, ok = sections["NAME"]
	if !ok {
		return nil, errors.New("not found section NAME")
	}
	face.Name = section.getString()

	// section FAMI
	section, ok = sections["FAMI"]
	if !ok {
		return nil, errors.New("not found section FAMI")
	}
	face.Family = section.getString()

	// section WEIG
	section, ok = sections["WEIG"]
	if !ok {
		return nil, errors.New("not found section WEIG")
	}
	face.Weight = section.getString()

	// section SLAN
	section, ok = sections["SLAN"]
	if !ok {
		return nil, errors.New("not found section SLAN")
	}
	face.Slant = section.getString()

	// section PTSZ
	section, ok = sections["PTSZ"]
	if !ok {
		return nil, errors.New("not found section PTSZ")
	}
	face.PointSize = int(section.getUint16BE())

	// section MAXW
	section, ok = sections["MAXW"]
	if !ok {
		return nil, errors.New("not found section MAXW")
	}
	face.MaxWidth = int(section.getUint16BE())

	// section MAXH
	section, ok = sections["MAXH"]
	if !ok {
		return nil, errors.New("not found section MAXH")
	}
	face.MaxHeight = int(section.getUint16BE())

	// section ASCE
	section, ok = sections["ASCE"]
	if !ok {
		return nil, errors.New("not found section ASCE")
	}
	face.Ascent = int(section.getUint16BE())

	// section DESC
	section, ok = sections["DESC"]
	if !ok {
		return nil, errors.New("not found section DESC")
	}
	face.Descent = int(section.getUint16BE())

	// section CHIX
	section, ok = sections["CHIX"]
	if !ok {
		return nil, errors.New("not found section CHIX")
	}

	chix, err := parseCharIndexes(section.data)
	if err != nil {
		return nil, err
	}
	face.CharIndexes = chix
	face.br = br
	return &face, nil
}

type section struct {
	name string
	data []byte
}

func (s *section) getString() string {
	data := bytes.TrimRight(s.data, "\x00")
	return string(data)
}

func (s *section) getUint16BE() uint16 {
	return binary.BigEndian.Uint16(s.data)
}

func parseSection(r io.Reader) (*section, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	s := &section{}
	s.name = string(buf)

	var length uint32
	binary.Read(r, binary.BigEndian, &length)

	if s.name == "DATA" {
		return nil, io.EOF
	}

	buf = make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, fmt.Errorf("readSection section %q length %d err: %v", s.name, length, err)
	}
	s.data = buf

	return s, nil
}

func parseCharIndexes(data []byte) ([]charIndex, error) {
	count := len(data) / (4 + 4 + 1)
	r := bytes.NewReader(data)
	result := make([]charIndex, 0, count)

	for {
		elem, err := parseCharIndex(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		result = append(result, elem)
	}
	return result, nil
}

type charIndex struct {
	unicodeCodePoint uint32
	flag             byte
	offset           uint32
}

func parseCharIndex(r *bytes.Reader) (charIndex, error) {
	var unicodeCodePoint uint32
	err := binary.Read(r, binary.BigEndian, &unicodeCodePoint)
	if err != nil {
		return charIndex{}, err
	}

	flag, err := r.ReadByte()
	if err != nil {
		return charIndex{}, err
	}

	var offset uint32
	err = binary.Read(r, binary.BigEndian, &offset)
	if err != nil {
		return charIndex{}, err
	}
	return charIndex{
		unicodeCodePoint: unicodeCodePoint,
		flag:             flag,
		offset:           offset,
	}, nil
}

func (f *Face) Close() error {
	return nil
}

func (f *Face) findCharIndex(r rune) *charIndex {
	for _, chixItem := range f.CharIndexes {
		if rune(chixItem.unicodeCodePoint) == r {
			return &chixItem
		}
	}
	return nil
}

func (f *Face) findChar(r rune) *CharInfo {
	charIdx := f.findCharIndex(r)
	if charIdx == nil {
		return nil
	}

	_, err := f.br.Seek(int64(charIdx.offset), io.SeekStart)
	if err != nil {
		return nil
	}
	charInfo, err := parseCharInfo(f.br)
	if err != nil {
		return nil
	}
	return charInfo
}

type CharInfo struct {
	width       uint16
	height      uint16
	xOffset     int16
	yOffset     int16
	deviceWidth int16
	mask        image.Image
}

func parseCharInfo(r io.Reader) (*CharInfo, error) {
	var d CharInfo
	err := binary.Read(r, binary.BigEndian, &d.width)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &d.height)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &d.xOffset)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &d.yOffset)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &d.deviceWidth)
	if err != nil {
		return nil, err
	}

	count := 0
	var b byte
	buf := make([]byte, 1)

	img := image.NewAlpha(image.Rect(0, 0, int(d.width), int(d.height)))

	for y := 0; y < int(d.height); y++ {
		for x := 0; x < int(d.width); x++ {
			if count == 0 {
				_, err = r.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				b = buf[0]
			}

			val := b & (1 << uint(7-count))
			if val != 0 {
				img.SetAlpha(x, y, color.Alpha{A: 255})
			}

			if count == 7 {
				// reset count
				count = 0
			} else {
				count++
			}
		}
	}

	d.mask = img
	return &d, nil
}

func (f *Face) Height() int {
	return f.Ascent + f.Descent
}
