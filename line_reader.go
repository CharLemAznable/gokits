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
    length := 0
    var c byte = 0
    skipWhiteSpace := true
    isCommentLine := false
    isNewLine := true
    appendedLineBegin := false
    precedingBackslash := false
    skipLF := false

    for true {
        ret, lth, err := lineReader.readInByteBuf(
            length, isCommentLine, precedingBackslash)
        if ret {
            return lth, err
        } else {
            length = lth
        }

        c = lineReader.inByteBuf[lineReader.inOff]
        lineReader.inOff++

        if lineReader.readSkip(&skipLF, &c,
            &skipWhiteSpace, &appendedLineBegin,
            &isNewLine, &isCommentLine) {
            continue
        }

        if c != '\n' && c != '\r' {
            lineReader.readLine(&length, c, &precedingBackslash)
        } else {
            cont, lth := lineReader.readEOL(&isCommentLine,
                &length, &isNewLine, &skipWhiteSpace)
            length = lth
            if cont {
                continue
            }

            ret, lth, err := lineReader.readNext(&precedingBackslash, &length,
                &skipWhiteSpace, &appendedLineBegin, c, &skipLF)
            if ret {
                return lth, err
            } else {
                length = lth
            }
        }
    }

    return -1, &LineReaderException{
        Message: "Unexpected Error",
    }
}

func (lineReader *LineReader) readInByteBuf(length int, isCommentLine bool, precedingBackslash bool) (bool, int, error) {
    if lineReader.inOff >= lineReader.inLimit {
        n, err := lineReader.reader.Read(lineReader.inByteBuf)
        lineReader.inLimit = n
        lineReader.inOff = 0
        if nil != err || lineReader.inLimit <= 0 {
            if length == 0 || isCommentLine {
                return true, -1, err
            }
            if precedingBackslash {
                length--
            }
            return true, length, nil
        }
    }
    return false, length, nil
}

func (lineReader *LineReader) readSkip(skipLF *bool, c *byte,
    skipWhiteSpace *bool, appendedLineBegin *bool,
    isNewLine *bool, isCommentLine *bool) bool {

    if *skipLF {
        *skipLF = false
        if *c == '\n' {
            return true
        }
    }
    if *skipWhiteSpace {
        if *c == ' ' || *c == '\t' || *c == '\f' {
            return true
        }
        if !*appendedLineBegin && (*c == '\r' || *c == '\n') {
            return true
        }
        *skipWhiteSpace = false
        *appendedLineBegin = false
    }
    if *isNewLine {
        *isNewLine = false
        if *c == '#' || *c == '!' {
            *isCommentLine = true
            return true
        }
    }
    return false
}

func (lineReader *LineReader) readLine(length *int, c byte, precedingBackslash *bool) {
    lineReader.lineBuf[*length] = c
    *length++
    if *length == len(lineReader.lineBuf) {
        newLength := *length * 2
        if newLength < 0 {
            newLength = int(^uint(0) >> 1)
        }
        buf := make([]byte, newLength)
        copy(buf, lineReader.lineBuf)
        lineReader.lineBuf = buf
    }
    // flip the preceding backslash flag
    if c == '\\' {
        *precedingBackslash = !*precedingBackslash
    } else {
        *precedingBackslash = false
    }
}

func (lineReader *LineReader) readEOL(isCommentLine *bool, length *int,
    isNewLine *bool, skipWhiteSpace *bool) (bool, int) {
    // reached EOL
    if *isCommentLine || *length == 0 {
        *isCommentLine = false
        *isNewLine = true
        *skipWhiteSpace = true
        *length = 0
        return false, *length
    }
    return true, *length
}

func (lineReader *LineReader) readNext(precedingBackslash *bool, length *int,
    skipWhiteSpace *bool, appendedLineBegin *bool, c byte, skipLF *bool) (bool, int, error) {
    if lineReader.inOff >= lineReader.inLimit {
        n, err := lineReader.reader.Read(lineReader.inByteBuf)
        lineReader.inLimit = n
        lineReader.inOff = 0
        if nil != err || lineReader.inLimit <= 0 {
            if *precedingBackslash {
                *length--
            }
            return true, *length, err
        }
    }
    if *precedingBackslash {
        *length -= 1
        // skip the leading whitespace characters in following line
        *skipWhiteSpace = true
        *appendedLineBegin = true
        *precedingBackslash = false
        if c == '\r' {
            *skipLF = true
        }
    } else {
        return true, *length, nil
    }
    return false, *length, nil
}

type LineReaderException struct {
    Message string
}

func (e *LineReaderException) Error() string {
    return fmt.Sprintf("LineReaderException: %s", e.Message)
}
