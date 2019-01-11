package gokits

import (
    "testing"
)

func TestStrFromInt(t *testing.T) {
    var i = 32
    if "32" != StrFromInt(i) {
        t.Fail()
    }

    var i64 int64 = 2147483648
    if "2147483648" != StrFromInt64(i64) {
        t.Fail()
    }
}

func TestIntFromStr(t *testing.T) {
    i, _ := IntFromStr("32")
    if 32 != i {
        t.Fail()
    }

    i64, _ := Int64FromStr("2147483648")
    if 2147483648 != i64 {
        t.Fail()
    }

    ifail, _ := IntFromStr("AB")
    if 0 != ifail {
        t.Fail()
    }

    i64fail, _ := Int64FromStr("<{[()]}>")
    if 0 != i64fail {
        t.Fail()
    }
}
