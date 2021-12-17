package util

import "encoding/json"

func StructToString(in interface{}) string {
	bytes, _ := json.Marshal(in)
	return string(bytes)
}
