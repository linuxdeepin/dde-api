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

#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "property.h"

#define MAX_BUF_LEN 1000

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

    if (nitems == NULL) {
        fprintf(stderr, "[get_prop] Invalid item number for %d\n", id);
        return NULL;
    }

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[get_prop] Open display failed for %d\n", id);
        return NULL;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[get_prop] Intern atom %s failed\n", prop);
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
        return NULL;
    }

    *nitems = (int)num_items;
    XCloseDisplay(disp);
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
    Display* disp = XOpenDisplay(NULL);
    if (!disp) {
        fprintf(stderr, "[set_prop_float] open display failed\n");
        return -1;
    }

    Atom type = XInternAtom(disp, "FLOAT", False);
    XCloseDisplay(disp);
    if (type == None) {
        fprintf(stderr, "[set_prop_float] Intern 'FLOAT' atom failed\n");
        return -1;
    }
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

    Display* disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "[set_prop] Open display failed for %d\n", id);
        return -1;
    }

    Atom prop_id = XInternAtom(disp, prop, True);
    if (prop_id == None) {
        XCloseDisplay(disp);
        fprintf(stderr, "[set_prop] Intern atom %s failed\n", prop);
        return -1;
    }

    XIChangeProperty(disp, id, prop_id, type, format,
                     XIPropModeReplace, data, nitems);
    /* XFree(&prop_id); */
    XCloseDisplay(disp);
    return 0;
}
