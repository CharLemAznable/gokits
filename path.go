package gokits

import (
    "path"
    "strings"
)

func PathJoin(elem ...string) string {
    result := path.Join(elem...)

    if len(elem) >= 1 {
        last := elem[len(elem)-1]
        if strings.HasSuffix(last, "/") &&
            !strings.HasSuffix(result, "/") {
            result = result + "/"
        }
    }
    return result
}
