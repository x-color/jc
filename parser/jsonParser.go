package parser

import (
	"encoding/json"
)

func ParseJSON(rowJSONData []byte) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal(rowJSONData, &data); err != nil {
		return nil, err
	}
	return data, nil
}
