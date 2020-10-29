package gokits

import (
    "github.com/stretchr/testify/assert"
    "sort"
    "strings"
    "testing"
)

func TestNewHashtable(t *testing.T) {
    a := assert.New(t)
    h := NewHashtable()
    h.Put("key0", "val1")
    h.Put("key1", 1000)
    h.Put("key2", 'a')

    a.Equal(3, h.Size())
    keys := h.Keys()
    sort.Slice(keys, func(i, j int) bool {
        return strings.Compare(keys[i].(string), keys[j].(string)) < 0
    })
    a.False("key0" != keys[0] || "key1" != keys[1] || "key2" != keys[2])
    a.Equal("val1", h.Get("key0"))
    a.Equal(1000, h.Get("key1"))
    a.Equal('a', h.Get("key2"))

    a.Equal(1000, h.Remove("key1"))
    a.Equal(2, h.Size())
    keys = h.Keys()
    sort.Slice(keys, func(i, j int) bool {
        return strings.Compare(keys[i].(string), keys[j].(string)) < 0
    })
    a.False("key0" != keys[0] || "key2" != keys[1])
}
