/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

#include <gtk/gtk.h>
#include <glib/gprintf.h>
#include <glib.h>
#include <stdlib.h>

void append_page(GtkNotebook * notebook, const char *tab_label)
{
	GtkWidget *tab_header, *label, *close_button, *child;

	tab_header = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 5);
	child = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);

	label = gtk_label_new(tab_label);
	close_button =
	    gtk_button_new_from_icon_name("window-close-symbolic",
					  GTK_ICON_SIZE_MENU);
	gtk_button_set_relief(GTK_BUTTON(close_button), GTK_RELIEF_NONE);
	gtk_box_pack_start(GTK_BOX(tab_header), label, TRUE, TRUE, 0);
	gtk_box_pack_end(GTK_BOX(tab_header), close_button, FALSE, FALSE, 0);

	gtk_notebook_append_page(notebook, child, tab_header);
	gtk_container_child_set(GTK_CONTAINER(notebook), child, "tab-expand",
				TRUE, NULL);

	gtk_widget_show_all(tab_header);
}

void add_icon_button(GtkHeaderBar * header_bar, const char *icon_name)
{
	GtkWidget *button =
	    gtk_button_new_from_icon_name(icon_name, GTK_ICON_SIZE_BUTTON);
	gtk_header_bar_pack_end(GTK_HEADER_BAR(header_bar), button);
}

static void capture(GtkOffscreenWindow * w, GdkEvent * ev, gpointer user_data)
{
	char *dest = (char *)user_data;

	GdkWindow *tmp_window = gtk_widget_get_window(GTK_WIDGET(w));
	cairo_surface_t *tmp_surface =
	    gdk_offscreen_window_get_surface(tmp_window);
	if (!tmp_surface) {
		g_warning("Get offscreen surface failed");
		return;
	}

	cairo_status_t status = cairo_surface_write_to_png(tmp_surface, dest);
	g_printf("write png status: %s\n", cairo_status_to_string(status));
	gtk_main_quit();
}

void gtk_thumbnail(char *theme_name, char *dest, int width, int min_height)
{
	g_setenv("GTK_THEME", theme_name, TRUE);

	GtkWidget *window;
	gboolean initialized = gtk_init_check(NULL, NULL);
	if (!initialized) {
		return;
	}

	window = gtk_offscreen_window_new();
	gtk_window_set_title(GTK_WINDOW(window), "Window");
	gtk_window_set_default_size(GTK_WINDOW(window), width, min_height);

	// header bar
	GtkWidget *header = gtk_header_bar_new();
	gtk_header_bar_set_title(GTK_HEADER_BAR(header), NULL);
	gtk_header_bar_set_show_close_button(GTK_HEADER_BAR(header), TRUE);
	gtk_header_bar_set_has_subtitle(GTK_HEADER_BAR(header), FALSE);

	add_icon_button(GTK_HEADER_BAR(header), "open-menu-symbolic");
	add_icon_button(GTK_HEADER_BAR(header), "system-search-symbolic");

	// notebook pages
	GtkWidget *notebook = gtk_notebook_new();
	append_page(GTK_NOTEBOOK(notebook), "Tab 1");
	append_page(GTK_NOTEBOOK(notebook), "Tab 2");

	GtkWidget *vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
	gtk_box_pack_start(GTK_BOX(vbox), header, TRUE, TRUE, 0);
	gtk_box_pack_start(GTK_BOX(vbox), notebook, TRUE, TRUE, 0);
	gtk_container_add(GTK_CONTAINER(window), vbox);

	g_signal_connect(G_OBJECT(window), "damage-event", G_CALLBACK(capture),
			 dest);

	gtk_widget_show_all(window);
	gtk_main();
}
