/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

func (op *cacheUtil) setCurrent(year int32) {
	if op.current != year {
		op.current = year
		op.clearCache()
	}
}

func (op *cacheUtil) setCache(key string, value interface{}) interface{} {
	if cacheMap == nil {
		cacheMap = make(map[string]interface{})
	}
	cacheMap[key] = value
	return value
}

func (op *cacheUtil) getCache(key string) (interface{}, bool) {
	if cacheMap == nil {
		cacheMap = make(map[string]interface{})
	}

	v, ok := cacheMap[key]
	return v, ok
}

func (op *cacheUtil) clearCache() {
	for key, _ := range cacheMap {
		delete(cacheMap, key)
	}
}

func newCache() *cacheUtil {
	m := &cacheUtil{}
	m.current = 0

	return m
}
