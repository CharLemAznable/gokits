package gokits

import (
    "testing"
)

func TestCondition(t *testing.T) {
    trueVal := Condition(true, 1, 2)
    if 1 != trueVal {
        t.Fail()
    }

    falseVal := Condition(false, func() interface{} {
        return "true"
    }, func() interface{} {
        return "false"
    })
    if "false" != falseVal {
        t.Fail()
    }
}

func TestIfAndUnless(t *testing.T) {
    trueVal := ""
    If(true, func() {
        trueVal = "true"
    })
    if "true" != trueVal {
        t.Fail()
    }

    If(false, func() {
        t.Fail()
    })

    falseVal := "false"
    Unless(false, func() {
        falseVal = "false"
    })
    if "false" != falseVal {
        t.Fail()
    }
    Unless(true, func() {
        t.Fail()
    })
}

func TestDefaultIfNil(t *testing.T) {
    var nilVal interface{} = nil

    trueVal := DefaultIfNil(nilVal, "true")
    if "true" != trueVal {
        t.Fail()
    }
    trueVal = DefaultIfNil("T", "true")
    if "T" != trueVal {
        t.Fail()
    }

    falseVal := DefaultIfNil(nil, func() interface{} {
        return "false"
    })
    if "false" != falseVal {
        t.Fail()
    }
    falseVal = DefaultIfNil("F", "false")
    if "F" != falseVal {
        t.Fail()
    }
}
