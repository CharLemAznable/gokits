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
    "bytes"
    "fmt"
    "io"
    "sort"
    "strconv"
    "strings"
)

// A YamlNode is a YAML YamlNode which can be a YamlMap, YamlList or YamlScalar.
type YamlNode interface {
    write(io.Writer, int, int)
}

// A YamlMap is a YAML Mapping which maps Strings to Nodes.
type YamlMap map[string]YamlNode

var _ YamlNode = YamlMap{}

// Key returns the value associeted with the key in the map.
func (node YamlMap) Key(key string) YamlNode {
    return node[key]
}

func (node YamlMap) write(out io.Writer, firstind, nextind int) {
    indent := bytes.Repeat([]byte{' '}, nextind)
    ind := firstind

    width := 0
    var scalarkeys []string
    var objectkeys []string
    for key, value := range node {
        if _, ok := value.(YamlScalar); ok {
            if swid := len(key); swid > width {
                width = swid
            }
            scalarkeys = append(scalarkeys, key)
            continue
        }
        objectkeys = append(objectkeys, key)
    }
    sort.Strings(scalarkeys)
    sort.Strings(objectkeys)

    for _, key := range scalarkeys {
        value := node[key].(YamlScalar)
        _, _ = out.Write(indent[:ind])
        _, _ = fmt.Fprintf(out, "%-*s %s\n", width+1, key+":", string(value))
        ind = nextind
    }
    for _, key := range objectkeys {
        _, _ = out.Write(indent[:ind])
        if node[key] == nil {
            _, _ = fmt.Fprintf(out, "%s: <nil>\n", key)
            continue
        }
        _, _ = fmt.Fprintf(out, "%s:\n", key)
        ind = nextind
        node[key].write(out, ind+2, ind+2)
    }
}

// A YamlList is a YAML Sequence of Nodes.
type YamlList []YamlNode

var _ YamlNode = YamlList{}

// GetString the number of items in the YamlList.
func (node YamlList) Len() int {
    return len(node)
}

// GetString the idx'th item from the YamlList.
func (node YamlList) Item(idx int) YamlNode {
    if idx >= 0 && idx < len(node) {
        return node[idx]
    }
    return nil
}

func (node YamlList) write(out io.Writer, firstind, nextind int) {
    indent := bytes.Repeat([]byte{' '}, nextind)
    ind := firstind

    for _, value := range node {
        _, _ = out.Write(indent[:ind])
        _, _ = fmt.Fprint(out, "- ")
        ind = nextind
        value.write(out, 0, ind+2)
    }
}

// A YamlScalar is a YAML YamlScalar.
type YamlScalar string

var _ YamlNode = YamlScalar("")

// String returns the string represented by this YamlScalar.
func (node YamlScalar) String() (string, error) { return string(node), nil }

func (node YamlScalar) Int() (int64, error) {
    str, err := node.String()
    if err != nil {
        return 0, err
    }

    i, err := strconv.ParseInt(str, 10, 64)
    if err != nil {
        return 0, err
    }

    return i, nil
}

func (node YamlScalar) Bool() (bool, error) {
    str, err := node.String()
    if err != nil {
        return false, err
    }

    b, err := strconv.ParseBool(str)
    if err != nil {
        return false, err
    }

    return b, nil
}

func (node YamlScalar) write(out io.Writer, ind, _ int) {
    _, _ = fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", ind), string(node))
}

// YamlRender returns a string of the node as a YAML document.  Note that
// Scalars will have a newline appended if they are rendered directly.
func YamlRender(node YamlNode) string {
    buf := bytes.NewBuffer(nil)
    node.write(buf, 0, 0)
    return buf.String()
}
