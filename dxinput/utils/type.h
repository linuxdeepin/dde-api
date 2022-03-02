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

#ifndef TYPE_H
#define TYPE_H

enum DEVICE_TYPE {
    TYPE_UNKNOWN,
    TYPE_KEYBOARD,
    TYPE_MOUSE,
    TYPE_TOUCHPAD,
    TYPE_WACOM,
    TYPE_TOUCHSCREEN,
};

void setErrorHandler();
int listener_error_handler(Display * display, XErrorEvent * event);
int listener_ioerror_handler(Display * display);

int query_device_type(int deviceid);
int is_property_exist(int deviceid, const char* prop);

#endif
