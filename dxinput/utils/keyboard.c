// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

/**
 * Set keyboard repeat
 **/

#include <stdio.h>
#include <pthread.h>

#include <X11/Xlib.h>
#include <X11/XKBlib.h>
#include "type.h"
#include "x11_mutex.h"

int
set_keyboard_repeat(int repeated, unsigned int delay, unsigned int interval)
{
    pthread_mutex_lock(&x11_global_mutex);
    setErrorHandler();

    Display *disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open display failed\n");
        pthread_mutex_unlock(&x11_global_mutex);
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

    pthread_mutex_unlock(&x11_global_mutex);

    return ret;
}
