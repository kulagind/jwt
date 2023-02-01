package utils

import "encoding/json"

func UnsafetyToJson(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return data
}

func UnsafetyParseJson(parsedJson []byte, obj *interface{}) {
	err := json.Unmarshal(parsedJson, obj)
	if err != nil {
		panic(err)
	}
}
