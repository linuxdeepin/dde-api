// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef __CURSOR_H__
#define __CURSOR_H__

int init_gtk();
void set_gtk_cursor(char* name);
int set_qt_cursor(const char* name);

#endif
