package gokits

type RecFunc func(interface{}) interface{}

func YComb(f func(func(interface{}) interface{}) func(interface{}) interface{}) func(interface{}) interface{} {
    return func(n interface{}) interface{} {
        return f(YComb(f))(n)
    }
}
