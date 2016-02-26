/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __CURSOR_H__
#define __CURSOR_H__

int init_gtk();
void set_gtk_cursor(char* name);
int set_qt_cursor(const char* name);

#endif
