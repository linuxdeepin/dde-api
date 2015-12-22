package dxinput

import (
	"fmt"
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	propMidBtnEmulation      string = "Evdev Middle Button Emulation"
	propWheelEmulation              = "Evdev Wheel Emulation"
	propWheelEmulationButton        = "Evdev Wheel Emulation Button"
)

type Mouse struct {
	Id   int32
	Name string
}

func NewMouse(id int32) (*Mouse, error) {
	info := utils.ListDevice().Get(id)
	if info == nil {
		return nil, fmt.Errorf("Invalid device id: %v", id)
	}

	if info.Type != utils.DevTypeMouse {
		return nil, fmt.Errorf("Device id '%v' not a mouse", id)
	}

	return &Mouse{
		Id:   info.Id,
		Name: info.Name,
	}, nil
}

func (m *Mouse) Enable(enabled bool) error {
	return enableDevice(m.Id, enabled)
}

func (m *Mouse) IsEnabled() bool {
	return isDeviceEnabled(m.Id)
}

func (m *Mouse) EnableLeftHanded(enabled bool) error {
	if enabled == m.CanLeftHanded() {
		return nil
	}

	return utils.SetLeftHanded(uint32(m.Id), m.Name, enabled)
}

func (m *Mouse) CanLeftHanded() bool {
	return utils.CanLeftHanded(uint32(m.Id), m.Name)
}

/*
 * Evdev Middle Button Emulation
 *     1 boolean value (8 bit, 0 or 1).
 */
func (m *Mouse) EnableMiddleButtonEmulation(enabled bool) error {
	if enabled == m.CanMiddleButtonEmulation() {
		return nil
	}

	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}

	return utils.SetInt8Prop(m.Id, propMidBtnEmulation, values)
}

func (m *Mouse) CanMiddleButtonEmulation() bool {
	values, err := getInt32Prop(m.Id, propMidBtnEmulation, 1)
	if err != nil {
		return false
	}

	return (values[0] == 1)
}

/*
 * Evdev Wheel Emulation
 *    1 boolean value (8 bit, 0 or 1).
 */
func (m *Mouse) EnableWheelEmulation(enabled bool) error {
	if enabled == m.CanWheelEmulation() {
		return nil
	}
	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}
	return utils.SetInt8Prop(m.Id, propWheelEmulation, values)
}

func (m *Mouse) CanWheelEmulation() bool {
	values, err := getInt8Prop(m.Id, propWheelEmulation, 1)
	if err != nil {
		return false
	}
	return (values[0] == 1)
}

/*
 * Evdev Wheel Emulation Button
 *    1 8-bit value, allowed range 0-32, 0 disables the button.
 */
func (m *Mouse) SetWheelEmulationButton(btnNum int8) error {
	old, _ := m.GetWheelEmulationButton()
	if btnNum == old {
		return nil
	}
	return utils.SetInt8Prop(m.Id, propWheelEmulationButton,
		[]int8{btnNum})
}

func (m *Mouse) GetWheelEmulationButton() (int8, error) {
	values, err := getInt8Prop(m.Id, propWheelEmulationButton, 1)
	if err != nil {
		return -1, err
	}
	return values[0], nil
}

func (m *Mouse) EnableNaturalScroll(enabled bool) error {
	if enabled == m.CanNaturalScroll() {
		return nil
	}

	btnMap, err := utils.GetButtonMap(uint32(m.Id), m.Name)
	if err != nil {
		return err
	}

	if len(btnMap) < 5 {
		return fmt.Errorf("Invalid mouse device: button number < 5")
	}

	if enabled {
		btnMap[3], btnMap[4] = 5, 4
	} else {
		btnMap[3], btnMap[4] = 4, 5
	}

	return utils.SetButtonMap(uint32(m.Id), m.Name, btnMap)
}

func (m *Mouse) CanNaturalScroll() bool {
	btnMap, err := utils.GetButtonMap(uint32(m.Id), m.Name)
	if err != nil {
		return false
	}

	if len(btnMap) < 5 {
		return false
	}

	if btnMap[3] == 5 && btnMap[4] == 4 {
		return true
	}

	return false
}

func (m *Mouse) SetMotionAcceleration(accel float32) error {
	return setMotionAcceleration(m.Id, accel)
}

func (m *Mouse) GetMotionAcceleration() (float32, error) {
	return getMotionAcceleration(m.Id)
}

func (m *Mouse) SetMotionThreshold(thres float32) error {
	return setMotionThreshold(m.Id, thres)
}

func (m *Mouse) GetMotionThreshold() (float32, error) {
	return getMotionThreshold(m.Id)
}
