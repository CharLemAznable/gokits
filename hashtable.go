package gokits

import (
    "sync"
)

type Hashtable struct {
    mutex  sync.RWMutex
    mapper map[interface{}]interface{}
}

func NewHashtable() *Hashtable {
    return &Hashtable{mapper: make(map[interface{}]interface{})}
}

func (hashtable *Hashtable) Put(key, value interface{}) interface{} {
    hashtable.mutex.Lock()
    defer hashtable.mutex.Unlock()

    old := hashtable.mapper[key]
    hashtable.mapper[key] = value

    return old
}

func (hashtable *Hashtable) Get(key interface{}) interface{} {
    hashtable.mutex.RLock()
    defer hashtable.mutex.RUnlock()

    return hashtable.mapper[key]
}

func (hashtable *Hashtable) Remove(key interface{}) interface{} {
    hashtable.mutex.Lock()
    defer hashtable.mutex.Unlock()

    old := hashtable.mapper[key]
    delete(hashtable.mapper, key)

    return old
}

func (hashtable *Hashtable) Size() int {
    hashtable.mutex.RLock()
    defer hashtable.mutex.RUnlock()

    return len(hashtable.mapper)
}

func (hashtable *Hashtable) Keys() []interface{} {
    hashtable.mutex.Lock()
    defer hashtable.mutex.Unlock()

    keys := make([]interface{}, 0, len(hashtable.mapper))
    for key := range hashtable.mapper {
        keys = append(keys, key)
    }

    return keys
}
