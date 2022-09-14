// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"errors"
	"os/exec"
	"sync"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-lib/dbusutil"

	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
)

//go:generate dbusutil-gen em -type Device

const (
	dbusServiceName                 = "com.deepin.api.Device"
	dbusPath                        = "/com/deepin/api/Device"
	dbusInterface                   = dbusServiceName
	rfkillBin                       = "rfkill"
	rfkillDeviceTypeBluetooth       = "bluetooth"
	unblockBluetoothDevicesActionId = "com.deepin.api.device.unblock-bluetooth-devices"
)

type Device struct {
	service      *dbusutil.Service
	mu           sync.Mutex
	callingCount int
}

func (d *Device) incCallingCount() {
	d.mu.Lock()
	d.callingCount++
	d.mu.Unlock()
}

func (d *Device) decCallingCount() {
	d.mu.Lock()
	d.callingCount--
	d.mu.Unlock()
}

func (d *Device) canQuit() bool {
	d.mu.Lock()
	count := d.callingCount
	d.mu.Unlock()
	return count == 0
}

func (*Device) GetInterfaceName() string {
	return dbusInterface
}

//  UnblockBluetoothDevice unblock bluetooth devices through rfkill
func (d *Device) UnblockBluetoothDevices(sender dbus.Sender) *dbus.Error {
	d.service.DelayAutoQuit()
	d.incCallingCount()
	defer d.decCallingCount()
	err := d.unblockBluetoothDevices(sender)
	return dbusutil.ToError(err)
}

func (d *Device) unblockBluetoothDevices(sender dbus.Sender) error {
	ok, err := checkAuthorization(unblockBluetoothDevicesActionId, string(sender))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized access")
	}

	return exec.Command(rfkillBin, "unblock", rfkillDeviceTypeBluetooth).Run()
}

func checkAuthorization(actionId string, sysBusName string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", sysBusName)

	ret, err := authority.CheckAuthorization(0, subject, actionId,
		nil, polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}

type rfkillItem struct {
	Id     json.RawMessage `json:"id"`
	Type   string          `json:"type"`
	Device string          `json:"device"`
	Soft   string          `json:"soft"`
	Hard   string          `json:"hard"`
}

func getRfkillItems() ([]rfkillItem, error) {
	output, err := exec.Command(rfkillBin, "-J").Output()
	if err != nil {
		return nil, err
	}
	var v map[string][]rfkillItem
	err = json.Unmarshal(output, &v)
	if err != nil {
		return nil, err
	}
	return v[""], nil
}

func (d *Device) HasBluetoothDeviceBlocked() (has bool, busErr *dbus.Error) {
	d.service.DelayAutoQuit()
	d.incCallingCount()
	defer d.decCallingCount()

	has, err := d.hasBluetoothDeviceBlocked()
	return has, dbusutil.ToError(err)
}

func (d *Device) hasBluetoothDeviceBlocked() (bool, error) {
	items, err := getRfkillItems()
	if err != nil {
		logger.Warning(err)
		return false, err
	}
	logger.Debug(items)
	for _, item := range items {
		if item.Type == rfkillDeviceTypeBluetooth && item.Soft == "blocked" {
			return true, nil
		}
	}
	return false, nil
}
