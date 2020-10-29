package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestCondition(t *testing.T) {
    a := assert.New(t)
    trueVal := Condition(true, 1, 2)
    a.Equal(1, trueVal)

    falseVal := Condition(false, func() interface{} {
        return "true"
    }, func() interface{} {
        return "false"
    })
    a.Equal("false", falseVal)
}

func TestIfAndUnless(t *testing.T) {
    a := assert.New(t)
    trueVal := ""
    If(true, func() {
        trueVal = "true"
    })
    a.Equal("true", trueVal)

    If(false, func() {
        a.Fail("error condition")
    })

    falseVal := "false"
    Unless(false, func() {
        falseVal = "false"
    })
    a.Equal("false", falseVal)
    Unless(true, func() {
        a.Fail("error condition")
    })
}

func TestDefaultIfNil(t *testing.T) {
    a := assert.New(t)
    var nilVal interface{} = nil

    trueVal := DefaultIfNil(nilVal, "true")
    a.Equal("true", trueVal)
    trueVal = DefaultIfNil("T", "true")
    a.Equal("T", trueVal)

    falseVal := DefaultIfNil(nil, func() interface{} {
        return "false"
    })
    a.Equal("false", falseVal)
    falseVal = DefaultIfNil("F", "false")
    a.Equal("F", falseVal)
}
