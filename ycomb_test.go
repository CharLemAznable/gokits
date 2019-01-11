package gokits

import (
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
    if 3628800 != resultFac {
        test.Error("YComb(fac) Error: fac(10) should be 3628800 but ", resultFac)
    }

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
    if 55 != resultFib {
        test.Error("YComb(fib) Error: fib(10) should be 55 but ", resultFib)
    }
}
