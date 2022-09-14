// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
