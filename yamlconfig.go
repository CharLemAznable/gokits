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

    // TODO(kevlar): Add a cache?
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
    defer fin.Close()

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
func ConfigYamlFile(filename string) *YamlFile {
    f, err := ReadYamlFile(filename)
    if err != nil {
        panic(err)
    }
    return f
}

func (f *YamlFile) Get(spec string) (YamlNode, error) {
    node, err := YamlChild(f.Root, spec)
    if err != nil {
        return nil, err
    }

    if node == nil {
        return nil, &YamlNodeNotFound{
            Full: spec,
            Spec: spec,
        }
    }

    return node, nil
}

func ScalarOfYaml(node YamlNode, spec string) (YamlScalar, error) {
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

func StringOfYaml(node YamlNode, spec string) (string, error) {
    sc, err := ScalarOfYaml(node, spec)
    if err != nil {
        return "", err
    }

    return sc.String()
}

func IntOfYaml(node YamlNode, spec string) (int64, error) {
    sc, err := ScalarOfYaml(node, spec)
    if err != nil {
        return 0, err
    }

    return sc.Int()
}

func BoolOfYaml(node YamlNode, spec string) (bool, error) {
    sc, err := ScalarOfYaml(node, spec)
    if err != nil {
        return false, err
    }

    return sc.Bool()
}

func ListOfYaml(node YamlNode, spec string) (YamlList, error) {
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

func MapOfYaml(node YamlNode, spec string) (YamlMap, error) {
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

// GetString retrieves a scalar from the file specified by a string of the same
// format as that expected by YamlChild.  If the final node is not a YamlScalar, GetString
// will return an error.
func (f *YamlFile) GetScalar(spec string) (YamlScalar, error) {
    node, err := f.Get(spec)
    if err != nil {
        return "", err
    }

    return ScalarOfYaml(node, spec)
}

func (f *YamlFile) GetString(spec string) (string, error) {
    scalar, err := f.GetScalar(spec)
    if err != nil {
        return "", err
    }

    return scalar.String()
}

func (f *YamlFile) GetInt(spec string) (int64, error) {
    scalar, err := f.GetScalar(spec)
    if err != nil {
        return 0, err
    }

    return scalar.Int()
}

func (f *YamlFile) GetBool(spec string) (bool, error) {
    scalar, err := f.GetScalar(spec)
    if err != nil {
        return false, err
    }

    return scalar.Bool()
}

func (f *YamlFile) GetList(spec string) (YamlList, error) {
    node, err := f.Get(spec)
    if err != nil {
        return nil, err
    }

    return ListOfYaml(node, spec)
}

func (f *YamlFile) GetMap(spec string) (YamlMap, error) {
    node, err := f.Get(spec)
    if err != nil {
        return nil, err
    }

    return MapOfYaml(node, spec)
}

// Count retrieves a the number of elements in the specified list from the file
// using the same format as that expected by YamlChild.  If the final node is not a
// YamlList, Count will return an error.
func (f *YamlFile) Count(spec string) (int, error) {
    node, err := YamlChild(f.Root, spec)
    if err != nil {
        return -1, err
    }

    if node == nil {
        return -1, &YamlNodeNotFound{
            Full: spec,
            Spec: spec,
        }
    }

    lst, ok := node.(YamlList)
    if !ok {
        return -1, &YamlNodeTypeMismatch{
            Full:     spec,
            Spec:     spec,
            Token:    "$",
            Expected: "gokits.YamlList",
            Node:     node,
        }
    }
    return lst.Len(), nil
}

// Require retrieves a scalar from the file specified by a string of the same
// format as that expected by YamlChild.  If the final node is not a YamlScalar, String
// will panic.  This is a convenience function for use in initializers.
func (f *YamlFile) Require(spec string) string {
    str, err := f.GetString(spec)
    if err != nil {
        panic(err)
    }
    return str
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
        default:
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
    }
    return recur(root, "", spec)
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
