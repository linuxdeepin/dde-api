package dxinput

import (
	"fmt"
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	propOff            string = "Synaptics Off"
	propScrollDistance        = "Synaptics Scrolling Distance"
	propEdgeScroll            = "Synaptics Edge Scrolling"
	propTwoFingerScrol        = "Synaptics Two-Finger Scrolling"
	propTapAction             = "Synaptics Tap Action"
)

type Touchpad struct {
	Id   int32
	Name string
}

/**
 * touchpad properties see:
 * http://www.x.org/archive/X11R7.5/doc/man/man4/synaptics.4.html#sect4
 *
 * Also use 'xinput list-props <id>' to list these props.
 **/
func NewTouchpad(id int32) (*Touchpad, error) {
	info := utils.ListDevice().Get(id)
	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %v", id)
	}

	if info.Type != utils.DevTypeTouchpad {
		return nil, fmt.Errorf("Device id '%v' not a touchpad", id)
	}

	return &Touchpad{
		Id:   info.Id,
		Name: info.Name,
	}, nil
}

/**
 * Property 'Synaptics Off' 8 bit, valid values (0, 1, 2):
 *	Value 0: Touchpad is enabled
 *	Value 1: Touchpad is switched off
 *	Value 2: Only tapping and scrolling is switched off
 **/
func (tpad *Touchpad) Enable(enabled bool) error {
	err := enableDevice(tpad.Id, enabled)
	if err != nil {
		return err
	}

	if enabled == tpad.IsEnabled() {
		return nil
	}

	var values []int8
	if enabled {
		values = []int8{0}
	} else {
		values = []int8{1}
	}

	return utils.SetInt8Prop(tpad.Id, propOff, values)
}

func (tpad *Touchpad) IsEnabled() bool {
	if !isDeviceEnabled(tpad.Id) {
		return false
	}

	values, err := getInt8Prop(tpad.Id, propOff, 1)
	if err != nil {
		return false
	}

	return (values[0] == 0)
}

func (tpad *Touchpad) EnableLeftHanded(enabled bool) error {
	if enabled == tpad.CanLeftHanded() {
		return nil
	}

	return utils.SetLeftHanded(uint32(tpad.Id), tpad.Name, enabled)
}

func (tpad *Touchpad) CanLeftHanded() bool {
	return utils.CanLeftHanded(uint32(tpad.Id), tpad.Name)
}

/**
 * Property 'Synaptics Tap Action' 8 bit,
 * up to MAX_TAP values (see synaptics.h), 0 disables an element.
 * order: RT, RB, LT, LB, F1, F2, F3.
 *	Option "RTCornerButton" "integer":
 *		Which mouse button is reported on a right top corner tap.
 *	Option "RBCornerButton" "integer":
 *		Which mouse button is reported on a right bottom corner tap.
 *	Option "LTCornerButton" "integer":
 *		Which mouse button is reported on a left top corner tap.
 *	Option "LBCornerButton" "integer":
 *		Which mouse button is reported on a left bottom corner tap.
 *	Option "TapButton1" "integer":
 *		Which mouse button is reported on a non-corner one-finger tap.
 *	Option "TapButton2" "integer":
 *		Which mouse button is reported on a non-corner two-finger tap.
 *	Option "TapButton3" "integer":
 *		Which mouse button is reported on a non-corner
 *		three-finger tap.
 **/
func (tpad *Touchpad) EnableTapToClick(enabled bool) error {
	if enabled == tpad.CanTapToClick() {
		return nil
	}

	values, err := getInt8Prop(tpad.Id, propTapAction, 7)
	if err != nil {
		return err
	}

	if !enabled {
		values[4], values[5], values[6] = 0, 0, 0
	} else {
		if tpad.CanLeftHanded() {
			values[4], values[5], values[6] = 3, 1, 2
		} else {
			values[4], values[5], values[6] = 1, 3, 2
		}
	}

	return utils.SetInt8Prop(tpad.Id, propTapAction, values)
}

func (tpad *Touchpad) CanTapToClick() bool {
	values, err := getInt8Prop(tpad.Id, propTapAction, 7)
	if err != nil {
		return false
	}

	if tpad.CanLeftHanded() {
		if values[4] == 3 && values[5] == 1 && values[6] == 2 {
			return true
		}
	} else {
		if values[4] == 1 && values[5] == 3 && values[6] == 2 {
			return true
		}
	}

	return false
}

/**
 * Property "Synaptics Edge Scrolling" 8 bit (BOOL), 3 values, vertical,
 * horizontal, corner. :
 *	Option "VertEdgeScroll" "boolean":
 *		Enable vertical scrolling when dragging along the right edge.
 *	Option "HorizEdgeScroll" "boolean" :
 *		Enable horizontal scrolling when dragging along
 *		the bottom edge.
 *	Option "CornerCoasting" "boolean":
 *		Enable edge scrolling to continue while the finger stays
 *		in an edge corner.
 **/
func (tpad *Touchpad) EnableEdgeScroll(enabled bool) error {
	if enabled == tpad.CanEdgeScroll() {
		return nil
	}

	values, err := getInt8Prop(tpad.Id, propEdgeScroll, 3)
	if err != nil {
		return err
	}

	if enabled {
		values[0], values[1], values[2] = 1, 1, 1
	} else {
		values[0], values[1], values[2] = 0, 0, 0
	}

	return utils.SetInt8Prop(tpad.Id, propEdgeScroll, values)
}

func (tpad *Touchpad) CanEdgeScroll() bool {
	values, err := getInt8Prop(tpad.Id, propEdgeScroll, 3)
	if err != nil {
		return false
	}

	if values[0] != 1 || values[1] != 1 {
		return false
	}

	return true
}

/**
 * Property 'Synaptics Two-Finger Scrolling' 8 bit (BOOL),
 * 2 values, vertical, horizontal.
 *	Option "VertTwoFingerScroll" "boolean":
 *		Enable vertical scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 *	Option "HorizTwoFingerScroll" "boolean" :
 *		Enable horizontal scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 **/
func (tpad *Touchpad) EnableTwoFingerScroll(vert, horiz bool) error {
	oldVert, oldHoriz := tpad.CanTwoFingerScroll()
	if oldVert == vert && oldHoriz == horiz {
		return nil
	}

	var (
		newVert  int8 = 0
		newHoriz int8 = 0
	)
	if vert {
		newVert = 1
	}
	if horiz {
		newHoriz = 1
	}

	return utils.SetInt8Prop(tpad.Id, propTwoFingerScrol,
		[]int8{newVert, newHoriz})
}

func (tpad *Touchpad) CanTwoFingerScroll() (bool, bool) {
	values, err := getInt8Prop(tpad.Id, propTwoFingerScrol, 2)
	if err != nil {
		return false, false
	}

	return (values[0] == 1), (values[1] == 1)
}

/**
 * Property "Synaptics Scrolling Distance" 32 bit, 2 values, vert, horiz.
 *	Option "VertScrollDelta" "integer":
 *		Move distance of the finger for a scroll event.
 *	Option "HorizScrollDelta" "integer" :
 *		Move distance of the finger for a scroll event.
 *
 * if delta = 0, use value from property getting
 **/
func (tpad *Touchpad) EnableNaturalScroll(enabled bool) error {
	if enabled == tpad.CanNaturalScroll() {
		return nil
	}

	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return err
	}

	if enabled {
		values[0], values[1] = -absInt32(values[0]), -absInt32(values[1])
	} else {
		values[0], values[1] = absInt32(values[0]), absInt32(values[1])
	}

	return utils.SetInt32Prop(tpad.Id, propScrollDistance, values)
}

func (tpad *Touchpad) CanNaturalScroll() bool {
	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return false
	}

	if values[0] < 0 || values[1] < 0 {
		return true
	}

	return false
}

func (tpad *Touchpad) SetScrollDistance(vert, horiz int32) error {
	oldVert, oldHoriz := tpad.GetScrollDistance()
	if oldVert == vert && oldHoriz == horiz {
		return nil
	}

	if tpad.CanNaturalScroll() {
		vert = -vert
		horiz = -horiz
	}

	return utils.SetInt32Prop(tpad.Id, propScrollDistance,
		[]int32{vert, horiz})
}

func (tpad *Touchpad) GetScrollDistance() (int32, int32) {
	values, err := getInt32Prop(tpad.Id, propScrollDistance, 2)
	if err != nil {
		return 0, 0
	}

	return absInt32(values[0]), absInt32(values[1])
}

func (tpad *Touchpad) SetMotionAcceleration(accel float32) error {
	return setMotionAcceleration(tpad.Id, accel)
}

func (tpad *Touchpad) GetMotionAcceleration() (float32, error) {
	return getMotionAcceleration(tpad.Id)
}

func (tpad *Touchpad) SetMotionThreshold(thres float32) error {
	return setMotionThreshold(tpad.Id, thres)
}

func (tpad *Touchpad) GetMotionThreshold() (float32, error) {
	return getMotionThreshold(tpad.Id)
}
