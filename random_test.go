package gokits

import (
    "testing"
)

func TestRandomString(t *testing.T) {
    random1 := RandomString(20)
    random2 := RandomString(20)
    if random1 == random2 {
        t.Fail()
    }
}
