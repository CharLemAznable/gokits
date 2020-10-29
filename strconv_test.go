package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestStrFromInt(t *testing.T) {
    var i = 32
    assert.Equal(t, "32", StrFromInt(i))

    var i64 int64 = 2147483648
    assert.Equal(t, "2147483648", StrFromInt64(i64))
}

func TestIntFromStr(t *testing.T) {
    i, _ := IntFromStr("32")
    assert.Equal(t, 32, i)

    i64, _ := Int64FromStr("2147483648")
    assert.EqualValues(t, 2147483648, i64)

    ifail, _ := IntFromStr("AB")
    assert.Equal(t, 0, ifail)

    i64fail, _ := Int64FromStr("<{[()]}>")
    assert.EqualValues(t, 0, i64fail)
}
