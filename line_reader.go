package gokits

import (
    "fmt"
    "io"
)

type LineReader struct {
    inByteBuf []byte
    lineBuf   []byte
    inLimit   int
    inOff     int
    reader    io.Reader
}

type readingLineTemp struct {
    length             int
    c                  byte
    skipWhiteSpace     bool
    isCommentLine      bool
    isNewLine          bool
    appendedLineBegin  bool
    precedingBackslash bool
    skipLF             bool
    err                error
}

func NewLineReader(reader io.Reader) *LineReader {
    return &LineReader{
        inByteBuf: make([]byte, 8192),
        lineBuf:   make([]byte, 1024),
        reader:    reader,
        inLimit:   0,
        inOff:     0,
    }
}

func (lineReader *LineReader) ReadLine() (int, error) {
    temp := &readingLineTemp{
        length:             0,
        c:                  0,
        skipWhiteSpace:     true,
        isCommentLine:      false,
        isNewLine:          true,
        appendedLineBegin:  false,
        precedingBackslash: false,
        skipLF:             false,
        err:                nil,
    }

    for true {
        ret := lineReader.readInByteBuf(temp)
        if ret {
            return temp.length, temp.err
        }

        temp.c = lineReader.inByteBuf[lineReader.inOff]
        lineReader.inOff++

        if lineReader.readSkip(temp) {
            continue
        }

        if temp.c != '\n' && temp.c != '\r' {
            lineReader.readLine(temp)
        } else {
            cont := lineReader.readEOL(temp)
            if cont {
                continue
            }

            ret := lineReader.readNext(temp)
            if ret {
                return temp.length, temp.err
            }
        }
    }

    return -1, &LineReaderException{
        Message: "Unexpected Error",
    }
}

func (lineReader *LineReader) readInByteBuf(temp *readingLineTemp) bool {
    if lineReader.inOff >= lineReader.inLimit {
        n, err := lineReader.reader.Read(lineReader.inByteBuf)
        lineReader.inLimit = n
        lineReader.inOff = 0
        if nil != err || lineReader.inLimit <= 0 {
            if temp.length == 0 || temp.isCommentLine {
                temp.length = -1
                temp.err = err
                return true
            }
            if temp.precedingBackslash {
                temp.length--
            }
            return true
        }
    }
    return false
}

func (lineReader *LineReader) readSkip(temp *readingLineTemp) bool {
    if temp.skipLF {
        temp.skipLF = false
        if temp.c == '\n' {
            return true
        }
    }
    if temp.skipWhiteSpace {
        if temp.c == ' ' || temp.c == '\t' || temp.c == '\f' {
            return true
        }
        if !temp.appendedLineBegin && (temp.c == '\r' || temp.c == '\n') {
            return true
        }
        temp.skipWhiteSpace = false
        temp.appendedLineBegin = false
    }
    if temp.isNewLine {
        temp.isNewLine = false
        if temp.c == '#' || temp.c == '!' {
            temp.isCommentLine = true
            return true
        }
    }
    return false
}

func (lineReader *LineReader) readLine(temp *readingLineTemp) {
    lineReader.lineBuf[temp.length] = temp.c
    temp.length++
    if temp.length == len(lineReader.lineBuf) {
        newLength := temp.length * 2
        if newLength < 0 {
            newLength = int(^uint(0) >> 1)
        }
        buf := make([]byte, newLength)
        copy(buf, lineReader.lineBuf)
        lineReader.lineBuf = buf
    }
    // flip the preceding backslash flag
    if temp.c == '\\' {
        temp.precedingBackslash = !temp.precedingBackslash
    } else {
        temp.precedingBackslash = false
    }
}

func (lineReader *LineReader) readEOL(temp *readingLineTemp) bool {
    // reached EOL
    if temp.isCommentLine || temp.length == 0 {
        temp.isCommentLine = false
        temp.isNewLine = true
        temp.skipWhiteSpace = true
        temp.length = 0
        return true
    }
    return false
}

func (lineReader *LineReader) readNext(temp *readingLineTemp) bool {
    if lineReader.inOff >= lineReader.inLimit {
        n, err := lineReader.reader.Read(lineReader.inByteBuf)
        lineReader.inLimit = n
        lineReader.inOff = 0
        if nil != err || lineReader.inLimit <= 0 {
            if temp.precedingBackslash {
                temp.length--
            }
            return true
        }
    }
    if temp.precedingBackslash {
        temp.length -= 1
        // skip the leading whitespace characters in following line
        temp.skipWhiteSpace = true
        temp.appendedLineBegin = true
        temp.precedingBackslash = false
        if temp.c == '\r' {
            temp.skipLF = true
        }
    } else {
        return true
    }
    return false
}

type LineReaderException struct {
    Message string
}

func (e *LineReaderException) Error() string {
    return fmt.Sprintf("LineReaderException: %s", e.Message)
}
