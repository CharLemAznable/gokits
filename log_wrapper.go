// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package gokits

import (
    "errors"
    "fmt"
    "os"
    "strings"
)

type logWrapper struct {
    logger Logger
}

var (
    LOG logWrapper
)

func init() {
    LOG = logWrapper{}
    LOG.logger = NewDefaultLogger(DEBUG)
}

// Wrapper for (*Logger).LoadConfiguration
func (wrapper logWrapper) LoadConfiguration(filename string) {
    wrapper.logger.LoadConfiguration(filename)
}

// Wrapper for (*Logger).AddFilter
func (wrapper logWrapper) AddFilter(name string, lvl level, writer LogWriter) {
    wrapper.logger.AddFilter(name, lvl, writer)
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func (wrapper logWrapper) Close() {
    wrapper.logger.Close()
}

func (wrapper logWrapper) Crash(args ...interface{}) {
    if len(args) > 0 {
        wrapper.logger.intLogf(CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
    }
    panic(args)
}

// Logs the given message and crashes the program
func (wrapper logWrapper) Crashf(format string, args ...interface{}) {
    wrapper.logger.intLogf(CRITICAL, format, args...)
    wrapper.logger.Close() // so that hopefully the messages get logged
    panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func (wrapper logWrapper) Exit(args ...interface{}) {
    if len(args) > 0 {
        wrapper.logger.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
    }
    wrapper.logger.Close() // so that hopefully the messages get logged
    os.Exit(0)
}

// Compatibility with `log`
func (wrapper logWrapper) Exitf(format string, args ...interface{}) {
    wrapper.logger.intLogf(ERROR, format, args...)
    wrapper.logger.Close() // so that hopefully the messages get logged
    os.Exit(0)
}

// Compatibility with `log`
func (wrapper logWrapper) Stderr(args ...interface{}) {
    if len(args) > 0 {
        wrapper.logger.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
    }
}

// Compatibility with `log`
func (wrapper logWrapper) Stderrf(format string, args ...interface{}) {
    wrapper.logger.intLogf(ERROR, format, args...)
}

// Compatibility with `log`
func (wrapper logWrapper) Stdout(args ...interface{}) {
    if len(args) > 0 {
        wrapper.logger.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
    }
}

// Compatibility with `log`
func (wrapper logWrapper) Stdoutf(format string, args ...interface{}) {
    wrapper.logger.intLogf(INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func (wrapper logWrapper) Log(lvl level, source, message string) {
    wrapper.logger.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func (wrapper logWrapper) Logf(lvl level, format string, args ...interface{}) {
    wrapper.logger.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func (wrapper logWrapper) Logc(lvl level, closure func() string) {
    wrapper.logger.intLogc(lvl, closure)
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func (wrapper logWrapper) Finest(arg0 interface{}, args ...interface{}) {
    const (
        lvl = FINEST
    )
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
    case func() string:
        // Log the closure (no other arguments used)
        wrapper.logger.intLogc(lvl, first)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func (wrapper logWrapper) Fine(arg0 interface{}, args ...interface{}) {
    const (
        lvl = FINE
    )
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
    case func() string:
        // Log the closure (no other arguments used)
        wrapper.logger.intLogc(lvl, first)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func (wrapper logWrapper) Debug(arg0 interface{}, args ...interface{}) {
    const (
        lvl = DEBUG
    )
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
    case func() string:
        // Log the closure (no other arguments used)
        wrapper.logger.intLogc(lvl, first)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func (wrapper logWrapper) Trace(arg0 interface{}, args ...interface{}) {
    const (
        lvl = TRACE
    )
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
    case func() string:
        // Log the closure (no other arguments used)
        wrapper.logger.intLogc(lvl, first)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func (wrapper logWrapper) Info(arg0 interface{}, args ...interface{}) {
    const (
        lvl = INFO
    )
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
    case func() string:
        // Log the closure (no other arguments used)
        wrapper.logger.intLogc(lvl, first)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func (wrapper logWrapper) Warn(arg0 interface{}, args ...interface{}) error {
    const (
        lvl = WARNING
    )
    return wrapper.errorInternal(lvl, arg0, args...)
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func (wrapper logWrapper) Error(arg0 interface{}, args ...interface{}) error {
    const (
        lvl = ERROR
    )
    return wrapper.errorInternal(lvl, arg0, args...)
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func (wrapper logWrapper) Critical(arg0 interface{}, args ...interface{}) error {
    const (
        lvl = CRITICAL
    )
    return wrapper.errorInternal(lvl, arg0, args...)
}

func (wrapper logWrapper) errorInternal(lvl level, arg0 interface{}, args ...interface{}) error {
    switch first := arg0.(type) {
    case string:
        // Use the string as a format string
        wrapper.logger.intLogf(lvl, first, args...)
        return errors.New(fmt.Sprintf(first, args...))
    case func() string:
        // Log the closure (no other arguments used)
        str := first()
        wrapper.logger.intLogf(lvl, "%s", str)
        return errors.New(str)
    default:
        // Build a format string so that it will be similar to Sprint
        wrapper.logger.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
        return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
    }
}
