//gcc theme_preview.c `pkg-config --libs --cflags gtk+-2.0  libmetacity-private `
//
/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/


#include <gtk/gtk.h>
#include <metacity-private/theme-parser.h>
#include <metacity-private/preview-widget.h>

#include "fix_old_gtk_version.h"

typedef struct _ThumbData {
	int width;
	int height;

	char* dest;
	char* background;
} ThumbData;

static GtkWidget* get_preview_from_meta(const char* name);
static void padding_thumbnail(const GtkFixed* fixed);
static void capture(GtkOffscreenWindow* w, GdkEvent* ev, gpointer user_data);

int
try_init()
{
    return gtk_init_check(NULL, NULL);
}

int
generate_thumbnail(const char* name, const char* dest, const char* bg,
                   int width, int height)
{
	if (!name || !dest) {
		g_warning("Invalid theme name or dest");
		return -1;
	}

	GtkWidget* w = gtk_offscreen_window_new();
	gtk_widget_set_size_request(w, width, height);
	GtkWidget* preview = get_preview_from_meta(name);
	if (!preview) {
		g_warning("get metacity theme preview failed");
		return -1;
	}

	gtk_container_add(GTK_CONTAINER(w), preview);
	GtkWidget* fixed = gtk_fixed_new();
	gtk_container_add(GTK_CONTAINER(preview), fixed);
	padding_thumbnail(GTK_FIXED(fixed));

	ThumbData data;
	data.width = width;
	data.height = height;
	data.dest = (char*)dest;
	data.background = (char*)bg;
	g_signal_connect(G_OBJECT(w), "damage-event",
	                 G_CALLBACK(capture), &data);
	gtk_widget_realize(fixed);
	gtk_widget_show_all(w);

	gtk_main();
	return 0;
}

static GtkWidget*
get_preview_from_meta(const char* name)
{
	if (!name) {
		g_warning("Theme name is null");
		return NULL;
	}

	// Init meta_current_theme, otherwise segmentation in metacity
	meta_theme_set_current("", TRUE);

	GError* error = NULL;
	MetaTheme* meta = NULL;
	meta = meta_theme_load(name, &error);
	if (error) {
		g_warning("Load meta theme '%s' failed: %s",
		          name, error->message);
		g_error_free(error);
		return NULL;
	}

	GtkWidget* preview = NULL;
	preview = meta_preview_new();
	if (!preview) {
		g_warning("New metacity preview failed");
		return NULL;
	}

	meta_preview_set_theme((MetaPreview*)preview, meta);
	/*meta_theme_free(meta);*/
	meta_preview_set_title((MetaPreview*)preview, "");

	return preview;
}

static void
padding_thumbnail(const GtkFixed* fixed)
{
	//TODO: Should handle gtk2/gtk3 themes
}

static void
capture(GtkOffscreenWindow* w, GdkEvent* ev, gpointer user_data)
{
	ThumbData* data = (ThumbData*)user_data;
	int width = data->width;
	int height = data->height;
	char* dest = data->dest;
	char* bg = data->background;

	cairo_surface_t* surface = NULL;
	if (bg) {
		surface = cairo_image_surface_create_from_png(bg);
	} else {
		surface = cairo_image_surface_create(
		              CAIRO_FORMAT_ARGB32,
		              width, height);
	}
	if (!surface) {
		g_warning("Create surface failed\n");
		return;
	}

	cairo_t* cairo = cairo_create(surface);
	GdkPixbuf* pbuf = gtk_offscreen_window_get_pixbuf(w);

	gdk_cairo_set_source_pixbuf(cairo, pbuf, -15, 15);
	cairo_paint(cairo);
	cairo_surface_write_to_png(surface, dest);

	g_object_unref(G_OBJECT(pbuf));
	cairo_destroy(cairo);
	cairo_surface_destroy(surface);

	gtk_main_quit();
}
