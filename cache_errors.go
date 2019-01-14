/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package gokits

import (
    "errors"
)

var (
    // ErrCacheKeyNotFound gets returned when a specific key couldn't be found
    ErrCacheKeyNotFound = errors.New("key not found in cache")
    // ErrCacheKeyNotFoundOrLoadable gets returned when a specific key couldn't be
    // found and loading via the data-loader callback also failed
    ErrCacheKeyNotFoundOrLoadable = errors.New("key not found and could not be loaded into cache")
)
