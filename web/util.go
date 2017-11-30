package web

import "encoding/json"

func ToJson(jsonData interface{}) ([]byte,error ){
	return json.Marshal(jsonData)
}

func ToJsonStr(jsonData interface{}) string {
	if data,err := ToJson(jsonData); err == nil{
		return string(data)
	}
	return ""
}