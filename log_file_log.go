// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package gokits

import (
    "fmt"
    "os"
    "time"
)

// This log writer sends output to a file
//noinspection GoSnakeCaseUsage
type FileLogWriter struct {
    rec chan *LogRecord
    rot chan bool

    // The opened file
    filename string
    file     *os.File

    // The logging format
    format string

    // File header/trailer
    header, trailer string

    // Rotate at linecount
    maxlines          int
    maxlines_curlines int

    // Rotate at size
    maxsize         int
    maxsize_cursize int

    // Rotate daily
    daily          bool
    daily_opendate int

    // Keep old logfiles (.001, .002, etc)
    rotate bool
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
    w.rec <- rec
}

func (w *FileLogWriter) Close() {
    close(w.rec)
}

const FileLogWriterFormat = "FileLogWriter(%q): %s\n"
// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fname string, rotate bool) *FileLogWriter {
    w := &FileLogWriter{
        rec:      make(chan *LogRecord, LogBufferLength),
        rot:      make(chan bool),
        filename: fname,
        format:   "[%D %T] [%L] (%S) %M",
        rotate:   rotate,
    }

    // open the file for the first time
    if err := w.intRotate(); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, FileLogWriterFormat, w.filename, err)
        return nil
    }

    go w.loopInternal()

    return w
}

func (w *FileLogWriter) loopInternal() {
    defer w.deferInternal()

    for true {
        select {
        case <-w.rot:
            if err := w.intRotate(); err != nil {
                _, _ = fmt.Fprintf(os.Stderr, FileLogWriterFormat, w.filename, err)
                return
            }
        case rec, ok := <-w.rec:
            if ret := w.recInternal(ok, rec); ret {
                return
            }
        }
    }
}

func (w *FileLogWriter) deferInternal() {
    if w.file != nil {
        _, _ = fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
        _ = w.file.Close()
    }
}

func (w *FileLogWriter) recInternal(ok bool, rec *LogRecord) bool {
    if !ok {
        // Flush
        _ = w.file.Sync()
        return true
    }
    now := time.Now()
    if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
        (w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
        (w.daily && now.Day() != w.daily_opendate) {
        if err := w.intRotate(); err != nil {
            _, _ = fmt.Fprintf(os.Stderr, FileLogWriterFormat, w.filename, err)
            return true
        }
    }

    // Perform the write
    n, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, FileLogWriterFormat, w.filename, err)
        return true
    }

    // Flush
    _ = w.file.Sync()
    // Update the counts
    w.maxlines_curlines++
    w.maxsize_cursize += n
    return false
}

// Request that the logs rotate
func (w *FileLogWriter) Rotate() {
    w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotate() error {
    // Close any log file that may be open
    if w.file != nil {
        _, _ = fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
        _ = w.file.Close()
    }

    // If we are keeping log files, move it to the next available number
    if w.rotate {
        _, err := os.Lstat(w.filename)
        if err == nil { // file exists
            // Find the next available number
            num := 1
            fname := ""
            for ; err == nil && num <= 999; num++ {
                fname = w.filename + fmt.Sprintf(".%03d", num)
                _, err = os.Lstat(fname)
            }
            // return error if the last file checked still existed
            if err == nil {
                return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
            }

            // Rename the file to its newfound home
            err = os.Rename(w.filename, fname)
            if err != nil {
                return fmt.Errorf("Rotate: %s\n", err)
            }
        }
    }

    // Open the log file
    fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
    if err != nil {
        return err
    }
    w.file = fd

    now := time.Now()
    _, _ = fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

    // Set the daily open date to the current date
    w.daily_opendate = now.Day()

    // initialize rotation values
    w.maxlines_curlines = 0
    w.maxsize_cursize = 0

    return nil
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
    w.format = format
    return w
}

// Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
    w.header, w.trailer = head, foot
    if w.maxlines_curlines == 0 {
        _, _ = fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
    }
    return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
    // fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
    w.maxlines = maxlines
    return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
    // fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
    w.maxsize = maxsize
    return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
    // fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
    w.daily = daily
    return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
    // fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
    w.rotate = rotate
    return w
}

// NewXMLLogWriter is a utility method for creating a FileLogWriter set up to
// output XML record log messages instead of line-based ones.
func NewXMLLogWriter(fname string, rotate bool) *FileLogWriter {
    return NewFileLogWriter(fname, rotate).SetFormat(
        `	<record level="%L">
		<timestamp>%D %T</timestamp>
		<source>%S</source>
		<message>%M</message>
	</record>`).SetHeadFoot("<log created=\"%D %T\">", "</log>")
}
