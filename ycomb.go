package gokits

type RecFunc func(n interface{}) interface{}

func YComb(b func(RecFunc) RecFunc) RecFunc {
    var g = func(f func(t RecFunc) RecFunc) RecFunc {
        var r = func(y interface{}) RecFunc {
            var w = y.(func(v interface{}) RecFunc)
            return f(func(n interface{}) interface{} {
                return w(w)(n)
            })
        }
        return r(r)
    }
    return g(b)
}