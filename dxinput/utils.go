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
	"fmt"
	"os"

	"pkg.deepin.io/dde/api/dxinput/utils"
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
	datas, num := utils.GetProperty(id, prop)
	if len(datas) == 0 || num != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' values failed",
			id, prop)
	}

	return utils.ReadInt8(datas, nitems), nil
}

func getInt16Prop(id int32, prop string, nitems int32) ([]int16, error) {
	datas, num := utils.GetProperty(id, prop)
	if len(datas) == 0 || num != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' values failed",
			id, prop)
	}

	return utils.ReadInt16(datas, nitems), nil
}

func getInt32Prop(id int32, prop string, nitems int32) ([]int32, error) {
	datas, num := utils.GetProperty(id, prop)
	if len(datas) == 0 || num != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' values failed",
			id, prop)
	}

	return utils.ReadInt32(datas, nitems), nil
}

func getFloat32Prop(id int32, prop string, nitems int32) ([]float32, error) {
	datas, num := utils.GetProperty(id, prop)
	if len(datas) == 0 || num != nitems {
		return nil, fmt.Errorf("Get prop '%v -- %s' values failed",
			id, prop)
	}

	return utils.ReadFloat32(datas, nitems), nil
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
