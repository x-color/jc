package parser

import (
	"errors"
	"strconv"
	"strings"
)

type Key struct {
	name    string
	number  int
	isSlice bool
}

func ParseKeys(rowKeys []string) ([]Key, error) {
	keys := []Key{}
	for _, rowKey := range rowKeys {
		if strings.HasSuffix(rowKey, "]") {
			data := strings.Split(strings.TrimSuffix(rowKey, "]"), "[")
			if len(data) != 2 {
				return nil, errors.New("can't parse args")
			}
			name := data[0]
			num, err := strconv.Atoi(data[1])
			if err != nil {
				return nil, err
			}
			if num < 0 {
				return nil, errors.New("index number is required more than 0")
			}
			if name != "" {
				keys = append(keys, Key{name: name, isSlice: false})
			}
			keys = append(keys, Key{number: num, isSlice: true})
		} else {
			keys = append(keys, Key{name: rowKey, isSlice: false})
		}
	}
	return keys, nil
}

func ChoiceFromJSON(jsonData interface{}, keys []Key) (interface{}, error) {
	data := jsonData
	for _, key := range keys {
		switch convertedData := data.(type) {
		case []interface{}:
			if !key.isSlice {
				return nil, errors.New("index number is required, not key")
			}
			if len(convertedData) <= key.number {
				return nil, errors.New("index number is out of range")
			}
			data = convertedData[key.number]
		case map[string]interface{}:
			var ok bool
			data, ok = convertedData[key.name]
			if !ok {
				return nil, errors.New("key does not exist")
			}
		default:
			return nil, errors.New("can't access JSON object")
		}
	}
	return data, nil
}
