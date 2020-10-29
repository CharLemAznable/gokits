package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestArrayIndexOf(t *testing.T) {
    a := assert.New(t)
    arr := []string{"aa", "bb", "cc"}
    a.False(ArrayIndexOf("bb", arr) < 0)
    a.False(ArrayIndexOf("dd", arr) >= 0)
    a.False(ArrayIndexOf("BB", arr) >= 0)
    a.False(ArrayIndexOf("DD", arr) >= 0)

    a.False(ArrayIndexOfIgnoreCase("bb", arr) < 0)
    a.False(ArrayIndexOfIgnoreCase("dd", arr) >= 0)
    a.False(ArrayIndexOfIgnoreCase("BB", arr) < 0)
    a.False(ArrayIndexOfIgnoreCase("DD", arr) >= 0)
}

func TestArrayContains(t *testing.T) {
    a := assert.New(t)
    arr := []string{"aa", "bb", "cc"}
    a.True(ArrayContains("bb", arr))
    a.False(ArrayContains("dd", arr))
    a.False(ArrayContains("BB", arr))
    a.False(ArrayContains("DD", arr))

    a.True(ArrayContainsIgnoreCase("bb", arr))
    a.False(ArrayContainsIgnoreCase("dd", arr))
    a.True(ArrayContainsIgnoreCase("BB", arr))
    a.False(ArrayContainsIgnoreCase("DD", arr))
}
