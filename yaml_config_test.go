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
    "testing"
)

var dummyConfigFile = `
mapping:
  key1: value1
  key2: value2
  key3: 5
  key4: true
  key5: false
list:
  - item1
  - item2
config:
  server:
    - www.google.com
    - www.cnn.com
    - www.example.com
  admin:
    - username: god
      password: z3u5
    - username: lowly
      password: f!r3m3
`

var configGetTests = []struct {
    Spec string
    Want string
    Err  string
}{
    {"mapping.key1", "value1", ""},
    {"mapping.key2", "value2", ""},
    {"list[0]", "item1", ""},
    {"list[1]", "item2", ""},
    {"list", "", `yaml: list: type mismatch: "list" is gokits.YamlList, want gokits.YamlScalar (at "$")`},
    {"list.0", "", `yaml: .list.0: type mismatch: ".list" is gokits.YamlList, want gokits.YamlMap (at ".0")`},
    {"config.server[0]", "www.google.com", ""},
    {"config.server[1]", "www.cnn.com", ""},
    {"config.server[2]", "www.example.com", ""},
    {"config.server[3]", "", `yaml: .config.server[3]: ".config.server[3]" not found`},
    {"config.listen[0]", "", `yaml: .config.listen[0]: ".config.listen" not found`},
    {"config.admin[0].username", "god", ""},
    {"config.admin[1].username", "lowly", ""},
    {"config.admin[2].username", "", `yaml: .config.admin[2].username: ".config.admin[2]" not found`},
}

func TestYamlGet(t *testing.T) {
    config := ConfigYamlString(dummyConfigFile)

    for _, test := range configGetTests {
        got, err := config.GetString(test.Spec)
        if want := test.Want; got != want {
            t.Errorf("GetString(%q) = %q, want %q", test.Spec, got, want)
        }

        switch err {
        case nil:
            got = ""
        default:
            got = err.Error()
        }
        if want := test.Err; got != want {
            t.Errorf("GetString(%q) error %#q, want %#q", test.Spec, got, want)
        }
    }

    i, err := config.GetInt("mapping.key3")
    if err != nil || i != 5 {
        t.Errorf("GetInt mapping.key3 wrong")
    }

    b, err := config.GetBool("mapping.key4")
    if err != nil || b != true {
        t.Errorf("GetBool mapping.key4 wrong")
    }

    b, err = config.GetBool("mapping.key5")
    if err != nil || b != false {
        t.Errorf("GetBool mapping.key5 wrong")
    }
}

var configGetOfTests = []struct {
    Spec string
    Want string
    Err  string
}{
    {"mapping.key1", "value1", ""},
    {"mapping.key2", "value2", ""},
    {"list[0]", "item1", ""},
    {"list[1]", "item2", ""},
    {"list", "", `yaml: list: type mismatch: "list" is gokits.YamlList, want gokits.YamlScalar (at "$")`},
    {"list.0", "", `yaml: list.0: type mismatch: "list.0" is <nil>, want gokits.YamlScalar (at "$")`},
    {"config.server[0]", "www.google.com", ""},
    {"config.server[1]", "www.cnn.com", ""},
    {"config.server[2]", "www.example.com", ""},
    {"config.server[3]", "", `yaml: config.server[3]: type mismatch: "config.server[3]" is <nil>, want gokits.YamlScalar (at "$")`},
    {"config.listen[0]", "", `yaml: config.listen[0]: type mismatch: "config.listen[0]" is <nil>, want gokits.YamlScalar (at "$")`},
    {"config.admin[0].username", "god", ""},
    {"config.admin[1].username", "lowly", ""},
    {"config.admin[2].username", "", `yaml: config.admin[2].username: type mismatch: "config.admin[2].username" is <nil>, want gokits.YamlScalar (at "$")`},
}

func TestYamlGetOf(t *testing.T) {
    config := ConfigYamlString(dummyConfigFile)

    for _, test := range configGetOfTests {
        node, _ := config.Get(test.Spec)
        got, err := StringOfYaml(node, test.Spec)
        if want := test.Want; got != want {
            t.Errorf("GetString(%q) = %q, want %q", test.Spec, got, want)
        }

        switch err {
        case nil:
            got = ""
        default:
            got = err.Error()
        }
        if want := test.Err; got != want {
            t.Errorf("GetString(%q) error %#q, want %#q", test.Spec, got, want)
        }
    }

    node3, _ := config.Get("mapping.key3")
    i, err := IntOfYaml(node3, "mapping.key3")
    if err != nil || i != 5 {
        t.Errorf("GetInt mapping.key3 wrong")
    }

    node4, _ := config.Get("mapping.key4")
    b, err := BoolOfYaml(node4, "mapping.key4")
    if err != nil || b != true {
        t.Errorf("GetBool mapping.key4 wrong")
    }

    node5, _ := config.Get("mapping.key5")
    b, err = BoolOfYaml(node5, "mapping.key5")
    if err != nil || b != false {
        t.Errorf("GetBool mapping.key5 wrong")
    }
}

func TestYamlCountAndRequire(t *testing.T) {
    config := ConfigYamlString(dummyConfigFile)

    count, err := config.Count("list")
    if nil != err && count != 2 {
        t.Errorf("Count list wrong")
    }

    countErr, err := config.Count("config")
    if nil == err || countErr != -1 {
        t.Errorf("Count config wrong")
    }
}

func TestYamlTypes(t *testing.T) {
    config := ConfigYamlString(dummyConfigFile)

    mapping, _ := config.GetMap("mapping")
    value1, _ := StringOfYaml(mapping.Key("key1"), "mapping.key1")
    if "value1" != value1 {
        t.Errorf("YamlMap.Key() wrong")
    }

    list, _ := config.GetList("list")
    if 2 != list.Len() {
        t.Errorf("YamlList.Len() wrong")
    }
    item1, _ := StringOfYaml(list.Item(0), "list[0]")
    if "item1" != item1 {
        t.Errorf("YamlList.Item() wrong")
    }
}