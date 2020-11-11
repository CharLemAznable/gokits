package gokits

import (
    "github.com/BurntSushi/toml"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

type TestTime struct {
    Duration Duration
}

func TestDuration(t *testing.T) {
    testTime := &TestTime{}
    _, _ = toml.DecodeFile("time_test.toml", testTime)
    assert.Equal(t, time.Minute*4+time.Second*2, testTime.Duration.Duration)
}
