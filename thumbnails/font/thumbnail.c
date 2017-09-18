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
 * Font file thumbnail generator.
 *
 * Reference gnome-font-viewer
 **/
#include <ft2build.h>
#include FT_FREETYPE_H

#include <cairo-ft.h>
#include <glib.h>

#include <math.h>


#define PADDING_VERTICAL 2
#define PADDING_HORIZONTAL 4

#define DEFAULT_THUMB_STR "Aa"

static cairo_t* create_cairo_with_white_bg(gint size);
static void draw_text(cairo_t* cr, FT_Face face, gchar* text, gint size);
static gsize read_file(gchar* file, gchar** contents);
static gdouble calculate_scale(gint width, gint height, gint size);

static FT_Face create_font_face(FT_Library library,
                                gchar* contents, gsize length);
static const gchar* ft_strerror (FT_Error error);
static gboolean check_font_contain_text(FT_Face face, const gchar* text);
static gchar* check_for_ascii_glyph_numbers(FT_Face face,
                                              gboolean* found_ascii);
static void destroy_ft_face(FT_Face face);
static void destroy_ft_library(FT_Library library);
static gchar* build_fallback_thumbstr(FT_Face face);

int
font_thumbnail(char* file, char* dest, int size)
{
        gchar* contents;
        gsize length = read_file(file, &contents);
        if (length == -1) {
                return -1;
        }

        FT_Library library;
        FT_Error error = FT_Init_FreeType(&library);
        if (error) {
                g_free(contents);
                g_printerr("Could not init freetype: %s\n", ft_strerror(error));
                return -1;
        }

        FT_Face face = create_font_face(library, contents, length);
        if (!face) {
                g_free(contents);
                destroy_ft_library(library);
                return -1;
        }

        cairo_t* cr = create_cairo_with_white_bg(size);
        draw_text(cr, face, DEFAULT_THUMB_STR, size);
        cairo_surface_write_to_png(cairo_get_target(cr), dest);
        cairo_destroy(cr);

        // Must free at end, otherwise no text on thumbnail image
        g_free(contents);
        destroy_ft_face(face);
        destroy_ft_library(library);
        return 0;
}

static cairo_t*
create_cairo_with_white_bg(gint size)
{
    cairo_surface_t* surface = cairo_image_surface_create(
        CAIRO_FORMAT_ARGB32,
        size,
        size);
    cairo_t* cr = cairo_create(surface);
    cairo_surface_destroy(surface);

    //background color: white
    cairo_set_source_rgb(cr, 1, 1, 1);
    cairo_paint(cr);

    return cr;
}

static void
draw_text(cairo_t* cr, FT_Face face, gchar* text, gint size)
{
        cairo_font_face_t* font = cairo_ft_font_face_create_for_ft_face(face, 0);
        cairo_set_font_face(cr, font);
        cairo_font_face_destroy(font);
        gint font_size = size - 2 * PADDING_VERTICAL;
        cairo_set_font_size(cr, font_size);

        gchar* str;
        if (check_font_contain_text(face, text)) {
                str = g_strdup(text);
        }else {
                str = build_fallback_thumbstr(face);
        }

        cairo_text_extents_t extents;
        cairo_text_extents(cr, str, &extents);
        gdouble scale = calculate_scale(extents.width, extents.height, size);
        cairo_scale(cr, scale, scale);

        cairo_translate(cr,
                        PADDING_HORIZONTAL - extents.x_bearing + (size - scale * extents.width) / 2.0,
                        PADDING_VERTICAL - extents.y_bearing + (size - scale * extents.height) / 2.0);

        // black
        cairo_set_source_rgba(cr, 0, 0, 0, 1.0);
        cairo_show_text(cr, str);

        g_free(str);
}

static gsize
read_file(gchar* file, gchar** contents)
{
        GError* error = NULL;
        gsize length;

        g_file_get_contents(file, contents, &length, &error);
        if (error) {
                g_printerr("Read '%s' contents failed: %s\n",
                           file, error->message);
                g_error_free(error);
                return -1;
        }

        return length;
}

static gdouble
calculate_scale(gint width, gint height, gint size)
{
        gdouble scale_x, scale_y;

        if (width > (size - 2 * PADDING_HORIZONTAL)) {
                scale_x = (gdouble)(size - 2 * PADDING_HORIZONTAL) / width;
        } else {
                scale_x = 1.0;
        }

        if (height > (size - 2 * PADDING_VERTICAL)) {
                scale_y = (gdouble)(size - 2 * PADDING_VERTICAL) / height;
        } else {
                scale_y = 1.0;
        }

        return MIN(scale_x, scale_y);
}

static FT_Face
create_font_face(FT_Library library, gchar* contents, gsize length)
{
        FT_Face face;
        FT_Error error = FT_New_Memory_Face(library,
                                            (const FT_Byte*)contents,
                                            (FT_Long)length,
                                            0, &face);
        if (error) {
                g_printerr("Create font face failed: %s\n", ft_strerror(error));
                return NULL;
        }

        return face;
}

static gboolean
check_font_contain_text(FT_Face face, const gchar* text)
{
        glong len, idx, map;
        gboolean retval;
        gunichar* str = g_utf8_to_ucs4_fast(text, -1, &len);

        FT_CharMap charmap;
        for (map = 0; map < face->num_charmaps; map++) {
                charmap = face->charmaps[map];
                FT_Set_Charmap(face, charmap);

                retval = TRUE;

                for (idx = 0; idx < len; idx++) {
                        gunichar c = str[idx];
                        if (!FT_Get_Char_Index(face, c)) {
                                retval = FALSE;
                                break;
                        }
                }

                if (retval) {
                        break;
                }
        }

        g_free(str);
        return retval;
}

static gchar*
check_for_ascii_glyph_numbers(FT_Face face,gboolean* found_ascii)
{
        *found_ascii = FALSE;

        GString *ascii_str = g_string_new(NULL);
        GString *str = g_string_new(NULL);
        guint glyph, found = 0;
        gulong c = FT_Get_First_Char(face, &glyph);

        do {
                if (glyph == 65 || glyph == 97) {
                        g_string_append_unichar(ascii_str, (gunichar)c);
                        found++;
                }

                if (found  == 2) {
                        break;
                }

                g_string_append_unichar(str, (gunichar)c);
                c = FT_Get_Next_Char(face, c, &glyph);
        } while (glyph != 0);

        if (found == 2) {
                *found_ascii = TRUE;
                g_string_free(str, TRUE);
                return g_string_free(ascii_str, FALSE);
        } else {
                g_string_free(ascii_str, TRUE);
                return g_string_free(str, FALSE);
        }

}

static gchar*
build_fallback_thumbstr(FT_Face face)
{
        gboolean found_ascii = FALSE;
        gchar* chars = check_for_ascii_glyph_numbers(face, &found_ascii);

        if (found_ascii) {
                return chars;
        }

        gint idx = 0;
        GString* retval = g_string_new(NULL);
        gint total_chars = g_utf8_strlen(chars, -1);

        gchar *ptr, *end;
        while(idx < 2) {
                total_chars = (gint)floor(total_chars / 2.0);
                ptr = g_utf8_offset_to_pointer(chars, total_chars);
                end = g_utf8_find_next_char(ptr, NULL);

                g_string_append_len(retval, ptr, end - ptr);
                idx++;
        }

        return g_string_free(retval, FALSE);
}

/**
 * FT_ERRORS_H is a special header file which is used to define
 * the handling of FT2 enumeration constants, and can also be
 * used to generate error message strings with a small macro trick.
 * much more details in freetype2/fterrors.h.
 **/
static const gchar *
ft_strerror (FT_Error error)
{
#undef __FTERRORS_H__
#define FT_ERRORDEF(e,v,s) case e: return s;
#define FT_ERROR_START_LIST
#define FT_ERROR_END_LIST
        switch (error)
        {
#include FT_ERRORS_H
        default:
                return "unknown";
        }
}

static void
destroy_ft_face(FT_Face face)
{
        FT_Error error = FT_Done_Face(face);
        if (error) {
                g_printerr("Could not finalize font face: %s\n",
                           ft_strerror(error));
        }
}

static void
destroy_ft_library(FT_Library library)
{
        FT_Error error = FT_Done_FreeType(library);
        if (error) {
                g_printerr("Could not finalize freetype library:%s \n",
                           ft_strerror(error));
        }
}
