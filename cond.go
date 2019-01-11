package gokits

func Condition(cond bool, trueVal, falseVal interface{}) interface{} {
    if cond {
        return ValOrFunc(trueVal)
    }
    return ValOrFunc(falseVal)
}

func If(cond bool, trueFunc func()) {
    if cond {
        trueFunc()
    }
}

func Unless(cond bool, falseFunc func()) {
    if !cond {
        falseFunc()
    }
}

func DefaultIfNil(val, def interface{}) interface{} {
    if nil != val {
        return val
    }
    return ValOrFunc(def)
}

func ValOrFunc(val interface{}) interface{} {
    switch val.(type) {
    case func() interface{}:
        return (val.(func() interface{}))()
    default:
        return val
    }
}
