// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <gdk/gdk.h>
#include "text.h"

#define FONT_NAME "monospace"

static void do_show_text(cairo_t* cr, char** text, ThumbInfo* info);

int
text_thumbnail(char** text, char* dest, ThumbInfo* info)
{
	if (!gdk_init_check(NULL, NULL)) {
		g_warning("Init gdk failed");
		return -1;
	}

	cairo_surface_t* surface = cairo_image_surface_create(
	    CAIRO_FORMAT_ARGB32, info->width, info->height);
	if (!surface) {
		g_warning("Create surface failed");
		return -1;
	}

	cairo_t* cr = cairo_create(surface);
	cairo_surface_destroy(surface);
	if (!cr) {
		g_warning("Create cairo failed");
		return -1;
	}

	cairo_set_source_rgba(cr, 1.0, 1.0, 1.0, 0);
	cairo_paint(cr);
	do_show_text(cr, text, info);

	cairo_status_t status = cairo_surface_write_to_png(
			cairo_get_target(cr),
			dest);
	cairo_destroy(cr);
	if (status != CAIRO_STATUS_SUCCESS) {
		g_warning("Write cairo to file '%s' failed", dest);
		return -1;
	}

	return 0;
}

static void
do_show_text(cairo_t* cr, char** text, ThumbInfo* info)
{
	cairo_select_font_face(cr, FONT_NAME,
	                       CAIRO_FONT_SLANT_NORMAL,
	                       CAIRO_FONT_WEIGHT_BOLD);

	cairo_set_font_size(cr, info->fontSize);

	// text color: black
	cairo_set_source_rgb(cr, 0, 0, 0);

	int i = 0;
	int y_pos = info->yborder;
	while (text[i]) {
		if (y_pos > info->canvasHeight) {
			break;
		}

		cairo_move_to(cr, info->xborder, y_pos);
		cairo_show_text(cr, text[i]);
		y_pos += info->fontSize;
		i++;
	}
}
