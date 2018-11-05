/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

#include <stdio.h>
#include <string.h>

#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "type.h"

static int is_mouse_device(int deviceid);
static int is_touchpad_device(int deviceid);
static int is_touchscreen_device(int deviceid);
static int is_wacom_device(int deviceid);
static XIDeviceInfo* get_xdevice_by_id(int deviceid);

int
query_device_type(int deviceid)
{
    if (is_wacom_device(deviceid)) {
        return TYPE_WACOM;
    }

    if (is_touchscreen_device(deviceid)) {
        return TYPE_TOUCHSCREEN;
    }

    if (is_touchpad_device(deviceid)) {
        return TYPE_TOUCHPAD;
    }

    if (is_mouse_device(deviceid)) {
        return TYPE_MOUSE;
    }

    return TYPE_UNKNOWN;
}

int
is_property_exist(int deviceid, const char* prop)
{
    if (!prop) {
        return 0;
    }

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed at check prop exist\n");
        return 0;
    }

    int nprops = 0;
    Atom *props = XIListProperties(disp, deviceid, &nprops);
    if (!props) {
        XCloseDisplay(disp);
        fprintf(stderr, "List '%d' properties failed\n", deviceid);
        return 0;
    }

    int exist = 0;
    while (nprops--) {
        char* name = XGetAtomName(disp, props[nprops]);
        if (strcmp(name, prop) == 0) {
            exist = 1;
        }
        XFree(name);

        if (exist == 1) {
            break;
        }
    }
    XCloseDisplay(disp);
    XFree(props);

    return exist;
}

static int
is_mouse_device(int deviceid)
{
    return (is_property_exist(deviceid, "Button Labels") ||
            is_property_exist(deviceid, "libinput Button Scrolling Button"));
}

static int
is_touchpad_device(int deviceid)
{
    return (is_property_exist(deviceid, "Synaptics Off") ||
            is_property_exist(deviceid, "libinput Tapping Enabled"));
}

// TODO: support libinput
static int
is_wacom_device(int deviceid)
{
    return is_property_exist(deviceid, "Wacom Tool Type");
}

static int
is_touchscreen_device(int deviceid)
{
    // for libinput
    if (is_property_exist(deviceid, "libinput Calibration Matrix")) {
        return 1;
    }
    // Now XInput2 library detect touchscreen as mouse
    if (!is_mouse_device(deviceid)) {
        return 0;
    }

    XIDeviceInfo* xinfo = get_xdevice_by_id(deviceid);
    if (!xinfo) {
        return 0;
    }

    if (xinfo->num_classes <= 0) {
        XIFreeDeviceInfo(xinfo);
        return 0;
    }

    int i = 0;
    int flag = 0;
    for (; i < xinfo->num_classes; i++) {
        XIAnyClassInfo* any = xinfo->classes[i];
        switch (any->type) {
        case XIValuatorClass: {
            XIValuatorClassInfo* val = (XIValuatorClassInfo*)any;
            // Absolute mode is touchscreen, relative mode not
            if (val->mode == XIModeAbsolute) {
                flag = 1;
                break;
            }
        }
        }
    }

    XIFreeDeviceInfo(xinfo);
    return flag;
}

static XIDeviceInfo*
get_xdevice_by_id(int deviceid)
{
    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed at query device\n");
        return NULL;
    }

    int num = 0;
    XIDeviceInfo* xinfo = XIQueryDevice(disp, deviceid, &num);
    XCloseDisplay(disp);
    if (!xinfo) {
        fprintf(stderr, "Query device by id '%d' failed\n", deviceid);
        return NULL;
    }

    if (num != 1) {
        XIFreeDeviceInfo(xinfo);
        return NULL;
    }

    return xinfo;
}
