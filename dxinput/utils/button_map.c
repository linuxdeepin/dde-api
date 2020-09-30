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
#include <stdlib.h>
#include <string.h>

#include <X11/Xlib.h>
#include <X11/extensions/XInput.h>

#include "button_map.h"

static int get_button_number(Display* disp, const char* name);
static const XDeviceInfo* find_device_by_name(const XDeviceInfo* devs,
                                              int n_dev, const char* name);
static int get_device_button_number(const XDeviceInfo* dev);
static unsigned char* do_get_button_map(Display* disp,
                                        unsigned long xid, int nbuttons);

unsigned char*
get_button_map(unsigned long xid, const char* name, int* nbuttons)
{
    if (!name) {
        fprintf(stderr, "[get_button_map] empty device name for %lu\n",
                xid);
        return NULL;
    }

    Display* disp = XOpenDisplay(NULL);
    if (!disp) {
        fprintf(stderr, "[get_button_map] open display failed for %lu %s\n",
                xid, name);
        return NULL;
    }

    int num_btn = get_button_number(disp, name);
    if (num_btn == -1) {
        XCloseDisplay(disp);
        fprintf(stderr, "[get_button_map] get button number failed for %lu %s\n",
                xid, name);
        return NULL;
    }

    *nbuttons = num_btn;
    unsigned char* map = do_get_button_map(disp, xid, num_btn);
    XCloseDisplay(disp);
    return map;
}

int
set_button_map(unsigned long xid, const char* name,
               unsigned char* map, int nbuttons)
{
    if (!name || !map) {
        fprintf(stderr, "[set_button_map] invalid device name or value\n");
        return -1;
    }

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[set_button_map] open display failed: %lu %s\n",
                xid, name);
        return -1;
    }

    XDevice* dev = XOpenDevice(disp, xid);
    if (!dev) {
        XCloseDisplay(disp);
        fprintf(stderr, "[set_button_map] open device failed for %lu %s\n",
                xid, name);
        return -1;
    }

    // map: no two elements can have the same nonzero value,
    // or a BadValue error results.
    int ret = XSetDeviceButtonMapping(disp, dev, map, nbuttons);
    XCloseDevice(disp, dev);
    XCloseDisplay(disp);
    // TODO: if ret == MappingBusy, try again
    if (ret != MappingSuccess) {
        return -1;
    }

    return 0;
}

static unsigned char*
do_get_button_map(Display* disp, unsigned long xid, int nbuttons)
{
    XDevice* dev = XOpenDevice(disp, xid);
    if (!dev) {
        fprintf(stderr, "[do_get_button_map] open device %lu failed\n", xid);
        return NULL;
    }

    unsigned char* map = (unsigned char*)calloc(nbuttons,
                                                sizeof(unsigned char));
    if (!map) {
        XCloseDevice(disp, dev);
        fprintf(stderr, "[do_get_button_map] alloc memory failed for %lu\n",
                xid);
        return NULL;
    }

    int rc = XGetDeviceButtonMapping(disp, dev, map, nbuttons);
    XCloseDevice(disp, dev);
    if (rc != nbuttons) {
        free(map);
        fprintf(stderr, "[do_get_button_map] get button map failed for %lu\n",
                xid);
        return NULL;
    }

    return map;
}

static int
get_button_number(Display* disp, const char* name)
{
    int n_dev = 0;
    XDeviceInfo* devs = XListInputDevices(disp, &n_dev);
    if (!devs) {
        fprintf(stderr, "[get_button_number] list device failed for %s\n", name);
        return -1;
    }

    const XDeviceInfo* info = find_device_by_name(devs, n_dev, name);
    if (!info) {
        XFreeDeviceList(devs);
        fprintf(stderr, "[get_button_number] not found device for %s\n", name);
        return -1;
    }

    int num_btn = get_device_button_number(info);
    XFreeDeviceList(devs);
    return num_btn;
}

static const XDeviceInfo*
find_device_by_name(const XDeviceInfo* devs, int n_dev, const char* name)
{
    int i = 0;
    for (; i < n_dev; i++) {
        if (devs[i].use != IsXExtensionPointer) {
            continue;
        }

        if (strcmp(name, devs[i].name) == 0) {
            return &(devs[i]);
        }
    }
    return NULL;
}

static int
get_device_button_number(const XDeviceInfo* dev)
{
    if (!dev) {
        return -1;
    }

    int i = 0;
    int num_btn = -1;
    XAnyClassPtr ptr = (XAnyClassPtr)dev->inputclassinfo;
    for (; i < dev->num_classes; i++) {
        if (ptr->class != ButtonClass) {
            ptr = (XAnyClassPtr)((char*)ptr + ptr->length);
            continue;
        }

        num_btn = ((XButtonInfoPtr)ptr)->num_buttons;
        break;
    }
    return num_btn;
}
