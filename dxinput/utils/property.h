// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef UTILS_H
#define UTILS_H

#include <X11/Xlib.h>

unsigned char* get_prop(int id, const char* prop, int* nitems);
int set_prop_int(int id, const char* prop, unsigned char* data, int nitems, int bit);
int set_prop_float(int id, const char* prop, unsigned char* data, int nitems);
int set_prop(int id, const char* prop, unsigned char* data, int nitems,
                    Atom type, Atom format);
int enable_left_handed(int id, const char* prop, int enabled);

#endif
