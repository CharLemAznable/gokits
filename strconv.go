package gokits

import (
    "strconv"
)

func IntFromStr(str string) (int, error) {
    return strconv.Atoi(str)
}

func Int64FromStr(str string) (int64, error) {
    return strconv.ParseInt(str, 10, 64)
}

func StrFromInt(i int) string {
    return strconv.Itoa(i)
}

func StrFromInt64(i int64) string {
    return strconv.FormatInt(i, 10)
}
