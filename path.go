package gokits

import (
    "strings"
)

func JoinPathComponent(aPath, bPath string, morePaths ...string) string {
    ap := aPath
    for strings.HasSuffix(ap, "/") {
        ap = ap[:len(ap)-1]
    }
    bp := bPath
    for strings.HasPrefix(bp, "/") {
        bp = bp[1:]
    }
    result := ap + "/" + bp

    return Condition(len(morePaths) > 0, func() interface{} {
        return JoinPathComponent(result, morePaths[0], morePaths[1:]...)
    }, result).(string)
}
