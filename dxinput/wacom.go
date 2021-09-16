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

package dxinput

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	. "pkg.deepin.io/dde/api/dxinput/common"
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	cmdXSetWacom = "xsetwacom"

	cmdKeyArea      string = "Area"
	cmdKeyResetArea string = "ResetArea"
	cmdKeyMode      string = "mode"
	cmdKeyButton    string = "Button"
	cmdKeyRotate    string = "Rotate"
	cmdKeySuppress  string = "Suppress"
	//(x1, y2, x2, y2) red(x1, y1), blue(x2, y2), green(Threshold)
	cmdKeyPressureCurve string = "PressureCurve"
	cmdKeyThreshold     string = "Threshold"
	cmdKeyRawSample     string = "RawSample"
	// such as 'VGA1'
	cmdKeyMapToOutput string = "MapToOutput"
)

const (
	WacomTypeUnknown = iota
	WacomTypeStylus
	WacomTypeEraser
	WacomTypePad
)

type Wacom struct {
	Id   int32
	Name string
}

func NewWacom(id int32) (*Wacom, error) {
	info := utils.ListDevice().Get(id)
	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %v", id)
	}
	return NewWacomFromDevInfo(info)
}

func NewWacomFromDevInfo(dev *DeviceInfo) (*Wacom, error) {
	if dev == nil || dev.Type != DevTypeWacom {
		return nil, fmt.Errorf("Not a wacom device(%d - %s)", dev.Id, dev.Name)
	}

	return &Wacom{
		Id:   dev.Id,
		Name: dev.Name,
	}, nil
}

func (w *Wacom) QueryType() int {
	nameLower := strings.ToLower(w.Name)
	switch {
	case strings.Contains(nameLower, "stylus"):
		return WacomTypeStylus
	case strings.Contains(nameLower, "eraser"):
		return WacomTypeEraser
	case strings.Contains(nameLower, "pad"):
		return WacomTypePad
	default:
		return WacomTypeUnknown
	}
}

// Area x1 y1 x2 y2
// Set  the tablet input area in device coordinates in the form top
// left x/y and bottom right x/y.
func (w *Wacom) SetArea(x1, y1, x2, y2 int) error {
	var cmd = fmt.Sprintf("%s set %v %s %v %v %v %v", cmdXSetWacom, w.Id,
		cmdKeyArea, x1, y1, x2, y2)
	return doAction(cmd)
}

func (w *Wacom) ResetArea() error {
	var cmd = fmt.Sprintf("%s set %v %s", cmdXSetWacom, w.Id, cmdKeyResetArea)
	return doAction(cmd)
}

func (w *Wacom) getIdAsStr() string {
	return strconv.FormatInt(int64(w.Id), 10)
}

// GetArea get the tablet input area
func (w *Wacom) GetArea() (x1, y1, x2, y2 int, err error) {
	var out []byte
	out, err = exec.Command(cmdXSetWacom, "get", w.getIdAsStr(), cmdKeyArea).Output()
	if err != nil {
		return
	}
	// parse out
	_, err = fmt.Fscanln(bytes.NewReader(out), &x1, &y1, &x2, &y2)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return
}

// Rotate valid values: none|half|cw|ccw
// none: the tablet is not rotated and uses its natural rotation
// half: the tablet is rotated by 180 degrees (upside-down)
// cw  : the tablet is rotated 90 degrees clockwise
// ccw : the tablet is rotated 90 degrees counter-clockwise
func (w *Wacom) SetRotate(value string) error {
	switch value {
	case "none", "half", "cw", "ccw":
		break
	default:
		return fmt.Errorf("Invalid value: %s", value)
	}

	var cmd = fmt.Sprintf("%s set %v %s %v", cmdXSetWacom, w.Id,
		cmdKeyRotate, value)
	return doAction(cmd)
}

// Button button-number [mapping]
// Set a mapping for the specified button-number.
// Numeric  button  mappings  indicate  what  X11 button number the
// given button-number should correspond to. For example, a mapping
// of  "3" means a press of the given button-number will produce as
// a press of X11 button 3 (i.e. right click).
//
// Action mappings allow button presses  to  perform  many  events.
// They  take  the  form of a string of keywords and arguments. For
// example, "key +a +shift b -shift -a" converts the button into  a
// series  of  keystrokes,  in  this example "press a, press shift,
// press and release b, release shift, release a".
func (w *Wacom) SetButton(btn int, value string) error {
	var cmd = fmt.Sprintf("%s set %v %s %v %s", cmdXSetWacom, w.Id,
		cmdKeyButton, btn, value)
	return doAction(cmd)
}

// Mode Absolute|Relative
// Set the device mode as either  Relative  or  Absolute.
// Relative means  pointer  tracking  for  the  device  will function like a
// mouse.
// Absolute means the pointer corresponds to the device's actual position on
// the tablet or tablet PC screen.
func (w *Wacom) SetMode(mode string) error {
	switch mode {
	case "Absolute", "Relative":
		break
	default:
		return fmt.Errorf("Invalid value: %s", mode)
	}
	var cmd = fmt.Sprintf("%s set %v %s %s", cmdXSetWacom, w.Id,
		cmdKeyMode, mode)
	return doAction(cmd)
}

// PressureCurve  x1 y1 x2 y2
// A  Bezier curve of third order, composed of two anchor points
// (0,0 and 100,100) and two user modifiable control points that
// define the curve's  shape.
// Raise the curve (x1<y1 x2<y2) to "soften" the feel and
// lower the curve (x1>y1  x2>y2) for a "firmer" feel.
// Sigmoid shaped curves are permitted (x1>y1 x2<y2 or x1<y1 x2>y2).
//
// Default:  0 0 100 100, a linear curve.
// range of 0 to 100 for all four values.
func (w *Wacom) SetPressureCurve(x1, y1, x2, y2 int) error {
	if (x1 < 0 || x1 > 100) || (y1 < 0 || y1 > 100) ||
		(x2 < 0 || x2 > 100) || (y2 < 0 || y2 > 100) {
		return fmt.Errorf("Invalid value: %v %v %v %v", x1, y1, x2, y2)
	}

	var cmd = fmt.Sprintf("%s set %v %s %v %v %v %v", cmdXSetWacom, w.Id,
		cmdKeyPressureCurve, x1, y1, x2, y2)
	return doAction(cmd)
}

// Suppress level
// Set the delta (difference) cutoff level for further processing
// of incoming input tool coordinate values.
// To disable suppression use a level of 0.
// Default:  2, range of 0 to 100.
func (w *Wacom) SetSuppress(value int) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("Invalid value: %v", value)
	}

	var cmd = fmt.Sprintf("%s set %v %s %v", cmdXSetWacom, w.Id,
		cmdKeySuppress, value)
	return doAction(cmd)
}

// Threshold level
// Set the minimum pressure necessary to generate a Button event
// for the stylus tip, eraser, or touch. The pressure levels of
// all tablets are normalized to 2048 levels irregardless of the
// actual hardware supported levels. This parameter is independent
// of the PressureCurve parameter.
// Default:  27, range of 0 to 2047.
func (w *Wacom) SetThreshold(thres int) error {
	if thres < 0 || thres > 2047 {
		return fmt.Errorf("Invalid value: %v", thres)
	}

	var cmd = fmt.Sprintf("%s set %v %s %v", cmdXSetWacom, w.Id,
		cmdKeyThreshold, thres)
	return doAction(cmd)
}

// The the window size for incoming input tool raw data points
// Default: 4, range of 1 to 20
func (w *Wacom) SetRawSample(sample uint32) error {
	if sample == 0 {
		return fmt.Errorf("Invalid raw sample: %v", sample)
	}

	var cmd = fmt.Sprintf("%s set %v %s %v", cmdXSetWacom, w.Id,
		cmdKeyRawSample, sample)
	return doAction(cmd)
}

// Mapping PC screen to tablet, such as "VGA1"
func (w *Wacom) MapToOutput(output string) error {
	if len(output) == 0 {
		return nil
	}

	var cmd = fmt.Sprintf("%s set %v %s %s", cmdXSetWacom, w.Id,
		cmdKeyMapToOutput, output)
	return doAction(cmd)
}

func doAction(cmd string) error {
	// #nosec G204
	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}
	return nil
}
