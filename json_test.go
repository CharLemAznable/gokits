package gokits

import (
    "testing"
)

type testJsonType struct {
    IntValue  int
    StrValue  string
    BoolValue bool
}

func TestJson(t *testing.T) {
    testJsonStruct := new(testJsonType)
    testJsonStruct.IntValue = 12
    testJsonStruct.StrValue = "34"
    testJsonStruct.BoolValue = true

    if "{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true}" != Json(testJsonStruct) {
        t.Fail()
    }
}

func TestUnJson(t *testing.T) {
    testJsonStr := "{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true}"
    testJsonStruct := UnJson(testJsonStr, new(testJsonType)).(*testJsonType)

    if 12 != testJsonStruct.IntValue ||
        "34" != testJsonStruct.StrValue ||
        !testJsonStruct.BoolValue {
        t.Fail()
    }
}
