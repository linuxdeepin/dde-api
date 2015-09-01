#ifndef __BUTTON_MAP_H__
#define __BUTTON_MAP_H__

unsigned char* get_button_map(unsigned long xid, const char* name, int* nbuttons);
int set_button_map(unsigned long xid, const char* name,
               unsigned char* map, int nbuttons);

#endif
