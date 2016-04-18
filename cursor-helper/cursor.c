/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <gtk/gtk.h>
#include <X11/Xlib.h>

#include "cursor.h"

int
init_gtk()
{
    static gboolean xcb_init = FALSE;

    if (!xcb_init) {
        XInitThreads();

        if (!gtk_init_check(NULL, NULL)) {
            return -1;
        }
    }
    xcb_init = TRUE;
    return 0;
}

void
set_gtk_cursor(char* name)
{
    GtkSettings* s = gtk_settings_get_default();
    g_object_set(G_OBJECT(s), "gtk-cursor-theme-name", name, NULL);
}
