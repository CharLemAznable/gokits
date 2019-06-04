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
    // readLine() (int, error)
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
    var length = 0
    var c byte = 0
    var skipWhiteSpace = true
    var isCommentLine = false
    var isNewLine = true
    var appendedLineBegin = false
    var precedingBackslash = false
    var skipLF = false

    for true {
        if lineReader.inOff >= lineReader.inLimit {
            n, err := lineReader.reader.Read(lineReader.inByteBuf)
            lineReader.inLimit = n
            lineReader.inOff = 0
            if nil != err || lineReader.inLimit <= 0 {
                if length == 0 || isCommentLine {
                    return -1, err
                }
                if precedingBackslash {
                    length--
                }
                return length, nil
            }
        }

        c = lineReader.inByteBuf[lineReader.inOff]
        lineReader.inOff++

        if skipLF {
            skipLF = false
            if c == '\n' {
                continue
            }
        }
        if skipWhiteSpace {
            if c == ' ' || c == '\t' || c == '\f' {
                continue
            }
            if !appendedLineBegin && (c == '\r' || c == '\n') {
                continue
            }
            skipWhiteSpace = false
            appendedLineBegin = false
        }
        if isNewLine {
            isNewLine = false
            if c == '#' || c == '!' {
                isCommentLine = true
                continue
            }
        }

        if c != '\n' && c != '\r' {
            lineReader.lineBuf[length] = c
            length++
            if length == len(lineReader.lineBuf) {
                var newLength = length * 2
                if newLength < 0 {
                    newLength = int(^uint(0) >> 1)
                }
                var buf = make([]byte, newLength)
                copy(buf, lineReader.lineBuf)
                lineReader.lineBuf = buf
            }
            // flip the preceding backslash flag
            if c == '\\' {
                precedingBackslash = !precedingBackslash
            } else {
                precedingBackslash = false
            }
        } else {
            // reached EOL
            if isCommentLine || length == 0 {
                isCommentLine = false
                isNewLine = true
                skipWhiteSpace = true
                length = 0
                continue
            }
            if lineReader.inOff >= lineReader.inLimit {
                n, err := lineReader.reader.Read(lineReader.inByteBuf)
                lineReader.inLimit = n
                lineReader.inOff = 0
                if nil != err || lineReader.inLimit <= 0 {
                    if precedingBackslash {
                        length--
                    }
                    return length, err
                }
            }
            if precedingBackslash {
                length -= 1
                // skip the leading whitespace characters in following line
                skipWhiteSpace = true
                appendedLineBegin = true
                precedingBackslash = false
                if c == '\r' {
                    skipLF = true
                }
            } else {
                return length, nil
            }
        }
    }

    return -1, &LineReaderException{
        Message: "Unexpected Error",
    }
}

type LineReaderException struct {
    Message string
}

func (e *LineReaderException) Error() string {
    return fmt.Sprintf("LineReaderException: %s", e.Message)
}
