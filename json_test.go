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

func TestUnJsonArray(t *testing.T) {
    testJsonStr := "[{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true},{\"IntValue\":56,\"StrValue\":\"78\",\"BoolValue\":true}]"
    var v = new([]testJsonType)
    testJsonStructArray := *(UnJson(testJsonStr, v).(*[]testJsonType))

    if 2 != len(testJsonStructArray) {
        t.Fail()
    }

    struct0 := testJsonStructArray[0]
    if 12 != struct0.IntValue ||
        "34" != struct0.StrValue ||
        !struct0.BoolValue {
        t.Fail()
    }

    struct1 := testJsonStructArray[1]
    if 56 != struct1.IntValue ||
        "78" != struct1.StrValue ||
        !struct1.BoolValue {
        t.Fail()
    }
}
