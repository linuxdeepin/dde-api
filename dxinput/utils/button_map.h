/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __BUTTON_MAP_H__
#define __BUTTON_MAP_H__

unsigned char* get_button_map(unsigned long xid, const char* name, int* nbuttons);
int set_button_map(unsigned long xid, const char* name,
               unsigned char* map, int nbuttons);

#endif
