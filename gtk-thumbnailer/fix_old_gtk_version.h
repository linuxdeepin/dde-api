/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <gtk/gtk.h>
#ifndef __CREATE_FROM_PIXBUF__
#define __CREATE_FROM_PIXBUF__


#if !GTK_CHECK_VERSION(3, 10, 0)
cairo_surface_t* gdk_cairo_surface_create_from_pixbuf(GdkPixbuf* pixbuf, int scale, GdkWindow* w);
#endif

#if !GTK_CHECK_VERSION(3, 8, 0)
void gtk_widget_set_opacity(GtkWidget* widget, double opacity);
#endif

#endif
