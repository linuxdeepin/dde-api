/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

#ifndef __RECORD_H__
#define __RECORD_H__

#define KEY_PRESS 1
#define KEY_RELEASE 0
#define BUTTON_PRESS 1
#define BUTTON_RELEASE 0

#include <X11/extensions/XInput2.h>

int start_listen();

int xi_mask_is_set(unsigned char*ptr, char mask);

//remove this;
void print_deviceevent(XIDeviceEvent* event);
void print_devicechangedevent(Display *dpy, XIDeviceChangedEvent *event);
void print_hierarchychangedevent(XIHierarchyEvent *event);
void print_rawevent(XIRawEvent *event);
void print_enterleave(XILeaveEvent* event);
void print_propertyevent(Display *display, XIPropertyEvent* event);
const char* type_to_name(int evtype);

#endif
