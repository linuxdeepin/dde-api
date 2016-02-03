/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

/**
 * Set keyboard repeat
 **/

#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/XKBlib.h>

int
set_keyboard_repeat(int repeated, unsigned int delay, unsigned int interval)
{
    Display *disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed\n");
        return -1;
    }

    int ret = 0;
    if (repeated) {
        XAutoRepeatOn(disp);

        // Use XKB in preference
        int rate_set = XkbSetAutoRepeatRate(disp, XkbUseCoreKbd,
                                            delay, interval);
        if (!rate_set) {
            ret = -1;
            fprintf(stderr, "Neither XKeyboard not Xfree86's\
				       	keyboard extensions are available,\
					\n no way to support keyboard\
				       	autorepeat rate settings\n");
        }
    } else {
        XAutoRepeatOff(disp);
    }

    XSync(disp, False);
    XCloseDisplay(disp);

    return ret;
}
