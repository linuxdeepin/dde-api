/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __RECORD_H__
#define __RECORD_H__

#define KEY_PRESS 1
#define KEY_RELEASE 0
#define BUTTON_PRESS 1
#define BUTTON_RELEASE 0

#include <X11/extensions/XInput2.h>

int start_listen();

int xi_mask_is_set(unsigned char*ptr, char mask);

#endif
