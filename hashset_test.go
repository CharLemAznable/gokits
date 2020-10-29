package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewHashset(t *testing.T) {
    a := assert.New(t)
    s := NewHashset()
    s.Add("val1")
    s.Add(1000)
    s.Add('a')

    a.Equal(3, s.Size())
    a.True(s.Contains("val1"))
    a.True(s.Contains(1000))
    a.True(s.Contains('a'))
    a.False(s.Contains("val2"))
    a.False(s.Contains(2000))
    a.False(s.Contains('b'))

    a.True(s.Remove(1000))
    a.Equal(2, s.Size())
}
