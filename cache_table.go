/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package gokits

import (
    "log"
    "sort"
    "sync"
    "time"
)

type ExpireStrategy int

const (
    ExpireAfterWrite  = 0
    ExpireAfterAccess = 1
)

// CacheTable is a table within the cache
type CacheTable struct {
    sync.RWMutex

    // The table's name.
    name string
    // All cached items.
    items map[interface{}]*CacheItem
    // CacheItem expire strategy
    strategy ExpireStrategy

    // Timer responsible for triggering cleanup.
    cleanupTimer *time.Timer
    // Current timer duration.
    cleanupInterval time.Duration

    // The logger used for this table.
    logger *log.Logger

    // Callback method triggered when trying to load a non-existing key.
    loadData func(key interface{}, args ...interface{}) (*CacheItem, error)
    // Callback method triggered when adding a new item to the cache.
    addedItem func(item *CacheItem)
    // Callback method triggered before deleting an item from the cache.
    aboutToDeleteItem func(item *CacheItem)
}

// Count returns how many items are currently stored in the cache.
func (table *CacheTable) Count() int {
    table.RLock()
    defer table.RUnlock()
    return len(table.items)
}

// Foreach all items
func (table *CacheTable) Foreach(trans func(key interface{}, item *CacheItem)) {
    table.RLock()
    defer table.RUnlock()

    for k, v := range table.items {
        trans(k, v)
    }
}

// SetDataLoader configures a data-loader callback, which will be called when
// trying to access a non-existing key. The key and 0...n additional arguments
// are passed to the callback function.
func (table *CacheTable) SetDataLoader(f func(interface{}, ...interface{}) (*CacheItem, error)) {
    table.Lock()
    defer table.Unlock()
    table.loadData = f
}

// SetAddedItemCallback configures a callback, which will be called every time
// a new item is added to the cache.
func (table *CacheTable) SetAddedItemCallback(f func(*CacheItem)) {
    table.Lock()
    defer table.Unlock()
    table.addedItem = f
}

// SetAboutToDeleteItemCallback configures a callback, which will be called
// every time an item is about to be removed from the cache.
func (table *CacheTable) SetAboutToDeleteItemCallback(f func(*CacheItem)) {
    table.Lock()
    defer table.Unlock()
    table.aboutToDeleteItem = f
}

// SetLogger sets the logger to be used by this cache table.
func (table *CacheTable) SetLogger(logger *log.Logger) {
    table.Lock()
    defer table.Unlock()
    table.logger = logger
}

// Expiration check loop, triggered by a self-adjusting timer.
func (table *CacheTable) expirationCheck() {
    table.Lock()
    if table.cleanupTimer != nil {
        table.cleanupTimer.Stop()
    }
    if table.cleanupInterval > 0 {
        table.log("Expiration check triggered after", table.cleanupInterval, "for table", table.name)
    } else {
        table.log("Expiration check installed for table", table.name)
    }
    table.expirationCheckInternal()
    table.Unlock()
}

func (table *CacheTable) expirationCheckInternal() {
    // To be more accurate with timers, we would need to update 'now' on every
    // loop iteration. Not sure it's really efficient though.
    isExpireAfterAccess := ExpireAfterAccess == table.strategy
    now := time.Now()
    smallestDuration := 0 * time.Second
    for key, item := range table.items {
        // Cache values so we don't keep blocking the mutex.
        item.RLock()
        lifeSpan := item.lifeSpan
        baseOn := item.createdOn
        if isExpireAfterAccess {
            baseOn = item.accessedOn
        }
        item.RUnlock()

        if lifeSpan == 0 {
            continue
        }
        if now.Sub(baseOn) >= lifeSpan {
            // Item has exceeded its lifespan.
            _, _ = table.deleteInternal(key)
        } else {
            // Find the item chronologically closest to its end-of-lifespan.
            if smallestDuration == 0 || lifeSpan-now.Sub(baseOn) < smallestDuration {
                smallestDuration = lifeSpan - now.Sub(baseOn)
            }
        }
    }

    // Setup the interval for the next cleanup run.
    table.cleanupInterval = smallestDuration
    if smallestDuration > 0 {
        table.cleanupTimer = time.AfterFunc(smallestDuration, func() {
            go table.expirationCheck()
        })
    }
}

func (table *CacheTable) addInternal(item *CacheItem) {
    // Careful: do not run this method unless the table-mutex is locked!
    // It will unlock it for the caller before running the callbacks and checks
    table.log("Adding item with key", item.key, "and lifespan of", item.lifeSpan, "to table", table.name)
    table.items[item.key] = item

    // Cache values so we don't keep blocking the mutex.
    expDur := table.cleanupInterval
    addedItem := table.addedItem
    table.Unlock()

    // Trigger callback after adding an item to cache.
    if addedItem != nil {
        addedItem(item)
    }

    // If we haven't set up any expiration check timer or found a more imminent item.
    if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
        table.expirationCheck()
    }
}

// Add adds a key/value pair to the cache.
// Parameter key is the item's cache-key.
// Parameter lifeSpan determines after which time period without an access the item
// will get removed from the cache.
// Parameter data is the item's value.
func (table *CacheTable) Add(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
    item := NewCacheItem(key, lifeSpan, data)

    // Add item to cache.
    table.Lock()
    table.addInternal(item)

    return item
}

func (table *CacheTable) deleteInternal(key interface{}) (*CacheItem, error) {
    r, ok := table.items[key]
    if !ok {
        return nil, ErrCacheKeyNotFound
    }

    // Cache value so we don't keep blocking the mutex.
    aboutToDeleteItem := table.aboutToDeleteItem
    table.Unlock()

    // Trigger callbacks before deleting an item from cache.
    if aboutToDeleteItem != nil {
        aboutToDeleteItem(r)
    }

    r.RLock()
    defer r.RUnlock()
    if r.aboutToExpire != nil {
        r.aboutToExpire(key)
    }

    table.Lock()
    table.log("Deleting item with key", key, "created on", r.createdOn, "and hit", r.accessCount, "times from table", table.name)
    delete(table.items, key)

    return r, nil
}

// Delete an item from the cache.
func (table *CacheTable) Delete(key interface{}) (*CacheItem, error) {
    table.Lock()
    defer table.Unlock()

    return table.deleteInternal(key)
}

// Exists returns whether an item exists in the cache. Unlike the Value method
// Exists neither tries to fetch data via the loadData callback nor does it
// keep the item alive in the cache.
func (table *CacheTable) Exists(key interface{}) bool {
    table.RLock()
    defer table.RUnlock()
    _, ok := table.items[key]

    return ok
}

// NotFoundAdd tests whether an item not found in the cache. Unlike the Exists
// method this also adds data if they key could not be found.
func (table *CacheTable) NotFoundAdd(key interface{}, lifeSpan time.Duration, data interface{}) bool {
    table.Lock()

    if _, ok := table.items[key]; ok {
        table.Unlock()
        return false
    }

    item := NewCacheItem(key, lifeSpan, data)
    table.addInternal(item)

    return true
}

func (table *CacheTable) LoadingAdd(key interface{}, args ...interface{}) (*CacheItem, error) {
    loadData := table.loadData

    if loadData != nil {
        table.Lock()

        r, ok := table.items[key]
        if ok {
            // loaded by other goroutine
            table.Unlock()
            return r, nil
        }

        item, err := loadData(key, args...)
        if err == nil && item != nil {
            table.addInternal(item)
            return item, nil
        }

        table.Unlock()
        if err == nil {
            err = ErrCacheKeyNotFoundOrLoadable
        }
        return nil, err
    }
    return nil, ErrCacheKeyNotFound
}

// Value returns an item from the cache and marks it to be kept alive. You can
// pass additional arguments to your DataLoader callback function.
func (table *CacheTable) Value(key interface{}, args ...interface{}) (*CacheItem, error) {
    table.RLock()
    r, ok := table.items[key]
    table.RUnlock()

    if ok {
        // Update access counter and timestamp.
        r.KeepAlive()
        return r, nil
    }

    // Item doesn't exist in cache. Try and fetch it with a data-loader.
    return table.LoadingAdd(key, args...)
}

// Flush deletes all items from this cache table.
func (table *CacheTable) Flush() {
    table.Lock()
    defer table.Unlock()

    table.log("Flushing table", table.name)

    table.items = make(map[interface{}]*CacheItem)
    table.cleanupInterval = 0
    if table.cleanupTimer != nil {
        table.cleanupTimer.Stop()
    }
}

// CacheItemPair maps key to access counter
type CacheItemPair struct {
    Key         interface{}
    AccessCount int64
}

// CacheItemPairList is a slice of CacheIemPairs that implements sort.
// Interface to sort by AccessCount.
type CacheItemPairList []CacheItemPair

func (p CacheItemPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p CacheItemPairList) Len() int           { return len(p) }
func (p CacheItemPairList) Less(i, j int) bool { return p[i].AccessCount > p[j].AccessCount }

// MostAccessed returns the most accessed items in this cache table
func (table *CacheTable) MostAccessed(count int64) []*CacheItem {
    table.RLock()
    defer table.RUnlock()

    p := make(CacheItemPairList, len(table.items))
    i := 0
    for k, v := range table.items {
        p[i] = CacheItemPair{k, v.accessCount}
        i++
    }
    sort.Sort(p)

    var r []*CacheItem
    c := int64(0)
    for _, v := range p {
        if c >= count {
            break
        }

        item, ok := table.items[v.Key]
        if ok {
            r = append(r, item)
        }
        c++
    }

    return r
}

// Internal logging method for convenience.
func (table *CacheTable) log(v ...interface{}) {
    if table.logger == nil {
        return
    }

    table.logger.Println(v...)
}
