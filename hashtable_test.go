package gokits

import (
    "sort"
    "strings"
    "testing"
)

func TestNewHashtable(t *testing.T) {
    h := NewHashtable()
    h.Put("key0", "val1")
    h.Put("key1", 1000)
    h.Put("key2", 'a')

    if 3 != h.Size() {
        t.Fail()
    }
    keys := h.Keys()
    sort.Slice(keys, func(i, j int) bool {
        return strings.Compare(keys[i].(string), keys[j].(string)) < 0
    })
    if "key0" != keys[0] || "key1" != keys[1] || "key2" != keys[2] {
        t.Fail()
    }
    if "val1" != h.Get("key0") {
        t.Fail()
    }
    if 1000 != h.Get("key1") {
        t.Fail()
    }
    if 'a' != h.Get("key2") {
        t.Fail()
    }

    if 1000 != h.Remove("key1") {
        t.Fail()
    }
    if 2 != h.Size() {
        t.Fail()
    }
    keys = h.Keys()
    sort.Slice(keys, func(i, j int) bool {
        return strings.Compare(keys[i].(string), keys[j].(string)) < 0
    })
    if "key0" != keys[0] || "key2" != keys[1] {
        t.Fail()
    }
}