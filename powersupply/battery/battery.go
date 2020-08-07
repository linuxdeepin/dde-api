/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package battery

import (
	"pkg.deepin.io/gir/gudev-1.0"
	"math"
)

const (
	prefixPS    = "POWER_SUPPLY_"
	pName       = prefixPS + "NAME"
	//nolint
	pStatus     = prefixPS + "STATUS"
	//nolint
	pPresent    = prefixPS + "PRESENT"
	pTechnology = prefixPS + "TECHNOLOGY"

	// voltage
	pVoltageMaxDesign = prefixPS + "VOLTAGE_MAX_DESIGN"
	pVoltageMinDesign = prefixPS + "VOLTAGE_MIN_DESIGN"
	pVoltageNow       = prefixPS + "VOLTAGE_NOW"
	pVoltagePresent   = prefixPS + "VOLTAGE_PRESENT"
	pVoltageAvg       = prefixPS + "VOLTAGE_AVG"

	pPowerNow   = prefixPS + "POWER_NOW"
	pCurrentNow = prefixPS + "CURRENT_NOW"

	// energy
	pEnergyFullDesign = prefixPS + "ENERGY_FULL_DESIGN"
	pEnergyFull       = prefixPS + "ENERGY_FULL"
	pEnergyNow        = prefixPS + "ENERGY_NOW"
	pEnergyAvg        = prefixPS + "ENERGY_AVG"

	// charge
	pChargeFullDesign = prefixPS + "CHARGE_FULL_DESIGN"
	pChargeFull       = prefixPS + "CHARGE_FULL"
	pChargeNow        = prefixPS + "CHARGE_NOW"
	pChargeAvg        = prefixPS + "CHARGE_NOW"

	pCapacity     = prefixPS + "CAPACITY"
	pModelName    = prefixPS + "MODEL_NAME"
	pManufacturer = prefixPS + "MANUFACTURER"
	pSerialNumber = prefixPS + "SERIAL_NUMBER"
)

type BatteryInfo struct {
	Manufacturer string
	ModelName    string
	SerialNumber string
	Name         string
	Technology   string

	Energy           float64
	EnergyFull       float64
	EnergyFullDesign float64
	EnergyRate       float64

	Voltage     float64
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

//Sunway notebook
//POWER_SUPPLY_NAME=BAT0
//POWER_SUPPLY_STATUS=Discharging
//POWER_SUPPLY_PRESENT=1
//POWER_SUPPLY_VOLTAGE_MAX_DESIGN=11100
//POWER_SUPPLY_VOLTAGE_NOW=12575
//POWER_SUPPLY_CURRENT_NOW=0
//POWER_SUPPLY_CHARGE_FULL_DESIGN=7000
//POWER_SUPPLY_CHARGE_FULL=7326
//POWER_SUPPLY_CHARGE_NOW=7326
//POWER_SUPPLY_CAPACITY=100
//POWER_SUPPLY_TIME_TO_EMPTY_AVG=65535
//POWER_SUPPLY_TIME_TO_FULL_AVG=65535

// notebook ?
//POWER_SUPPLY_NAME=it8568-0
//POWER_SUPPLY_STATUS=Charging
//POWER_SUPPLY_PRESENT=1
//POWER_SUPPLY_VOLTAGE_NOW=16902000
//POWER_SUPPLY_CURRENT_NOW=0
//POWER_SUPPLY_CAPACITY=96
//POWER_SUPPLY_TEMP=-2702
//POWER_SUPPLY_TECHNOLOGY=Li-ion
//POWER_SUPPLY_CHARGE_FULL=0
//POWER_SUPPLY_CHARGE_NOW=5176
//POWER_SUPPLY_CHARGE_FULL_DESIGN=24800
//POWER_SUPPLY_CYCLE_COUNT=0
//POWER_SUPPLY_ENERGY_NOW=0
//POWER_SUPPLY_POWER_AVG=0
//POWER_SUPPLY_HEALTH=Good
//POWER_SUPPLY_MANUFACTURER=EASCS

func getVoltageDesign(bat *gudev.Device) (voltage float64) {
	/* design maximum */
	voltage = bat.GetPropertyAsDouble(pVoltageMaxDesign) / 1e6
	if voltage > 1 {
		return
	}

	/* design minimum */
	voltage = bat.GetPropertyAsDouble(pVoltageMinDesign) / 1e6
	if voltage > 1 {
		return
	}

	/* current voltage */
	voltage = bat.GetPropertyAsDouble(pVoltagePresent) / 1e6
	if voltage > 1 {
		return
	}

	/* current voltage, alternate form */
	voltage = bat.GetPropertyAsDouble(pVoltageNow) / 1e6
	if voltage > 1 {
		return
	}

	/* completely guess, to avoid getting zero values */
	return 10
}

func GetBatteryInfo(bat *gudev.Device) *BatteryInfo {
	if bat.HasProperty("POWER_SUPPLY_PRESENT") {
		if !bat.GetPropertyAsBoolean("POWER_SUPPLY_PRESENT") {
			return nil
		}
	}
	/* when no present property exists, handle as present */

	name := bat.GetProperty(pName)
	technology := bat.GetProperty(pTechnology)
	manufacturer := bat.GetProperty(pManufacturer)
	modelName := bat.GetProperty(pModelName)
	serialNumber := bat.GetProperty(pSerialNumber)

	/* get the current charge */
	energy := bat.GetPropertyAsDouble(pEnergyNow) / 1e6
	if energy < 0.01 {
		energy = bat.GetPropertyAsDouble(pEnergyAvg) / 1e6
	}

	/* used to convert A to W later */
	voltageDesign := getVoltageDesign(bat)

	/* initial values */
	/* these don't change at runtime */
	energyFull := bat.GetPropertyAsDouble(pEnergyFull) / 1e6
	energyFullDesign := bat.GetPropertyAsDouble(pEnergyFullDesign) / 1e6

	/* convert charge to energy */
	if energyFull < 0.01 {
		energyFull = bat.GetPropertyAsDouble(pChargeFull) / 1e6 * voltageDesign
		energyFullDesign = bat.GetPropertyAsDouble(pChargeFullDesign) / 1e6 * voltageDesign
		// TODO
	}

	//if energyFull > energyFullDesign {
	// warning energyFull is greater than energyFullDesign
	//}

	if energyFull < 0.01 && energyFullDesign > 0.01 {
		// warning correcting energyFull using energyFullDesign
		energyFull = energyFullDesign
	}

	/* calculate how broken our battery is */
	var capacity float64 = 100
	if energyFull > 0 {
		capacity = energyFull / energyFullDesign * 100
		capacity = clamp(capacity, 0, 100)
	}

	/* this is the new value in uW */
	energyRate := math.Abs(bat.GetPropertyAsDouble(pPowerNow) / 1e6)
	if energyRate < 0.01 {
		var chargeFull float64

		/* convert charge to energy */
		if energy < 0.01 {
			energy = bat.GetPropertyAsDouble(pChargeNow) / 1e6
			if energy < 0.01 {
				energy = bat.GetPropertyAsDouble(pChargeAvg) / 1e6
			}
			energy *= voltageDesign
		}

		chargeFull = bat.GetPropertyAsDouble(pChargeFull) / 1e6
		if chargeFull < 0.01 {
			chargeFull = bat.GetPropertyAsDouble(pChargeFullDesign) / 1e6
		}

		/* If chargeFull exists, then CURRENT_NOW is always reported in uA.
		 * In the legacy case, where energy only units exist, and POWER_NOW isn't present
		 * CURRENT_NOW is power in uW. */
		energyRate = math.Abs(bat.GetPropertyAsDouble(pCurrentNow) / 1e6)
		if chargeFull != 0 {
			energyRate *= voltageDesign
		}
	}

	/* some batteries don't update last_full attribute */
	if energy > energyFull {
		// warning energy bigger than full
		energyFull = energy
	}

	/* present voltage  */
	voltage := bat.GetPropertyAsDouble(pVoltageNow) / 1e6
	if voltage < 0.01 {
		voltage = bat.GetPropertyAsDouble(pVoltageAvg) / 1e6
	}

	/* sanity check to less than 100W */
	if energyRate > 100 {
		energyRate = 0
	}

	/* the hardware reporting failed --try to calculate this */
	// if energyRate < 0.01 {
	// TODO
	// }

	/* get a precise percentage */
	var percentage float64 = 0
	if bat.HasProperty(pCapacity) {
		percentage = bat.GetPropertyAsDouble(pCapacity)
		percentage = clamp(percentage, 0, 100)

		/* for devices which provide capacity, but not {energy,charge}_now */
		if energy < 0.1 && energyFull > 0 {
			energy = energyFull * percentage / 100
		}
	} else if energyFull > 0 {
		percentage = 100 * energy / energyFull
		percentage = clamp(percentage, 0, 100)
	}

	/* some batteries give out massive rate values when nearyly empty */
	if energy < 0.1 {
		energyRate = 0
	}

	// TODO
	status := parseStatus(bat.GetProperty("POWER_SUPPLY_STATUS"))

	/* calculate a quick and dirty time remaining value */
	var timeToEmpty, timeToFull uint64
	if energyRate > 0 {
		if status == StatusDischarging {
			timeToEmpty = uint64(3600 * (energy / energyRate))
		} else if status == StatusCharging {
			timeToFull = uint64(3600 * ((energyFull - energy) / energyRate))
		}
	}

	/* check the remaining thime is under a set limit, to deal with broken
	primary batteries rate */
	if timeToEmpty > 240*60*60 { /* ten days for discharging */
		timeToEmpty = 0
	}
	if timeToFull > 20*60*60 { /* 20 hours for charging */
		timeToFull = 0
	}

	return &BatteryInfo{
		Manufacturer: manufacturer,
		ModelName:    modelName,
		SerialNumber: serialNumber,
		Name:         name,
		Technology:   technology,

		Energy:           energy,
		EnergyFull:       energyFull,
		EnergyFullDesign: energyFullDesign,
		EnergyRate:       energyRate,

		Voltage:     voltage,
		Percentage:  percentage,
		Capacity:    capacity,
		Status:      status,
		TimeToEmpty: timeToEmpty,
		TimeToFull:  timeToFull,
	}
}

func clamp(val, min, max float64) float64 {
	if val < min {
		val = min
	} else if val > max {
		val = max
	}
	return val
}
