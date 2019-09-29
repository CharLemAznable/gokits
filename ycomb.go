package gokits

type RecFunc func(interface{}) interface{}

func YComb(f func(RecFunc) RecFunc) RecFunc {
    return func(n interface{}) interface{} {
        return f(YComb(f))(n)
    }
}
