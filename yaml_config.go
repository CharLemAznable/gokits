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
    "os"
    "strconv"
    "strings"
)

// A YamlFile represents the top-level YAML node found in a file.  It is intended
// for use as a configuration file.
type YamlFile struct {
    Root YamlNode

    // Add a cache?
}

func ReadYamlString(yamlconf string) (*YamlFile, error) {
    buf := bytes.NewBufferString(yamlconf)

    var err error
    f := new(YamlFile)
    f.Root, err = YamlParse(buf)
    if err != nil {
        return nil, err
    }

    return f, nil
}

// ReadYamlFile reads a YAML configuration file from the given filename.
func ReadYamlFile(filename string) (*YamlFile, error) {
    fin, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer func() { _ = fin.Close() }()

    f := new(YamlFile)
    f.Root, err = YamlParse(fin)
    if err != nil {
        return nil, err
    }

    return f, nil
}

// ConfigYamlString reads a YAML configuration from a static string.  If an error is
// found, it will panic.  This is a utility function and is intended for use in
// initializers.
//noinspection GoUnusedExportedFunction
func ConfigYamlString(yamlconf string) *YamlFile {
    f, err := ReadYamlString(yamlconf)
    if err != nil {
        panic(err)
    }

    return f
}

// ConfigYamlFile reads a YAML configuration file from the given filename and
// panics if an error is found.  This is a utility function and is intended for
// use in initializers.
//noinspection GoUnusedExportedFunction
func ConfigYamlFile(filename string) *YamlFile {
    f, err := ReadYamlFile(filename)
    if err != nil {
        panic(err)
    }
    return f
}

func (f *YamlFile) GetNode(spec string) (YamlNode, error) {
    return YamlChild(f.Root, spec)
}

func (f *YamlFile) GetScalar(spec string) (YamlScalar, error) {
    return ScalarChild(f.Root, spec)
}

func (f *YamlFile) GetString(spec string) (string, error) {
    return StringChild(f.Root, spec)
}

func (f *YamlFile) GetInt(spec string) (int64, error) {
    return IntChild(f.Root, spec)
}

func (f *YamlFile) GetBool(spec string) (bool, error) {
    return BoolChild(f.Root, spec)
}

func (f *YamlFile) GetList(spec string) (YamlList, error) {
    return ListChild(f.Root, spec)
}

func (f *YamlFile) GetListCount(spec string) (int, error) {
    return ListChildCount(f.Root, spec)
}

func (f *YamlFile) GetMap(spec string) (YamlMap, error) {
    return MapChild(f.Root, spec)
}

const rootSpec = "root"

func (f *YamlFile) RootScalar() (YamlScalar, error) {
    return Scalar(f.Root, rootSpec)
}

func (f *YamlFile) RootString() (string, error) {
    return String(f.Root, rootSpec)
}

func (f *YamlFile) RootInt() (int64, error) {
    return Int(f.Root, rootSpec)
}

func (f *YamlFile) RootBool() (bool, error) {
    return Bool(f.Root, rootSpec)
}

func (f *YamlFile) RootList() (YamlList, error) {
    return List(f.Root, rootSpec)
}

func (f *YamlFile) RooListCount() (int, error) {
    return ListCount(f.Root, rootSpec)
}

func (f *YamlFile) RootMap() (YamlMap, error) {
    return Map(f.Root, rootSpec)
}

func ScalarChild(root YamlNode, spec string) (YamlScalar, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return "", err
    }

    return Scalar(node, spec)
}

func StringChild(root YamlNode, spec string) (string, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return "", err
    }

    return String(node, spec)
}

func IntChild(root YamlNode, spec string) (int64, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return 0, err
    }

    return Int(node, spec)
}

func BoolChild(root YamlNode, spec string) (bool, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return false, err
    }

    return Bool(node, spec)
}

func ListChild(root YamlNode, spec string) (YamlList, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return nil, err
    }

    return List(node, spec)
}

func ListChildCount(root YamlNode, spec string) (int, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return -1, err
    }

    return ListCount(node, spec)
}

func MapChild(root YamlNode, spec string) (YamlMap, error) {
    node, err := YamlChild(root, spec)
    if err != nil {
        return nil, err
    }

    return Map(node, spec)
}

func Scalar(node YamlNode, spec string) (YamlScalar, error) {
    sc, ok := node.(YamlScalar)
    if !ok {
        return "", &YamlNodeTypeMismatch{
            Full:     spec,
            Spec:     spec,
            Token:    "$",
            Expected: "gokits.YamlScalar",
            Node:     node,
        }
    }
    return sc, nil
}

func String(node YamlNode, spec string) (string, error) {
    sc, err := Scalar(node, spec)
    if err != nil {
        return "", err
    }

    return sc.String()
}

func Int(node YamlNode, spec string) (int64, error) {
    sc, err := Scalar(node, spec)
    if err != nil {
        return 0, err
    }

    return sc.Int()
}

func Bool(node YamlNode, spec string) (bool, error) {
    sc, err := Scalar(node, spec)
    if err != nil {
        return false, err
    }

    return sc.Bool()
}

func List(node YamlNode, spec string) (YamlList, error) {
    ls, ok := node.(YamlList)
    if !ok {
        return nil, &YamlNodeTypeMismatch{
            Full:     spec,
            Spec:     spec,
            Token:    "$",
            Expected: "gokits.YamlList",
            Node:     node,
        }
    }
    return ls, nil
}

func ListCount(node YamlNode, spec string) (int, error) {
    lst, err := List(node, spec)
    if err != nil {
        return -1, err
    }

    return lst.Len(), nil
}

func Map(node YamlNode, spec string) (YamlMap, error) {
    mp, ok := node.(YamlMap)
    if !ok {
        return nil, &YamlNodeTypeMismatch{
            Full:     spec,
            Spec:     spec,
            Token:    "$",
            Expected: "gokits.YamlMap",
            Node:     node,
        }
    }
    return mp, nil
}

// YamlChild retrieves a child node from the specified node as follows:
//   .mapkey   - GetString the key 'mapkey' of the YamlNode, which must be a YamlMap
//   [idx]     - Choose the index from the current YamlNode, which must be a YamlList
//
// The above selectors may be applied recursively, and each successive selector
// applies to the result of the previous selector.  For convenience, a "." is
// implied as the first character if the first character is not a "." or "[".
// The node tree is walked from the given node, considering each token of the
// above format.  If a node along the evaluation path is not found, an error is
// returned. If a node is not the proper type, an error is returned.  If the
// final node is not a YamlScalar, an error is returned.
func YamlChild(root YamlNode, spec string) (YamlNode, error) {
    if len(spec) == 0 {
        return root, nil
    }

    if first := spec[0]; first != '.' && first != '[' {
        spec = "." + spec
    }

    var recur func(YamlNode, string, string) (YamlNode, error)
    recur = func(n YamlNode, last, s string) (YamlNode, error) {

        if len(s) == 0 {
            return n, nil
        }

        if n == nil {
            return nil, &YamlNodeNotFound{
                Full: spec,
                Spec: last,
            }
        }

        // Extract the next token
        delim := 1 + strings.IndexAny(s[1:], ".[")
        if delim <= 0 {
            delim = len(s)
        }
        tok := s[:delim]
        remain := s[delim:]

        switch s[0] {
        case '[':
            return parseYamlList(n, spec, last, tok, remain, recur)
        default:
            return parseDefault(n, spec, last, tok, remain, recur)
        }
    }
    return recur(root, "", spec)
}

func parseYamlList(n YamlNode, spec string, last string, tok string, remain string,
    recur func(YamlNode, string, string) (YamlNode, error)) (YamlNode, error) {
    s, ok := n.(YamlList)
    if !ok {
        return nil, &YamlNodeTypeMismatch{
            Node:     n,
            Expected: "gokits.YamlList",
            Full:     spec,
            Spec:     last,
            Token:    tok,
        }
    }

    if tok[0] == '[' && tok[len(tok)-1] == ']' {
        if num, err := strconv.Atoi(tok[1 : len(tok)-1]); err == nil {
            if num >= 0 && num < len(s) {
                return recur(s[num], last+tok, remain)
            }
        }
    }
    return nil, &YamlNodeNotFound{
        Full: spec,
        Spec: last + tok,
    }
}

func parseDefault(n YamlNode, spec string, last string, tok string, remain string,
    recur func(YamlNode, string, string) (YamlNode, error)) (YamlNode, error) {
    m, ok := n.(YamlMap)
    if !ok {
        return nil, &YamlNodeTypeMismatch{
            Node:     n,
            Expected: "gokits.YamlMap",
            Full:     spec,
            Spec:     last,
            Token:    tok,
        }
    }

    n, ok = m[tok[1:]]
    if !ok {
        return nil, &YamlNodeNotFound{
            Full: spec,
            Spec: last + tok,
        }
    }
    return recur(n, last+tok, remain)
}

type YamlNodeNotFound struct {
    Full string
    Spec string
}

func (e *YamlNodeNotFound) Error() string {
    return fmt.Sprintf("yaml: %s: %q not found", e.Full, e.Spec)
}

type YamlNodeTypeMismatch struct {
    Full     string
    Spec     string
    Token    string
    Node     YamlNode
    Expected string
}

func (e *YamlNodeTypeMismatch) Error() string {
    return fmt.Sprintf("yaml: %s: type mismatch: %q is %T, want %s (at %q)",
        e.Full, e.Spec, e.Node, e.Expected, e.Token)
}
