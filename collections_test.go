package gokits

import (
    "testing"
)

func TestArrayIndexOf(t *testing.T) {
    arr := []string{"aa", "bb", "cc"}
    if ArrayIndexOf("bb", arr) < 0 {
        t.Fail()
    }
    if ArrayIndexOf("dd", arr) >= 0 {
        t.Fail()
    }
    if ArrayIndexOf("BB", arr) >= 0 {
        t.Fail()
    }
    if ArrayIndexOf("DD", arr) >= 0 {
        t.Fail()
    }

    if ArrayIndexOfIgnoreCase("bb", arr) < 0 {
        t.Fail()
    }
    if ArrayIndexOfIgnoreCase("dd", arr) >= 0 {
        t.Fail()
    }
    if ArrayIndexOfIgnoreCase("BB", arr) < 0 {
        t.Fail()
    }
    if ArrayIndexOfIgnoreCase("DD", arr) >= 0 {
        t.Fail()
    }
}

func TestArrayContains(t *testing.T) {
    arr := []string{"aa", "bb", "cc"}
    if !ArrayContains("bb", arr) {
        t.Fail()
    }
    if ArrayContains("dd", arr) {
        t.Fail()
    }
    if ArrayContains("BB", arr) {
        t.Fail()
    }
    if ArrayContains("DD", arr) {
        t.Fail()
    }

    if !ArrayContainsIgnoreCase("bb", arr) {
        t.Fail()
    }
    if ArrayContainsIgnoreCase("dd", arr) {
        t.Fail()
    }
    if !ArrayContainsIgnoreCase("BB", arr) {
        t.Fail()
    }
    if ArrayContainsIgnoreCase("DD", arr) {
        t.Fail()
    }
}
