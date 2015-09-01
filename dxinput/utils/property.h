#ifndef UTILS_H
#define UTILS_H

#include <X11/Xlib.h>

unsigned char* get_prop(int id, const char* prop, int nitems);
int set_prop_bool(int id, const char* prop, unsigned char* data, int nitems);
int set_prop_int32(int id, const char* prop, unsigned char* data, int nitems);
int set_prop_float(int id, const char* prop, unsigned char* data, int nitems);
int set_prop(int id, const char* prop, unsigned char* data, int nitems,
                    Atom type, Atom format);
int enable_left_handed(int id, const char* prop, int enabled);

#endif
