/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dxinput

import (
	"fmt"
	"pkg.deepin.io/dde/api/dxinput/utils"
)

const (
	propDeviceEnabled        string = "Device Enabled"
	propConstantDeceleration        = "Device Accel Constant Deceleration"
	propAdaptiveDeceleration        = "Device Accel Adaptive Deceleration"
	propVelocityScaling             = "Device Accel Velocity Scaling"
)

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

func getStringProp(id int32, prop string) (string, error) {
	datas, num := utils.GetProperty(id, prop)
	if len(datas) == 0 {
		return "", fmt.Errorf("Get prop '%v -- %s' values failed",
			id, prop)
	}

	return string(datas[:num]), nil
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
