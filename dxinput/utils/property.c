// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdio.h>
#include <limits.h>
#include <pthread.h>

#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "property.h"
#include "type.h"
#include "x11_mutex.h"

/**
 *  The return data type if 'char' must be convert to 'int8_t*'
 * if 'int' must be convert to 'int32_t*'
 * if 'float' must be convert to 'float*'
 *
 * Returns the property data and sets nbytes to the actual byte length.
 * The caller is responsible for calling XFree() on the returned data.
 **/
unsigned char*
get_prop(int id, const char* prop, int* nbytes)
{
    if (!prop) {
        fprintf(stderr, "[get_prop] Empty property for %d\n", id);
        return NULL;
    }

    if (!nbytes) {
        fprintf(stderr, "[get_prop] Invalid nbytes pointer for %d\n", id);
        return NULL;
    }

    *nbytes = 0;

    pthread_mutex_lock(&x11_global_mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[get_prop] Open display failed for %d\n", id);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Intern atom %s failed\n", prop);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    Atom act_type;
    int act_format;
    unsigned long num_items, bytes_after;
    unsigned char* data = NULL;

    // Step 1: Query property size (length=0 to get bytes_after)
    int ret = XIGetProperty(disp, id, prop_id, 0, 0, False,
                            AnyPropertyType, &act_type, &act_format,
                            &num_items, &bytes_after, &data);
    if (ret != Success || act_type == None) {
        if (data) {
            XFree(data);
        }
        XCloseDisplay(disp);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    if (data) {
        XFree(data);
        data = NULL;
    }

    if (bytes_after == 0) {
        // Property exists but has no data
        XCloseDisplay(disp);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    // Step 2: Read all property data
    // length is in 32-bit units, so divide bytes_after by 4 (round up)
    unsigned long length = (bytes_after + 3) / 4;

    ret = XIGetProperty(disp, id, prop_id, 0, length, False,
                        AnyPropertyType, &act_type, &act_format,
                        &num_items, &bytes_after, &data);
    if (ret != Success) {
        if (data) {
            XFree(data);
        }
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Get %s data failed for %d\n", prop, id);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }

    // Calculate actual byte length based on format and num_items
    // format is 8, 16, or 32 bits per item
    // Guard against integer overflow when converting unsigned long to int
    unsigned long nbytes_ul = num_items * (act_format / 8);
    if (nbytes_ul > INT_MAX) {
        if (data) {
            XFree(data);
        }
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Property data too large: %lu bytes (max %d)\n", nbytes_ul, INT_MAX);
        pthread_mutex_unlock(&x11_global_mutex);
        return NULL;
    }
    *nbytes = (int)nbytes_ul;

    XCloseDisplay(disp);
    pthread_mutex_unlock(&x11_global_mutex);

    return data;
}

/**
 * Free the property data returned by get_prop.
 * This is a wrapper for XFree to be called from Go.
 */
void
free_prop_data(unsigned char* data)
{
    if (data) {
        XFree(data);
    }
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
    pthread_mutex_lock(&x11_global_mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(NULL);
    if (!disp) {
        fprintf(stderr, "[set_prop_float] open display failed\n");
        pthread_mutex_unlock(&x11_global_mutex);
        return -1;
    }

    Atom type = XInternAtom(disp, "FLOAT", False);
    XCloseDisplay(disp);
    if (type == None) {
        fprintf(stderr, "[set_prop_float] Intern 'FLOAT' atom failed\n");
        pthread_mutex_unlock(&x11_global_mutex);
        return -1;
    }

    pthread_mutex_unlock(&x11_global_mutex);

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

    pthread_mutex_lock(&x11_global_mutex);
    setErrorHandler();

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[set_prop] Open display failed for %d\n", id);
        pthread_mutex_unlock(&x11_global_mutex);
        return -1;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[set_prop] Intern atom %s failed\n", prop);
        pthread_mutex_unlock(&x11_global_mutex);
        return -1;
    }

    XIChangeProperty(disp, id, prop_id, type, format,
                     XIPropModeReplace, data, nitems);
    /* XFree(&prop_id); */
    XCloseDisplay(disp);

    pthread_mutex_unlock(&x11_global_mutex);

    return 0;
}
