package gokits

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestYComb(test *testing.T) {
    fac := YComb(func(recFunc RecFunc) RecFunc {
        return func(n interface{}) interface{} {
            if n.(int) <= 1 {
                return 1
            } else {
                return n.(int) * recFunc(n.(int) - 1).(int)
            }
        }
    })
    resultFac := fac(10)
    assert.Equal(test, 3628800, resultFac)

    fib := YComb(func(recFunc RecFunc) RecFunc {
        return func(n interface{}) interface{} {
            if n.(int) <= 2 {
                return 1
            } else {
                return recFunc(n.(int) - 1).(int) + recFunc(n.(int) - 2).(int)
            }
        }
    })
    resultFib := fib(10)
    assert.Equal(test, 55, resultFib)
}
