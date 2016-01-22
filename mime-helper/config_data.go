package main

import (
	"encoding/json"
	"io/ioutil"
)

type defaultAppTable struct {
	Apps defaultAppInfos `json:"DefaultApps"`
}

type defaultAppInfo struct {
	AppId   string   `json:"AppId"`
	AppType string   `json:"AppType"`
	Types   []string `json:"SupportedType"`
}
type defaultAppInfos []*defaultAppInfo

func unmarshal(file string) (*defaultAppTable, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var table defaultAppTable
	err = json.Unmarshal(content, &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func marshal(v interface{}) (string, error) {
	content, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func genMimeAppsFile(data string) error {
	table, err := unmarshal(data)
	if err != nil {
		logger.Warning("[genMimeAppsFile] unmarshal failed:", err)
		return err
	}

	for _, info := range table.Apps {
		for _, ty := range info.Types {
			err := SetAppInfo(ty, info.AppId)
			if err != nil {
				logger.Warningf("[genMimeAppsFile] set '%s' to parse '%s' failed: %v\n",
					info.AppId, ty, err)
			}
		}
	}

	return nil
}
