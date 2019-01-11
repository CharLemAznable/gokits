package gokits

import (
    "encoding/json"
)

func Json(v interface{}) string {
    bytes, err := json.Marshal(v)
    if nil != err {
        return ""
    }
    return string(bytes)
}

func UnJson(str string, v interface{}) interface{} {
    err := json.Unmarshal([]byte(str), v)
    return Condition(nil != err, nil, v)
}
