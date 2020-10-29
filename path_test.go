package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPathJoin(t *testing.T) {
    a := assert.New(t)

    a.Equal("a/b", PathJoin("a", "b"))
    a.Equal("/a/b", PathJoin("/a", "/b"))
    a.Equal("/a/b/", PathJoin("/a", "/b/"))
    a.Equal("/", PathJoin("/a", ".."))
    a.Equal(".", PathJoin("."))
}
