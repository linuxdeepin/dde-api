/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
