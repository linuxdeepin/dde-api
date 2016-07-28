/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package battery

import (
	"gir/gudev-1.0"
)

type BatteryInfo struct {
	Name             string
	Technology       string
	Manufacturer     string
	ModelName        string
	SerialNumber     string
	EnergyFullDesign uint64
	VoltageMinDesign uint64

	EnergyFull  uint64
	EnergyNow   uint64
	PowerNow    uint64
	VoltageNow  uint64
	Percentage  float64
	Capacity    float64
	Status      Status
	TimeToEmpty uint64
	TimeToFull  uint64
}

// uevent file sample:
// Deepin
// POWER_SUPPLY_NAME=BAT0
// POWER_SUPPLY_STATUS=Discharging
// POWER_SUPPLY_PRESENT=1
// POWER_SUPPLY_TECHNOLOGY=Li-ion
// POWER_SUPPLY_CYCLE_COUNT=0
// POWER_SUPPLY_VOLTAGE_MIN_DESIGN=11400000
// POWER_SUPPLY_VOLTAGE_NOW=12690000
// POWER_SUPPLY_POWER_NOW=0
// POWER_SUPPLY_ENERGY_FULL_DESIGN=46970000
// POWER_SUPPLY_ENERGY_FULL=41250000
// POWER_SUPPLY_ENERGY_NOW=41250000
// POWER_SUPPLY_CAPACITY=100
// POWER_SUPPLY_CAPACITY_LEVEL=Normal
// POWER_SUPPLY_MODEL_NAME=LNV-45N1
// POWER_SUPPLY_MANUFACTURER=LGC
// POWER_SUPPLY_SERIAL_NUMBER=44853

// Arch Linux
// POWER_SUPPLY_NAME=BAT0
// POWER_SUPPLY_STATUS=Full
// POWER_SUPPLY_PRESENT=1
// POWER_SUPPLY_TECHNOLOGY=Li-ion
// POWER_SUPPLY_CYCLE_COUNT=0
// POWER_SUPPLY_VOLTAGE_MIN_DESIGN=14800000
// POWER_SUPPLY_VOLTAGE_NOW=16636000
// POWER_SUPPLY_CURRENT_NOW=0
// POWER_SUPPLY_CHARGE_FULL_DESIGN=2200000
// POWER_SUPPLY_CHARGE_FULL=2167000
// POWER_SUPPLY_CHARGE_NOW=2167000
// POWER_SUPPLY_CAPACITY=100
// POWER_SUPPLY_CAPACITY_LEVEL=Full
// POWER_SUPPLY_MODEL_NAME=MWL31b
// POWER_SUPPLY_MANUFACTURER=SMP-SDI2
// POWER_SUPPLY_SERIAL_NUMBER=

func GetBatteryInfo(bat *gudev.Device) *BatteryInfo {
	if bat.HasProperty("POWER_SUPPLY_PRESENT") {
		if !bat.GetPropertyAsBoolean("POWER_SUPPLY_PRESENT") {
			return nil
		}
	}
	/* when no present property exists, handle as present */

	name := bat.GetProperty("POWER_SUPPLY_NAME")
	technology := bat.GetProperty("POWER_SUPPLY_TECHNOLOGY")
	manufacturer := bat.GetProperty("POWER_SUPPLY_MANUFACTURER")
	modelName := bat.GetProperty("POWER_SUPPLY_MODEL_NAME")
	serialNumber := bat.GetProperty("POWER_SUPPLY_SERIAL_NUMBER")
	energyFullDesign := bat.GetPropertyAsUint64("POWER_SUPPLY_ENERGY_FULL_DESIGN")
	if energyFullDesign < 1 {
		energyFullDesign = bat.GetPropertyAsUint64("POWER_SUPPLY_CHARGE_FULL_DESIGN")
	}
	voltageMinDesign := bat.GetPropertyAsUint64("POWER_SUPPLY_VOLTAGE_MIN_DESIGN")

	energyNow := bat.GetPropertyAsUint64("POWER_SUPPLY_ENERGY_NOW")
	if energyNow < 1 {
		energyNow = bat.GetPropertyAsUint64("POWER_SUPPLY_CHARGE_NOW")
	}

	energyFull := bat.GetPropertyAsUint64("POWER_SUPPLY_ENERGY_FULL")
	if energyFull < 1 {
		energyFull = bat.GetPropertyAsUint64("POWER_SUPPLY_CHARGE_FULL")
	}
	if energyFull < 1 && energyFullDesign > 1 {
		energyFull = energyFullDesign
	}
	/* some batteries don't update last_full attribute */
	if energyNow > energyFull {
		energyFull = energyNow
	}

	powerNow := bat.GetPropertyAsUint64("POWER_SUPPLY_POWER_NOW")
	if powerNow < 1 {
		powerNow = bat.GetPropertyAsUint64("POWER_SUPPLY_CURRENT_NOW")
	}

	voltageNow := bat.GetPropertyAsUint64("POWER_SUPPLY_VOLTAGE_NOW")

	percentage := float64(energyNow) / float64(energyFull) * 100.0

	/* calculate how broken our battery is */
	capacity := float64(energyFull) / float64(energyFullDesign) * 100

	status := parseStatus(bat.GetProperty("POWER_SUPPLY_STATUS"))

	var timeToEmpty, timeToFull uint64
	timeToEmpty = uint64(float64(energyNow) / float64(powerNow) * 3600)
	timeToFull = uint64(float64(energyFull-energyNow) / float64(powerNow) * 3600)

	return &BatteryInfo{
		Name:             name,
		Technology:       technology,
		Manufacturer:     manufacturer,
		ModelName:        modelName,
		SerialNumber:     serialNumber,
		EnergyFullDesign: energyFullDesign,
		VoltageMinDesign: voltageMinDesign,

		EnergyFull:  energyFull,
		EnergyNow:   energyNow,
		PowerNow:    powerNow,
		VoltageNow:  voltageNow,
		Percentage:  percentage,
		Capacity:    capacity,
		Status:      status,
		TimeToEmpty: timeToEmpty,
		TimeToFull:  timeToFull,
	}
}
