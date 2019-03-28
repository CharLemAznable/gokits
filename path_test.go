package gokits

import (
    "testing"
)

func TestJoinPathComponent(t *testing.T) {
    aPath := "//a//"
    bPath := "b//"
    cPath := "//c"
    dPath := "/d/"
    ePath := "e/"
    fPath := "/f"
    gPath := "g"

    join := JoinPathComponent(aPath, bPath, cPath, dPath, ePath, fPath, gPath)
    if "//a/b/c/d/e/f/g" != join {
        t.Fail()
    }
}
