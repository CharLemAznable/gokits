package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

type testJsonType struct {
    IntValue  int
    StrValue  string
    BoolValue bool
}

func TestJson(t *testing.T) {
    a := assert.New(t)
    testJsonStruct := new(testJsonType)
    testJsonStruct.IntValue = 12
    testJsonStruct.StrValue = "34"
    testJsonStruct.BoolValue = true

    a.Equal("{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true}", Json(testJsonStruct))
}

func TestUnJson(t *testing.T) {
    a := assert.New(t)
    testJsonStr := "{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true}"
    testJsonStruct := UnJson(testJsonStr, new(testJsonType)).(*testJsonType)

    a.Equal(12, testJsonStruct.IntValue)
    a.Equal("34", testJsonStruct.StrValue)
    a.True(testJsonStruct.BoolValue)
}

func TestUnJsonArray(t *testing.T) {
    a := assert.New(t)
    testJsonStr := "[{\"IntValue\":12,\"StrValue\":\"34\",\"BoolValue\":true},{\"IntValue\":56,\"StrValue\":\"78\",\"BoolValue\":true}]"
    var v = new([]testJsonType)
    testJsonStructArray := *(UnJson(testJsonStr, v).(*[]testJsonType))

    a.Equal(2, len(testJsonStructArray))

    struct0 := testJsonStructArray[0]
    a.Equal(12, struct0.IntValue)
    a.Equal("34", struct0.StrValue)
    a.True(struct0.BoolValue)

    struct1 := testJsonStructArray[1]
    a.Equal(56, struct1.IntValue)
    a.Equal("78", struct1.StrValue)
    a.True(struct1.BoolValue)
}
