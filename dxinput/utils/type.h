/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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

int query_device_type(int deviceid);
int is_property_exist(int deviceid, const char* prop);

#endif
