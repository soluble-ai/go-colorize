// Copyright 2019 Soluble Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colorize

import (
	"fmt"
	"github.com/fatih/color"
	"sort"
	"strings"
	"testing"
)

func TestDemoDefaults(t *testing.T) {
	t.SkipNow()
	var template strings.Builder
	params := make([]interface{}, 0, len(Styles))
	names := make([]string, 0, len(Styles))
	for k := range Styles {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		template.WriteString(fmt.Sprintf("{%s: %%12s }\n", k))
		params = append(params, k)
	}
	Colorize(template.String(), params...)
	Colorize("\n{primary:%s} eat {bg-success:pizza :pizza:}{bg-primary: and drink }{warning:beer} %s, but not {bg-danger:%s}\n\n",
		"Most folks like to", ":beer:", "everyone")
}

func TestBasic(t *testing.T) {
	s := SColorize("{primary:%s} %s", "hello", "world")
	if !strings.Contains(s, "hello") || !strings.Contains(s, "world") {
		t.Errorf("string contents wrong")
	}
	Styles["primary"].Enable(false)
	if Styles["primary"].IsEnabled() {
		t.Error("primary is enabled")
	}
	if Styles["nope"].IsEnabled() {
		t.Error("nope is enabled")
	}
	if SColorize("{primary:hello} world") != "hello world" {
		t.Error("disable style didn't work")
	}
	color.Output = &strings.Builder{}
	Colorize("hello, world")
	if color.Output.(*strings.Builder).String() != "hello, world" {
		t.Error("string didn't output")
	}
}

func TestEscapes(t *testing.T) {
	if e := SColorize("\\\n\r\t\v"); e != "\\\n\r\t\v" {
		t.Errorf("escape sequences wrong: %s", e)
	}
}

func TestIncomplete(t *testing.T) {
	if s := SColorize("{primary:foo"); s != "{primary:foo" {
		t.Error(s)
	}
	if s := SColorize("{nope:xxx\\}}"); s != "xxx}" {
		t.Error(s)
	}
	if s := SColorize("{nope:foo\\ bar"); s != "{nope:foo\\ bar" {
		t.Error(s)
	}
}
