#ifndef TYPE_H
#define TYPE_H

enum DEVICE_TYPE {
    TYPE_UNKNOWN,
    TYPE_KEYBOARD,
    TYPE_MOUSE,
    TYPE_TOUCHPAD,
    TYPE_WACOM,
    TYPE_TOUCHSCREEN,
};

int query_device_type(int deviceid);
int is_property_exist(int deviceid, const char* prop);

#endif
