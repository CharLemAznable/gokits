package gokits

import (
    "testing"
)

func TestPathJoin(t *testing.T) {
    if "a/b" != PathJoin("a", "b") {
        t.Fail()
    }

    if "/a/b" != PathJoin("/a", "/b") {
        t.Fail()
    }

    if "/a/b/" != PathJoin("/a", "/b/") {
        t.Fail()
    }

    if "/" != PathJoin("/a", "..") {
        t.Fail()
    }

    if "." != PathJoin(".") {
        t.Fail()
    }
}
