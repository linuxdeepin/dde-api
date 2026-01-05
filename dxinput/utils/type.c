// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdio.h>
#include <string.h>
#include <pthread.h>

#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "type.h"
#include "x11_mutex.h"

static int is_mouse_device(int deviceid);
static int is_touchpad_device(int deviceid);
static int is_touchscreen_device(int deviceid);
static int is_wacom_device(int deviceid);
static int is_keyboard_device(int deviceid);
static XIDeviceInfo* get_xdevice_by_id(int deviceid);
static int is_property_exist_unlocked(int deviceid, const char* prop);

int
listener_error_handler(Display * display, XErrorEvent * event)
{
    if(display && event){
        char msg[256];
        XGetErrorText(display, event->error_code, msg, 255);
        fprintf(stderr, "Ignore Xlib error : %s\n", msg);
    } else{
        fprintf(stderr, "listener_error_handler error\n");
    }
    return 0;
}

int
listener_ioerror_handler(Display * display)
{
    return 0;
}

void
setErrorHandler(){
    XSetErrorHandler(listener_error_handler);
    XSetIOErrorHandler(listener_ioerror_handler);
}

// Internal version without locking - assumes caller holds x11_global_mutex
int
query_device_type_unlocked(int deviceid)
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

    if (is_keyboard_device(deviceid)) {
        return TYPE_KEYBOARD;
    }

    return TYPE_UNKNOWN;
}

// External API - with locking
int
query_device_type(int deviceid)
{
    pthread_mutex_lock(&x11_global_mutex);
    
    int result = query_device_type_unlocked(deviceid);

    pthread_mutex_unlock(&x11_global_mutex);
    return result;
}

// Internal version without locking - assumes caller holds x11_global_mutex
static int
is_property_exist_unlocked(int deviceid, const char* prop)
{
    if (!prop) {
        return 0;
    }

    setErrorHandler();

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
        if (name && strcmp(name, prop) == 0) {
            exist = 1;
        }
        if (name) {
            XFree(name);
        }

        if (exist == 1) {
            break;
        }
    }
    XCloseDisplay(disp);
    XFree(props);

    return exist;
}

// External API - with locking
int
is_property_exist(int deviceid, const char* prop)
{
    pthread_mutex_lock(&x11_global_mutex);
    int result = is_property_exist_unlocked(deviceid, prop);
    pthread_mutex_unlock(&x11_global_mutex);
    return result;
}

static int
is_mouse_device(int deviceid)
{
    return (is_property_exist_unlocked(deviceid, "Button Labels") ||
            is_property_exist_unlocked(deviceid, "libinput Button Scrolling Button"));
}

static int
is_touchpad_device(int deviceid)
{
    return (is_property_exist_unlocked(deviceid, "Synaptics Off") ||
            is_property_exist_unlocked(deviceid, "libinput Tapping Enabled"));
}

static int 
is_keyboard_device(int deviceid)
{
    Display *display;
    int num_devices, i;

    // NOTE: No mutex lock here - called from query_device_type which already holds the lock
    // 打开 X11 显示
    display = XOpenDisplay(NULL);
    if (display == NULL) {
        fprintf(stderr, "Open display failed at check prop exist\n");
        return 0;
    }

    // 获取所有输入设备
    XIDeviceInfo *devices = XIQueryDevice(display, deviceid, &num_devices);
    if (devices == NULL || num_devices != 1) {
        fprintf(stderr, "Error getting device information.\n");
        XCloseDisplay(display);
        return 0;
    }

    if(devices[0].use != XISlaveKeyboard)
    {
        fprintf(stderr, "Device is not keyboard.\n");
        XIFreeDeviceInfo(devices);
        XCloseDisplay(display);
        return 0;
    }

    // 释放设备信息内存
    XIFreeDeviceInfo(devices);

    // 关闭 X11 显示
    XCloseDisplay(display);

    return 1;
}

// TODO: support libinput
static int
is_wacom_device(int deviceid)
{
    return is_property_exist_unlocked(deviceid, "Wacom Tool Type");
}

static int
is_touchscreen_device(int deviceid)
{
    // for libinput
    if (is_property_exist_unlocked(deviceid, "libinput Calibration Matrix")) {
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
