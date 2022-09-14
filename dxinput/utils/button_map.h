// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef __BUTTON_MAP_H__
#define __BUTTON_MAP_H__

unsigned char* get_button_map(unsigned long xid, const char* name, int* nbuttons);
int set_button_map(unsigned long xid, const char* name,
               unsigned char* map, int nbuttons);

#endif
