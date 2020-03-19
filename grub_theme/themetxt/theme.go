/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package themetxt

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	ComponentTypeBootMenu         = "boot_menu"
	ComponentTypeProgressBar      = "progress_bar"
	ComponentTypeCircularProgress = "circular_progress"
	ComponentTypeLabel            = "label"
	ComponentTypeImage            = "image"
	ComponentTypeHBox             = "hbox"
	ComponentTypeVBox             = "vbox"
	ComponentTypeCanvas           = "canvas"
)

type Property struct {
	name  string
	value interface{}
}

type Length interface {
	GetConvertFunc() func(val float64) float64
}

// 50
type AbsNum int

func (v AbsNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		return float64(v)
	}
}

// 50%
type RelNum int

func (v RelNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		return float64(v) / 100.0 * val
	}
}

// 50%-10
// rel: 50
// abs: 10
type CombinedNum struct {
	Rel int
	Abs int
	Op  CombinedNumOp
}

type CombinedNumOp int

const (
	CombinedNumAdd CombinedNumOp = iota
	CombinedNumSub
)

func (v CombinedNum) GetConvertFunc() func(val float64) float64 {
	return func(val float64) float64 {
		rel := float64(v.Rel) / 100.0 * val
		return rel - float64(v.Abs)
	}
}

type Component struct {
	Type     string
	Props    []*Property
	Children []*Component
}

func (c *Component) GetProp(name string) (interface{}, bool) {
	return getProp(c.Props, name)
}

func (c *Component) GetPropString(name string) (string, bool) {
	return getPropString(c.Props, name)
}

func (c *Component) GetPropLength(name string) (Length, bool) {
	return getPropLength(c.Props, name)
}

func (c *Component) GetPropBool(name string) (bool, bool) {
	return getPropBool(c.Props, name)
}

func (c *Component) GetPropInt(name string) (int, bool) {
	return getPropInt(c.Props, name)
}

func (c *Component) SetProp(name string, value interface{}) {
	for _, prop := range c.Props {
		if prop.name == name {
			prop.value = value
			return
		}
	}

	newProp := &Property{name: name, value: value}
	c.Props = append(c.Props, newProp)
}

func (c *Component) Dump(indent int) {
	indentStr := strings.Repeat(" ", indent*4)
	fmt.Printf("%s+ %s {\n", indentStr, c.Type)

	for _, prop := range c.Props {
		fmt.Printf("%s    %s = %T %#v\n", indentStr, prop.name, prop.value, prop.value)
	}

	for _, child := range c.Children {
		child.Dump(indent + 1)
	}

	fmt.Printf("%s}\n", indentStr)
}

func (c *Component) writeTo(w io.Writer, indent int) (n int64, err error) {
	indentStr := strings.Repeat(" ", indent*4)
	var pn int
	pn, err = fmt.Fprintf(w, "%s+ %s {\n", indentStr, c.Type)
	n += int64(pn)
	if err != nil {
		return
	}

	for _, prop := range c.Props {
		if strings.HasPrefix(prop.name, "_") {
			continue
		}
		pn, err = fmt.Fprintf(w, "%s    %s = %s\n", indentStr, prop.name,
			propValueToString(prop.value))
		n += int64(pn)
		if err != nil {
			return
		}
	}

	for _, child := range c.Children {
		var cn int64
		cn, err = child.writeTo(w, indent+1)
		n += cn
		if err != nil {
			return
		}
	}

	pn, err = fmt.Fprintf(w, "%s}\n", indentStr)
	n += int64(pn)
	return
}

func (c *Component) WriteTo(w io.Writer) (n int64, err error) {
	return c.writeTo(w, 0)
}

func propValueToString(value interface{}) string {
	switch val := value.(type) {
	case string:
		return strconv.Quote(val)
	case AbsNum:
		return strconv.Itoa(int(val))
	case RelNum:
		return strconv.Itoa(int(val)) + "%"
	case CombinedNum:
		var format string
		switch val.Op {
		case CombinedNumAdd:
			format = "%d%%+%d"
		case CombinedNumSub:
			format = "%d%%-%d"
		}
		return fmt.Sprintf(format, val.Rel, val.Abs)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func getPropString(props []*Property, name string) (string, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(string), true
	}
	return "", false
}

func getPropBool(props []*Property, name string) (bool, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(bool), true
	}
	return false, false
}

func getPropInt(props []*Property, name string) (int, bool) {
	v, ok := getProp(props, name)
	if ok {
		switch val := v.(type) {
		case int:
			return val, true
		case AbsNum:
			return int(val), true
		default:
			panic(fmt.Errorf("can not convert %#v to int", v))
		}
	}
	return 0, false
}

func getPropLength(props []*Property, name string) (Length, bool) {
	v, ok := getProp(props, name)
	if ok {
		return v.(Length), true
	}
	return nil, false
}

func getProp(props []*Property, name string) (interface{}, bool) {
	for _, prop := range props {
		if prop.name == name {
			return prop.value, true
		}
	}
	return nil, false
}

type Theme struct {
	Props      []*Property
	Components []*Component
}

func (t *Theme) SetProp(name string, value interface{}) {
	for _, prop := range t.Props {
		if prop.name == name {
			prop.value = value
			return
		}
	}

	newProp := &Property{name: name, value: value}
	t.Props = append(t.Props, newProp)
}

func (t *Theme) GetPropString(name string) (string, bool) {
	return getPropString(t.Props, name)
}

func (t *Theme) Dump() {
	for _, prop := range t.Props {
		fmt.Printf("%s : %T %#v\n", prop.name, prop.value, prop.value)
	}
	for _, comp := range t.Components {
		comp.Dump(0)
	}
}

func (t *Theme) WriteTo(w io.Writer) (n int64, err error) {
	for _, prop := range t.Props {
		var pn int
		pn, err = fmt.Fprintf(w, "%s: %s\n", prop.name, propValueToString(prop.value))
		n += int64(pn)
		if err != nil {
			return
		}
	}
	for _, comp := range t.Components {
		var cn int64
		cn, err = comp.WriteTo(w)
		n += cn
		if err != nil {
			return
		}
	}
	return
}

func ParseThemeFile(filename string) (*Theme, error) {
	v, err := ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return v.(*Theme), nil
}
