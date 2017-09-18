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
