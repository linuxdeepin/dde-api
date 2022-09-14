// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
