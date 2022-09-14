// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
