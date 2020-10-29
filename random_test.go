package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRandomString(t *testing.T) {
    random1 := RandomString(20)
    random2 := RandomString(20)
    assert.NotEqual(t, random1, random2)
}
