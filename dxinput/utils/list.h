#ifndef LIST_H
#define LIST_H

typedef struct _DeviceInfo {
    char *name;
    int id;
    int enabled;
    int ty; // type
} DeviceInfo;

DeviceInfo* list_device(int* num);
void free_device_list(DeviceInfo* devs, int num);

#endif
