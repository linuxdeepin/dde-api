/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef LIST_H
#define LIST_H

typedef struct _DeviceInfo {
    char *name;
    int id;
    int enabled;
    int ty; // type
} DeviceInfo;

DeviceInfo* list_device(int* num);
void free_device_list(DeviceInfo* devs, int num);

#endif
