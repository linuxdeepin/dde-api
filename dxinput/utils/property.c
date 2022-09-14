// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdio.h>
#include <pthread.h>

#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "property.h"
#include "type.h"

#define MAX_BUF_LEN 1000

static pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

/**
 *  The return data type if 'char' must be convert to 'int8_t*'
 * if 'int' must be convert to 'int32_t*'
 * if 'float' must be convert to 'float*'
 **/
unsigned char*
get_prop(int id, const char* prop, int* nitems)
{
    if (!prop) {
        fprintf(stderr, "[get_prop] Empty property for %d\n", id);
        return NULL;
    }

    if (!nitems) {
        fprintf(stderr, "[get_prop] Invalid item number for %d\n", id);
        return NULL;
    }

    pthread_mutex_lock(&mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[get_prop] Open display failed for %d\n", id);
        pthread_mutex_unlock(&mutex);
        return NULL;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Intern atom %s failed\n", prop);
        pthread_mutex_unlock(&mutex);
        return NULL;
    }

    Atom act_type;
    int act_format;
    unsigned long num_items, bytes_after;
    unsigned char* data = NULL;
    int ret = XIGetProperty(disp, id, prop_id, 0, MAX_BUF_LEN, False,
                            AnyPropertyType, &act_type, &act_format,
                            &num_items, &bytes_after, &data);
    if (ret != Success) {
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Get %s data failed for %d\n", prop, id);
        pthread_mutex_unlock(&mutex);
        return NULL;
    }

    *nitems = (int)num_items;
    XCloseDisplay(disp);

    pthread_mutex_unlock(&mutex);

    return data;
}

// bit: range(8,16,32)
int
set_prop_int(int id, const char* prop, unsigned char* data, int nitems, int bit)
{
    return set_prop(id, prop, data, nitems, XA_INTEGER, bit);
}

int
set_prop_float(int id, const char* prop, unsigned char* data, int nitems)
{
    pthread_mutex_lock(&mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(NULL);
    if (!disp) {
        fprintf(stderr, "[set_prop_float] open display failed\n");
        pthread_mutex_unlock(&mutex);
        return -1;
    }

    Atom type = XInternAtom(disp, "FLOAT", False);
    XCloseDisplay(disp);
    if (type == None) {
        fprintf(stderr, "[set_prop_float] Intern 'FLOAT' atom failed\n");
        pthread_mutex_unlock(&mutex);
        return -1;
    }

    pthread_mutex_unlock(&mutex);

    // Format must be 32
    int ret = set_prop(id, prop, data, nitems, type, 32);

    return ret;
}

int
set_prop(int id, const char* prop, unsigned char* data, int nitems,
         Atom type, Atom format)
{
    if (!prop) {
        fprintf(stderr, "[set_prop] Empty property for %d\n", id);
        return -1;
    }

    if (!data || nitems < 1) {
        fprintf(stderr, "[set_prop] Invalid data or item number for %d\n", id);
        return -1;
    }

    pthread_mutex_lock(&mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[set_prop] Open display failed for %d\n", id);
        pthread_mutex_unlock(&mutex);
        return -1;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[set_prop] Intern atom %s failed\n", prop);
        pthread_mutex_unlock(&mutex);
        return -1;
    }

    XIChangeProperty(disp, id, prop_id, type, format,
                     XIPropModeReplace, data, nitems);
    /* XFree(&prop_id); */
    XCloseDisplay(disp);

    pthread_mutex_unlock(&mutex);

    return 0;
}
