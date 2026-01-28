// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dxinput

import (
	"fmt"
	"os"

	"github.com/linuxdeepin/dde-api/dxinput/utils"
)

const (
	propDeviceEnabled        string = "Device Enabled"
	propConstantDeceleration string = "Device Accel Constant Deceleration"
	propAdaptiveDeceleration string = "Device Accel Adaptive Deceleration"
	propVelocityScaling      string = "Device Accel Velocity Scaling"
	propCoordTransMatrix     string = "Coordinate Transformation Matrix"
)

const (
	// see also randr
	RotationDirectionNormal   uint8 = 1
	RotationDirectionLeft     uint8 = 2
	RotationDirectionInverted uint8 = 4
	RotationDirectionRight    uint8 = 8
)

var (
	rotationNormal   = []float32{1, 0, 0, 0, 1, 0, 0, 0, 1}   // 0°
	rotationLeft     = []float32{0, -1, 1, 1, 0, 0, 0, 0, 1}  // clockwise 90°
	rotationInverted = []float32{-1, 0, 1, 0, -1, 1, 0, 0, 1} // clockwise or counterclockwise 180°
	rotationRight    = []float32{0, 1, 0, -1, 0, 1, 0, 0, 1}  // counterclockwise 90°
)

var (
	globalWayland bool
)

func init() {
	if len(os.Getenv("WAYLAND_DISPLAY")) != 0 {
		globalWayland = true
	}
}

/**
 * Prop: "Device Enabled", 8 bit, 1 value
 * valid values: 0, 1.
 **/
func enableDevice(id int32, enabled bool) error {
	if enabled == isDeviceEnabled(id) {
		return nil
	}

	var values []int8
	if enabled {
		values = []int8{1}
	} else {
		values = []int8{0}
	}

	return utils.SetInt8Prop(id, propDeviceEnabled, values)
}

func isDeviceEnabled(id int32) bool {
	values, err := getInt8Prop(id, propDeviceEnabled, 1)
	if err != nil {
		return false
	}

	return (values[0] == 1)
}

/**
 * Pointer device motion acceleration and threshold
 *
 * "Device Accel Constant Deceleration" 32 1 value float
 *
 * "Device Accel Adaptive Deceleration" 32 1 value float
 *
 * Detail info see:
 * http://510x.se/notes/posts/Changing_mouse_acceleration_in_Debian_and_Linux_in_general/
 **/

func setMotionAcceleration(id int32, accel float32) error {
	if accel <= 0 {
		return fmt.Errorf("Invalid accel value: %v, must > 0", accel)
	}
	value, err := getMotionAcceleration(id)
	if err != nil {
		return err
	}

	if accel > value-0.01 && accel < value+0.01 {
		return nil
	}

	return utils.SetFloat32Prop(id, propConstantDeceleration,
		[]float32{accel})
}

func getMotionAcceleration(id int32) (float32, error) {
	values, err := getFloat32Prop(id, propConstantDeceleration, 1)
	if err != nil {
		return 0, err
	}

	return values[0], nil
}

func setMotionThreshold(id int32, thres float32) error {
	value, err := getMotionThreshold(id)
	if err != nil {
		return err
	}

	if thres > value-0.01 && thres < value+0.01 {
		return nil
	}

	return utils.SetFloat32Prop(id, propAdaptiveDeceleration,
		[]float32{thres})
}

func getMotionThreshold(id int32) (float32, error) {
	values, err := getFloat32Prop(id, propAdaptiveDeceleration, 1)
	if err != nil {
		return 0, err
	}

	return values[0], nil
}

func setMotionScaling(id int32, scaling float32) error {
	old, err := getMotionScaling(id)
	if err != nil {
		return err
	}

	if scaling > old-0.01 && scaling < old+0.01 {
		return nil
	}
	return utils.SetFloat32Prop(id, propVelocityScaling, []float32{scaling})
}

func getMotionScaling(id int32) (float32, error) {
	values, err := getFloat32Prop(id, propVelocityScaling, 1)
	if err != nil {
		return 0, err
	}
	return values[0], nil
}

func setRotation(id int32, rotation uint8) error {
	var v []float32
	switch rotation {
	case RotationDirectionNormal:
		v = rotationNormal
	case RotationDirectionLeft:
		v = rotationLeft
	case RotationDirectionInverted:
		v = rotationInverted
	case RotationDirectionRight:
		v = rotationRight
	default:
		return fmt.Errorf("Invalid rotation value: %d", rotation)
	}

	values, _ := getFloat32Prop(id, propCoordTransMatrix, 9)
	if isFloat32ListEqual(v, values) {
		return nil
	}

	return utils.SetFloat32Prop(id, propCoordTransMatrix, v)
}

func setTransformationMatrix(id int32, m [9]float32) error {
	values, err := getFloat32Prop(id, propCoordTransMatrix, 9)
	if err != nil {
		return err
	}

	v := m[:]
	if isFloat32ListEqual(v, values) {
		return nil
	}

	return utils.SetFloat32Prop(id, propCoordTransMatrix, v)
}

func getInt8Prop(id int32, prop string, nitems int32) ([]int8, error) {
	datas, nBytes := utils.GetProperty(id, prop)
	if len(datas) == 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: property data is empty",
			id, prop)
	}
	// For int8, nBytes should equal nitems (1 byte per item)
	if nBytes != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: expected %d items but got %d bytes",
			id, prop, nitems, nBytes)
	}

	return utils.ReadInt8(datas, nitems), nil
}

func getInt16Prop(id int32, prop string, nitems int32) ([]int16, error) {
	datas, nBytes := utils.GetProperty(id, prop)
	if len(datas) == 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: property data is empty",
			id, prop)
	}
	// Check if nBytes is divisible by 2 (int16 size)
	if nBytes%2 != 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: byte count %d is not divisible by 2 (int16 size)",
			id, prop, nBytes)
	}
	// For int16, nBytes should equal nitems * 2 (2 bytes per item)
	if nBytes/2 != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: expected %d items but got %d bytes (%d items)",
			id, prop, nitems, nBytes, nBytes/2)
	}

	return utils.ReadInt16(datas, nitems), nil
}

func getInt32Prop(id int32, prop string, nitems int32) ([]int32, error) {
	datas, nBytes := utils.GetProperty(id, prop)
	if len(datas) == 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: property data is empty",
			id, prop)
	}
	// Check if nBytes is divisible by 4 (int32 size)
	if nBytes%4 != 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: byte count %d is not divisible by 4 (int32 size)",
			id, prop, nBytes)
	}
	// For int32, nBytes should equal nitems * 4 (4 bytes per item)
	if nBytes/4 != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: expected %d items but got %d bytes (%d items)",
			id, prop, nitems, nBytes, nBytes/4)
	}

	return utils.ReadInt32(datas, nitems), nil
}

func getFloat32Prop(id int32, prop string, nItems int32) ([]float32, error) {
	datas, nBytes := utils.GetProperty(id, prop)
	if len(datas) == 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: property data is empty",
			id, prop)
	}
	// Check if nBytes is divisible by 4 (float32 size)
	if nBytes%4 != 0 {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: byte count %d is not divisible by 4 (float32 size)",
			id, prop, nBytes)
	}
	// utils.GetProperty returns byte count, float32 occupies 4 bytes, divide by 4 to get item count
	if nBytes/4 != nItems {
		return nil, fmt.Errorf("Get prop '%v -- %s' failed: expected %d items but got %d bytes (%d items)",
			id, prop, nItems, nBytes, nBytes/4)
	}

	return utils.ReadFloat32(datas, nItems), nil
}

func absInt32(v int32) int32 {
	switch {
	case v < 0:
		return -v
	case v > 0:
		return v
	}
	return v
}

func isFloat32ListEqual(l1, l2 []float32) bool {
	len1, len2 := len(l1), len(l2)
	if len1 != len2 {
		return false
	}
	for i := 0; i < len1; i++ {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}
