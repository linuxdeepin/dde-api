/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef UTILS_H
#define UTILS_H

#include <X11/Xlib.h>

unsigned char* get_prop(int id, const char* prop, int nitems);
int set_prop_bool(int id, const char* prop, unsigned char* data, int nitems);
int set_prop_int32(int id, const char* prop, unsigned char* data, int nitems);
int set_prop_float(int id, const char* prop, unsigned char* data, int nitems);
int set_prop(int id, const char* prop, unsigned char* data, int nitems,
                    Atom type, Atom format);
int enable_left_handed(int id, const char* prop, int enabled);

#endif
