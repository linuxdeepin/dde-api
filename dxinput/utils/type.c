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

static int is_mouse_device(int deviceid);
static int is_touchpad_device(int deviceid);
static int is_touchscreen_device(int deviceid);
static int is_wacom_device(int deviceid);
static int is_keyboard_device(int deviceid);
static XIDeviceInfo* get_xdevice_by_id(int deviceid);

static pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

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

    if (is_keyboard_device(deviceid)) {
        return TYPE_KEYBOARD;
    }


    return TYPE_UNKNOWN;
}

int
is_property_exist(int deviceid, const char* prop)
{
    if (!prop) {
        return 0;
    }

    pthread_mutex_lock(&mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed at check prop exist\n");
        pthread_mutex_unlock(&mutex);
        return 0;
    }

    int nprops = 0;
    Atom *props = XIListProperties(disp, deviceid, &nprops);
    if (!props) {
        XCloseDisplay(disp);
        fprintf(stderr, "List '%d' properties failed\n", deviceid);
        pthread_mutex_unlock(&mutex);
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

    pthread_mutex_unlock(&mutex);

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

static int 
is_keyboard_device(int deviceid)
{
    Display *display;
    int num_devices, i;

    pthread_mutex_lock(&mutex);
    // 打开 X11 显示
    display = XOpenDisplay(NULL);
    if (display == NULL) {
        fprintf(stderr, "Open display failed at check prop exist\n");
        pthread_mutex_unlock(&mutex);
        return 0;
    }

    // 获取所有输入设备
    XIDeviceInfo *devices = XIQueryDevice(display, deviceid, &num_devices);
    if (devices == NULL || num_devices != 1) {
        fprintf(stderr, "Error getting device information.\n");
        pthread_mutex_unlock(&mutex);
        XCloseDisplay(display);
        return 0;
    }

    if(devices[0].use != XISlaveKeyboard)
    {
        fprintf(stderr, "Device is not keyboard.\n");
        pthread_mutex_unlock(&mutex);
        XIFreeDeviceInfo(devices);
        XCloseDisplay(display);
        return 0;
    }

    // 释放设备信息内存
    XIFreeDeviceInfo(devices);

    // 关闭 X11 显示
    XCloseDisplay(display);
    pthread_mutex_unlock(&mutex);

    return 1;
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
