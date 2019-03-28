package gokits

import (
    "strings"
)

func ArrayIndexOf(aStr string, arr []string) int {
    for index, value := range arr {
        if aStr == value {
            return index
        }
    }
    return -1
}

func ArrayIndexOfIgnoreCase(aStr string, arr []string) int {
    for index, value := range arr {
        if strings.ToLower(aStr) == strings.ToLower(value) {
            return index
        }
    }
    return -1
}

func ArrayContains(aStr string, arr []string) bool {
    return ArrayIndexOf(aStr, arr) >= 0
}

func ArrayContainsIgnoreCase(aStr string, arr []string) bool {
    return ArrayIndexOfIgnoreCase(aStr, arr) >= 0
}
