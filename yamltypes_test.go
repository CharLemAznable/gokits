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

var stringTests = []struct {
    Tree   YamlNode
    Expect string
}{
    {
        Tree: YamlScalar("test"),
        Expect: `test
`,
    },
    {
        Tree: YamlList{
            YamlScalar("One"),
            YamlScalar("Two"),
            YamlScalar("Three"),
        },
        Expect: `- One
- Two
- Three
`,
    },
    {
        Tree: YamlMap{
            "phonetic":     YamlScalar("true"),
            "organization": YamlScalar("Navy"),
            "alphabet": YamlList{
                YamlScalar("Alpha"),
                YamlScalar("Bravo"),
                YamlScalar("Charlie"),
            },
        },
        Expect: `organization: Navy
phonetic:     true
alphabet:
  - Alpha
  - Bravo
  - Charlie
`,
    },
    {
        Tree: YamlMap{
            "answer": YamlScalar("42"),
            "question": YamlList{
                YamlScalar("What do you get when you multiply six by nine?"),
                YamlScalar("How many roads must a man walk down?"),
            },
        },
        Expect: `answer: 42
question:
  - What do you get when you multiply six by nine?
  - How many roads must a man walk down?
`,
    },
    {
        Tree: YamlList{
            YamlMap{
                "name": YamlScalar("John Smith"),
                "age":  YamlScalar("42"),
            },
            YamlMap{
                "name": YamlScalar("Jane Smith"),
                "age":  YamlScalar("45"),
            },
        },
        Expect: `- age:  42
  name: John Smith
- age:  45
  name: Jane Smith
`,
    },
    {
        Tree: YamlList{
            YamlList{YamlScalar("one"), YamlScalar("two"), YamlScalar("three")},
            YamlList{YamlScalar("un"), YamlScalar("deux"), YamlScalar("trois")},
            YamlList{YamlScalar("ichi"), YamlScalar("ni"), YamlScalar("san")},
        },
        Expect: `- - one
  - two
  - three
- - un
  - deux
  - trois
- - ichi
  - ni
  - san
`,
    },
    {
        Tree: YamlMap{
            "yahoo":  YamlMap{"url": YamlScalar("http://yahoo.com/"), "company": YamlScalar("Yahoo! Inc.")},
            "google": YamlMap{"url": YamlScalar("http://google.com/"), "company": YamlScalar("Google, Inc.")},
        },
        Expect: `google:
  company: Google, Inc.
  url:     http://google.com/
yahoo:
  company: Yahoo! Inc.
  url:     http://yahoo.com/
`,
    },
}

func TestRender(t *testing.T) {
    for idx, test := range stringTests {
        if got, want := YamlRender(test.Tree), test.Expect; got != want {
            t.Errorf("%d. got:\n%s\n%d. want:\n%s\n", idx, got, idx, want)
        }
    }
}
