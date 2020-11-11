package gokits

import (
    "time"
)

type Duration struct {
    time.Duration
}

func (d *Duration) MarshalText() ([]byte, error) {
    return []byte(d.Duration.String()), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
    x, err := time.ParseDuration(string(text))
    if err != nil {
        return err
    }
    d.Duration = x
    return nil
}
