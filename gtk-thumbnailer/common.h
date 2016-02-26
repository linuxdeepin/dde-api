/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __COMMON_H__
#define __COMMON_H__

int gtk_thumbnail(const char *theme, const char *dest, const char *bg,
		int width, int height);

int try_init();

#endif
