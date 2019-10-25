package gokits

import (
    "bufio"
    "bytes"
    "errors"
    "io"
    "runtime"
    "sync"
    "time"
)

type Properties struct {
    Hashtable
    mutex    sync.RWMutex
    defaults *Properties
}

func NewProperties() *Properties {
    return &Properties{Hashtable: *NewHashtable()}
}

func NewPropertiesDefault(defaults *Properties) *Properties {
    properties := NewProperties()
    properties.defaults = defaults
    return properties
}

func (properties *Properties) SetProperty(key string, value string) {
    properties.mutex.Lock()
    defer properties.mutex.Unlock()
    properties.Put(key, value)
}

func (properties *Properties) GetProperty(key string) string {
    oval := properties.Get(key)
    if sval, ok := oval.(string); ok {
        return sval
    }
    if properties.defaults != nil {
        return properties.defaults.GetProperty(key)
    }
    return ""
}

func (properties *Properties) GetPropertyDefault(key, defaultValue string) string {
    value := properties.GetProperty(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func (properties *Properties) PropertyNames() []interface{} {
    hashtable := NewHashtable()
    properties.enumerate(hashtable)

    return hashtable.Keys()
}

func (properties *Properties) StringPropertyNames() []string {
    hashtable := NewHashtable()
    properties.enumerateStringProperties(hashtable)

    var set []string
    for _, key := range hashtable.Keys() {
        set = append(set, key.(string))
    }

    return set
}

func (properties *Properties) List(out io.Writer) {
    _, _ = out.Write(append([]byte("-- listing properties --"), newLine()...))

    hashtable := NewHashtable()
    properties.enumerate(hashtable)
    for _, key := range hashtable.Keys() {
        val := hashtable.Get(key)
        skey := key.(string)
        sval := val.(string)

        if len(sval) > 40 {
            sval = string([]byte(sval)[:37]) + "..."
        }
        _, _ = out.Write(append([]byte(skey+" = "+sval), newLine()...))
    }

    _, _ = out.Write(append([]byte("-- listing properties end --"), newLine()...))
}

func (properties *Properties) ToMap() map[interface{}]interface{} {
    var m = make(map[interface{}]interface{})
    keys := properties.Keys()
    for _, key := range keys {
        m[key] = properties.Get(key)
    }
    return m
}

func (properties *Properties) enumerate(hashtable *Hashtable) {
    properties.mutex.Lock()
    defer properties.mutex.Unlock()

    if properties.defaults != nil {
        properties.defaults.enumerate(hashtable)
    }

    for _, key := range properties.Hashtable.Keys() {
        hashtable.Put(key.(string), properties.Hashtable.Get(key))
    }
}

func (properties *Properties) enumerateStringProperties(hashtable *Hashtable) {
    properties.mutex.Lock()
    defer properties.mutex.Unlock()

    if properties.defaults != nil {
        properties.defaults.enumerateStringProperties(hashtable)
    }

    for _, key := range properties.Hashtable.Keys() {
        val := properties.Hashtable.Get(key)
        if _, ok := key.(string); ok {
            if _, ok = val.(string); ok {
                hashtable.Put(key.(string), val.(string))
            }
        }
    }
}

func (properties *Properties) Load(reader io.Reader) error {
    properties.mutex.Lock()
    defer properties.mutex.Unlock()

    return properties.load0(NewLineReader(reader))
}

type load0Temp struct {
    convtBuf                   []byte
    limit, keyLen, valueStart  int
    c                          byte
    hasSep, precedingBackslash bool
}

func (properties *Properties) load0(lr *LineReader) error {
    temp := &load0Temp{
        convtBuf: make([]byte, 1024),
    }

    for temp.limit, _ = lr.ReadLine(); temp.limit >= 0; temp.limit, _ = lr.ReadLine() {
        temp.c = 0
        temp.keyLen = 0
        temp.valueStart = temp.limit
        temp.hasSep = false

        properties.load0Loop(lr, temp)

        key, err := properties.loadConvert(lr.lineBuf, 0, temp.keyLen, temp.convtBuf)
        if err != nil {
            return err
        }
        value, err := properties.loadConvert(lr.lineBuf, temp.valueStart, temp.limit-temp.valueStart, temp.convtBuf)
        if err != nil {
            return err
        }
        properties.Put(key, value)
    }

    return nil
}

func (properties *Properties) load0Loop(lr *LineReader, temp *load0Temp) {
    // fmt.Println("line=<" + string(lr.(*lineReader).lineBuf[:limit]) + ">")
    temp.precedingBackslash = false
    for temp.keyLen < temp.limit {
        temp.c = lr.lineBuf[temp.keyLen]
        // need check if escaped.
        if (temp.c == '=' || temp.c == ':') && !temp.precedingBackslash {
            temp.valueStart = temp.keyLen + 1
            temp.hasSep = true
            break
        } else if (temp.c == ' ' || temp.c == '\t' || temp.c == '\f') && !temp.precedingBackslash {
            temp.valueStart = temp.keyLen + 1
            break
        }
        if temp.c == '\\' {
            temp.precedingBackslash = !temp.precedingBackslash
        } else {
            temp.precedingBackslash = false
        }
        temp.keyLen++
    }
    for temp.valueStart < temp.limit {
        temp.c = lr.lineBuf[temp.valueStart]
        if temp.c != ' ' && temp.c != '\t' && temp.c != '\f' {
            if !temp.hasSep && (temp.c == '=' || temp.c == ':') {
                temp.hasSep = true
            } else {
                break
            }
        }
        temp.valueStart++
    }
}

func (properties *Properties) loadConvert(in []byte, off, length int, convtBuf []byte) (string, error) {
    convtBuf = properties.checkBufLen(length, convtBuf)

    var aChar byte
    out := convtBuf
    outLen := 0
    end := off + length

    for off < end {
        aChar = in[off]
        off++
        if aChar == '\\' {
            aChar = in[off]
            off++
            if aChar == 'u' {
                value, err := properties.escapeUxxx(in, &off, &aChar)
                if nil != err {
                    return "", errors.New("malformed \\uxxxx encoding")
                }
                out[outLen] = byte(value)
                outLen++
            } else {
                out[outLen] = properties.escapeFormatByte(aChar)
                outLen++
            }
        } else {
            out[outLen] = aChar
            outLen++
        }
    }

    return string(out[:outLen]), nil
}

func (properties *Properties) checkBufLen(length int, convtBuf []byte) []byte {
    if len(convtBuf) < length {
        newLen := length * 2
        if newLen < 0 {
            newLen = int(^uint(0) >> 1)
        }
        return make([]byte, newLen)
    }
    return convtBuf
}

func (properties *Properties) escapeUxxx(in []byte, off *int, aChar *byte) (int, error) {
    // Read the xxxx
    value := 0

    for i := 0; i < 4; i++ {
        *aChar = in[*off]
        *off++
        switch *aChar {
        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
            value = (value << 4) + int(*aChar) - '0'
        case 'a', 'b', 'c', 'd', 'e', 'f':
            value = (value << 4) + 10 + int(*aChar) - 'a'
        case 'A', 'B', 'C', 'D', 'E', 'F':
            value = (value << 4) + 10 + int(*aChar) - 'A'
        default:
            return 0, errors.New("malformed \\uxxxx encoding")
        }
    }
    return value, nil
}

func (properties *Properties) escapeFormatByte(aChar byte) byte {
    if aChar == 't' {
        return '\t'
    } else
    if aChar == 'r' {
        return '\r'
    } else
    if aChar == 'n' {
        return '\n'
    } else
    if aChar == 'f' {
        return '\f'
    }
    return aChar
}

func (properties *Properties) Save(writer io.Writer, comments string) {
    _ = properties.Store(writer, comments)
}

func (properties *Properties) Store(writer io.Writer, comments string) error {
    return properties.store0(bufio.NewWriter(writer), comments, true)
}

func (properties *Properties) store0(bw *bufio.Writer, comments string, escUnicode bool) (err error) {
    if comments != "" {
        if err = writeComments(bw, comments); err != nil {
            return err
        }
    }
    if _, err = bw.WriteString("# " + time.Now().Format(time.UnixDate)); err != nil {
        return err
    }
    if _, err = bw.Write(newLine()); err != nil {
        return err
    }

    properties.mutex.Lock()
    defer properties.mutex.Unlock()

    for _, key := range properties.Keys() {
        val := properties.Get(key)

        skey := key.(string)
        sval := val.(string)

        skey = properties.saveConvert(skey, true, escUnicode)
        // No need to escape embedded and trailing spaces for value, hence
        // pass false to flag.
        sval = properties.saveConvert(sval, false, escUnicode)
        if _, err = bw.WriteString(skey + " = " + sval); err != nil {
            return err
        }
        if _, err = bw.Write(newLine()); err != nil {
            return err
        }
    }

    return bw.Flush()
}

func writeComments(bw *bufio.Writer, comments string) (err error) {
    if err = bw.WriteByte('#'); err != nil {
        return err
    }
    length := len(comments)
    current := 0
    last := 0
    uu := make([]byte, 6)
    uu[0] = '\\'
    uu[1] = 'u'
    for current < length {
        c := comments[current]
        if c > '\u00ff' || c == '\n' || c == '\r' {
            if last != current {
                if _, err = bw.Write([]byte(comments)[last:current]); err != nil {
                    return err
                }
            }
            if c > '\u00ff' {
                uu[2] = toHex(int(c>>12) & 0xf)
                uu[3] = toHex(int(c>>8) & 0xf)
                uu[4] = toHex(int(c>>4) & 0xf)
                uu[5] = toHex(int(c) & 0xf)
                if _, err = bw.Write(uu); err != nil {
                    return err
                }
            } else {
                if _, err = bw.Write(newLine()); err != nil {
                    return err
                }
                if c == '\r' && current != length-1 && comments[current+1] == '\n' {
                    current++
                }
                if current == length-1 || comments[current+1] != '#' && comments[current+1] != '!' {
                    if err = bw.WriteByte('#'); err != nil {
                        return err
                    }
                }
            }
            last = current + 1
        }
        current++
    }

    if last != current {
        if _, err = bw.Write([]byte(comments)[last:current]); err != nil {
            return err
        }
    }
    if _, err = bw.Write(newLine()); err != nil {
        return err
    }

    return nil
}

var hexDigit = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

func toHex(nibble int) byte {
    return hexDigit[nibble&0xF]
}

func newLine() []byte {
    const (
        CR = '\r'
        LF = '\n'
    )
    switch runtime.GOOS {
    case "windows":
        return []byte{CR, LF}
    case "linux":
        fallthrough
    default:
        return []byte{LF}
    }
}

func (properties *Properties) saveConvert(theString string, escapeSpace, escapeUnicode bool) string {
    length := len(theString)
    bufLen := length * 2
    if bufLen < 0 {
        bufLen = int(^uint(0) >> 1)
    }

    var outBuffer bytes.Buffer

    for x := 0; x < length; x++ {
        aChar := theString[x]
        // Handle common case first, selecting largest block that
        // avoids the specials below
        if (aChar > 61) && (aChar < 127) {
            writeEscapeBackSlash(&outBuffer, aChar)
            continue
        }
        switch aChar {
        case ' ':
            if x == 0 || escapeSpace {
                outBuffer.WriteByte('\\')
            }
            outBuffer.WriteByte(' ')
        case '\t':
            outBuffer.WriteByte('\\')
            outBuffer.WriteByte('t')
        case '\n':
            outBuffer.WriteByte('\\')
            outBuffer.WriteByte('n')
        case '\r':
            outBuffer.WriteByte('\\')
            outBuffer.WriteByte('r')
        case '\f':
            outBuffer.WriteByte('\\')
            outBuffer.WriteByte('f')
        case '=':
            fallthrough
        case ':', '#', '!':
            outBuffer.WriteByte('\\')
            outBuffer.WriteByte(aChar)
        default:
            writeDefault(&outBuffer, aChar, escapeUnicode)
        }
    }

    return outBuffer.String()
}

func writeEscapeBackSlash(outBuffer *bytes.Buffer, aChar uint8) {
    if aChar == '\\' {
        outBuffer.WriteByte('\\')
        outBuffer.WriteByte('\\')
        return
    }
    outBuffer.WriteByte(aChar)
}

func writeDefault(outBuffer *bytes.Buffer, aChar uint8, escapeUnicode bool) {
    if ((aChar < 0x0020) || (aChar > 0x007e)) && escapeUnicode {
        outBuffer.WriteByte('\\')
        outBuffer.WriteByte('u')
        outBuffer.WriteByte(toHex(int(aChar>>12) & 0xF))
        outBuffer.WriteByte(toHex(int(aChar>>8) & 0xF))
        outBuffer.WriteByte(toHex(int(aChar>>4) & 0xF))
        outBuffer.WriteByte(toHex(int(aChar) & 0xF))
    } else {
        outBuffer.WriteByte(aChar)
    }
}
