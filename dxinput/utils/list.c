// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>

#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "list.h"
#include "type.h"
#include "x11_mutex.h"

static int append_device(DeviceInfo** devs, XIDeviceInfo* xinfo, int idx);
static void free_device_info(DeviceInfo* dev);

DeviceInfo*
list_device(int* num)
{
    pthread_mutex_lock(&x11_global_mutex);
    setErrorHandler();

    if (!num) {
        fprintf(stderr, "list_device failed, !num\n");
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed\n");
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    int all_num = 0; 
    XIDeviceInfo* xinfos = XIQueryDevice(disp, XIAllDevices, &all_num);
    XCloseDisplay(disp);
    if (!xinfos) {
        fprintf(stderr, "List xinput device failed\n");
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    int i;
    int j = 0;
    DeviceInfo* devs = NULL;
    for (i = 0; i < all_num; i++) {
        if ((xinfos[i].use != XISlavePointer &&
                xinfos[i].use != XISlaveKeyboard &&
                xinfos[i].use != XIFloatingSlave)) {
            continue;
        }

        if (append_device(&devs, &xinfos[i], j) == -1) {
            continue;
        }

        j++;
    }

    XIFreeDeviceInfo(xinfos);
    *num = j;

    pthread_mutex_unlock(&x11_global_mutex);

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
    if(!devs || !xinfo || !xinfo->name){
        fprintf(stderr, "append_device failed for %d\n", idx);
        return -1;
    }

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
    tmp[idx].ty = query_device_type_unlocked(xinfo->deviceid);

    return 0;
}

static void
free_device_info(DeviceInfo* dev)
{
    if (!dev) {
        return;
    }

    if(dev->name){
        free(dev->name);
        dev->name = NULL;
    }
}
