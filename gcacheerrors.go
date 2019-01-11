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
    // ErrKeyNotFound gets returned when a specific key couldn't be found
    ErrKeyNotFound = errors.New("key not found in cache")
    // ErrKeyNotFoundOrLoadable gets returned when a specific key couldn't be
    // found and loading via the data-loader callback also failed
    ErrKeyNotFoundOrLoadable = errors.New("key not found and could not be loaded into cache")
)