// Copyright 2013 Google, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gokits

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "strings"
)

// YamlParse returns a root-level YamlNode parsed from the lines read from r.  In
// general, this will be done for you by one of the YamlFile constructors.
func YamlParse(r io.Reader) (node YamlNode, err error) {
    lb := &lineBuffer{
        Reader: bufio.NewReader(r),
    }

    defer func() {
        if r := recover(); r != nil {
            switch r := r.(type) {
            case error:
                err = r
            case string:
                err = errors.New(r)
            default:
                err = fmt.Errorf("%v", r)
            }
        }
    }()

    node = parseNode(lb, 0, nil)
    return
}

// Supporting types and constants

//noinspection GoUnusedConst
const (
    typUnknown = iota
    typSequence
    typMapping
    typScalar
)

var typNames = []string{
    "Unknown", "Sequence", "Mapping", "YamlScalar",
}

type lineReader interface {
    Next(minIndent int) *indentedLine
}

type indentedLine struct {
    lineno int
    indent int
    line   []byte
}

func (line *indentedLine) loopIndent() {
    for _, ch := range line.line {
        switch ch {
        case ' ':
            line.indent += 1
            continue
        default:
        }
        break
    }
    line.line = line.line[line.indent:]
}

func (line *indentedLine) String() string {
    return fmt.Sprintf("%2d: %s%s", line.indent,
        strings.Repeat(" ", 0*line.indent), string(line.line))
}

func parseNode(r lineReader, ind int, initial YamlNode) (node YamlNode) {
    first := true
    node = initial

    // read lines
    for {
        line := r.Next(ind)
        if line == nil {
            break
        }

        if len(line.line) == 0 {
            continue
        }

        node = parseNodeLine(r, &ind, &first, node, line)
    }
    return
}

func parseNodeLine(r lineReader, ind *int, first *bool, node YamlNode, line *indentedLine) YamlNode {
    if *first {
        *ind = line.indent
        *first = false
    }

    types, pieces := parseNodeInlineValue(r, line)

    return parsePrevNode(r, line, types, pieces, node)
}

func parseNodeInlineValue(r lineReader, line *indentedLine) ([]int, []string) {
    var types []int
    var pieces []string

    var inlineValue func([]byte)
    inlineValue = func(partial []byte) {
        // This can be a for loop now
        vtyp, brk := getType(partial)
        begin, end := partial[:brk], partial[brk:]

        if vtyp == typMapping {
            end = end[1:]
        }
        end = bytes.TrimLeft(end, " ")

        switch vtyp {
        case typScalar:
            types = append(types, typScalar)
            pieces = append(pieces, string(end))
            return
        case typMapping:
            types = append(types, typMapping)
            pieces = append(pieces, strings.TrimSpace(string(begin)))

            trimmed := bytes.TrimSpace(end)
            if len(trimmed) == 1 && trimmed[0] == '|' {
                text := parseTypMappingLoog(r)
                types = append(types, typScalar)
                pieces = append(pieces, string(text))
                return
            }
            inlineValue(end)
        case typSequence:
            types = append(types, typSequence)
            pieces = append(pieces, "-")
            inlineValue(end)
        }
    }

    inlineValue(line.line)
    return types, pieces
}

func parseTypMappingLoog(r lineReader) string {
    text := ""
    for {
        l := r.Next(1)
        if l == nil {
            break
        }

        s := string(l.line)
        s = strings.TrimSpace(s)
        if len(s) == 0 {
            break
        }
        text = text + "\n" + s
    }
    return text
}

func parsePrevNode(r lineReader, line *indentedLine, types []int, pieces []string, node YamlNode) YamlNode {
    var prev YamlNode

    // Nest inlines
    for len(types) > 0 {
        prev = parsePrevNodeLoop(r, line, &types, &pieces, node, prev)
    }

    return prev
}

func parsePrevNodeLoop(r lineReader, line *indentedLine,
    types *[]int, pieces *[]string, node YamlNode, prev YamlNode) YamlNode {

    last := len(*types) - 1
    typ, piece := (*types)[last], (*pieces)[last]

    var current YamlNode
    if last == 0 {
        current = node
    }
    // child := parseNode(r, line.indent+1, typUnknown) allow scalar only

    // Add to current node
    switch typ {
    case typScalar: // last will be == nil
        if _, ok := current.(YamlScalar); current != nil && !ok {
            panic("cannot append scalar to non-scalar node")
        }
        if current != nil {
            current = YamlScalar(piece) + " " + current.(YamlScalar)
            break
        }
        current = YamlScalar(piece)
    case typMapping:
        var mapNode YamlMap
        var ok bool
        var child YamlNode

        // GetString the current map, if there is one
        if mapNode, ok = current.(YamlMap); current != nil && !ok {
            _ = current.(YamlMap) // panic
        } else if current == nil {
            mapNode = make(YamlMap)
        }

        if _, inlineMap := prev.(YamlScalar); inlineMap && last > 0 {
            current = YamlMap{
                piece: prev,
            }
            break
        }

        child = parseNode(r, line.indent+1, prev)
        mapNode[piece] = child
        current = mapNode

    case typSequence:
        var listNode YamlList
        var ok bool
        var child YamlNode

        // GetString the current list, if there is one
        if listNode, ok = current.(YamlList); current != nil && !ok {
            _ = current.(YamlList) // panic
        } else if current == nil {
            listNode = make(YamlList, 0)
        }

        if _, inlineList := prev.(YamlScalar); inlineList && last > 0 {
            current = YamlList{
                prev,
            }
            break
        }

        child = parseNode(r, line.indent+1, prev)
        listNode = append(listNode, child)
        current = listNode
    }

    if last < 0 {
        last = 0
    }

    *types = (*types)[:last]
    *pieces = (*pieces)[:last]
    return current
}

func getType(line []byte) (typ, split int) {
    if len(line) == 0 {
        return
    }

    if line[0] == '-' {
        typ = typSequence
        split = 1
        return
    }

    typ = typScalar

    if line[0] == ' ' || line[0] == '"' {
        return
    }

    // the first character is real
    // need to iterate past the first word
    // things like "foo:" and "foo :" are mappings
    // everything else is a scalar

    idx := bytes.IndexAny(line, " \":")
    if idx < 0 {
        return
    }

    if line[idx] == '"' {
        return
    }

    if line[idx] == ':' {
        typ = typMapping
        split = idx
    } else if line[idx] == ' ' {
        seekColon(idx, line, &typ, &split)
    }

    if typ == typMapping && split+1 < len(line) && line[split+1] != ' ' {
        typ = typScalar
        split = 0
    }

    return
}

func seekColon(idx int, line []byte, typ *int, split *int) {
    // we have a space
    // need to see if its all spaces until a :
    for i := idx; i < len(line); i++ {
        switch ch := line[i]; ch {
        case ' ':
            continue
        case ':':
            // only split on colons followed by a space
            if i+1 < len(line) && line[i+1] != ' ' {
                continue
            }

            *typ = typMapping
            *split = i
            break
        default:
            break
        }
    }
}

// lineReader implementations

type lineBuffer struct {
    *bufio.Reader
    readLines int
    pending   *indentedLine
}

func (lb *lineBuffer) Next(min int) (next *indentedLine) {
    if lb.pending == nil {
        var (
            read []byte
            more bool
            err  error
        )

        l := new(indentedLine)
        l.lineno = lb.readLines
        more = true
        for more {
            read, more, err = lb.ReadLine()
            if err != nil {
                if err == io.EOF {
                    return nil
                }
                panic(err)
            }
            l.line = append(l.line, read...)
        }
        lb.readLines++

        l.loopIndent()
        // Ignore blank lines and comments.
        if len(l.line) == 0 || l.line[0] == '#' {
            return lb.Next(min)
        }

        lb.pending = l
    }
    next = lb.pending
    if next.indent < min {
        return nil
    }
    lb.pending = nil
    return
}

//noinspection GoUnusedType
type lineSlice []*indentedLine

func (ls *lineSlice) Next(min int) (next *indentedLine) {
    if len(*ls) == 0 {
        return nil
    }
    next = (*ls)[0]
    if next.indent < min {
        return nil
    }
    *ls = (*ls)[1:]
    return
}

func (ls *lineSlice) Push(line *indentedLine) {
    *ls = append(*ls, line)
}
