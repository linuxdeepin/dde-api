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
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "list.h"
#include "type.h"

static int append_device(DeviceInfo** devs, XIDeviceInfo* xinfo, int idx);
static void free_device_info(DeviceInfo* dev);

DeviceInfo*
list_device(int* num)
{
    if (!num) {
        return NULL;
    }

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed\n");
        return NULL;
    }

    int all_num;
    XIDeviceInfo* xinfos = XIQueryDevice(disp, XIAllDevices, &all_num);
    XCloseDisplay(disp);
    if (!xinfos) {
        fprintf(stderr, "List xinput device failed\n");
        return NULL;
    }

    int i;
    int j = 0;
    DeviceInfo* devs = NULL;
    for (i = 0; i < all_num; i++) {
        if ((xinfos[i].use != XISlavePointer &&
                xinfos[i].use != XISlaveKeyboard)) {
            continue;
        }

        if (append_device(&devs, &xinfos[i], j) == -1) {
            continue;
        }

        j++;
    }

    XIFreeDeviceInfo(xinfos);
    *num = j;

    return devs;
}

void
free_device_list(DeviceInfo* devs, int num)
{
    if (!devs) {
        return;
    }

    int i = 0;
    for (; i < num; i++) {
        free_device_info(&devs[i]);
    }

    free(devs);
}

static int
append_device(DeviceInfo** devs, XIDeviceInfo* xinfo, int idx)
{
    unsigned long size = strlen(xinfo->name);
    char* name = (char*)calloc(size+1, sizeof(char));
    if (!name) {
        fprintf(stderr, "Alloc info name memory failed for %d\n", idx);
        return -1;
    }
    memcpy(name, xinfo->name, size);

    DeviceInfo* tmp = (DeviceInfo*)realloc(*devs, (idx+1)*sizeof(DeviceInfo));
    if (!tmp) {
        fprintf(stderr, "Alloc memory failed for '%d' DeviceInfo\n", idx+1);
        free(name);
        return -1;
    }

    // if *devs == NULL
    *devs = tmp;
    // construct 'DeviceInfo' from 'XIDeviceInfo'
    tmp[idx].name = name;
    tmp[idx].id = xinfo->deviceid;
    tmp[idx].enabled = xinfo->enabled;
    tmp[idx].ty = query_device_type(xinfo->deviceid);

    return 0;
}

static void
free_device_info(DeviceInfo* dev)
{
    if (!dev) {
        return;
    }

    free(dev->name);
    dev->name = NULL;
}
