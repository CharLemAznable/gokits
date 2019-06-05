package gokits

import (
    "testing"
)

func TestNewHashset(t *testing.T) {
    s := NewHashset()
    s.Add("val1")
    s.Add(1000)
    s.Add('a')

    if 3 != s.Size() {
        t.Fail()
    }
    if !s.Contains("val1") {
        t.Fail()
    }
    if !s.Contains(1000) {
        t.Fail()
    }
    if !s.Contains('a') {
        t.Fail()
    }
    if s.Contains("val2") {
        t.Fail()
    }
    if s.Contains(2000) {
        t.Fail()
    }
    if s.Contains('b') {
        t.Fail()
    }

    if !s.Remove(1000) {
        t.Fail()
    }
    if 2 != s.Size() {
        t.Fail()
    }
}
