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

/*
colorize does ANSI colored output and emoji substitution.

Colorize takes a fmt.Printf() template string and looks for styles
in the format:

	{style-name:text}

And inserts the color esacpe sequences around the text.  It then passes
the template off to emoji.Sprintf() from https://github.com/kyokomi/emoji.

It uses the color attribute names, color.Output, and color.NoColor from
https://github.com/fatih/color.

The style names are modeled off the basic bootstrap color names.

Example use:

	colorize.Colorize("{primary:We} want {bg-success:beer :beer:} and {warning:pizza} :pizza:")

*/
package colorize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

const escape = "\x1b"

// Style is a ANSI escape code sequence
type Style struct {
	sequence string
	enable   *bool
}

// NewStyle creates a style from a set of attributes
func NewStyle(attributes ...color.Attribute) *Style {
	codes := make([]string, len(attributes))
	for i, attr := range attributes {
		codes[i] = strconv.Itoa(int(attr))
	}
	sequence := fmt.Sprintf("%s[%sm", escape, strings.Join(codes, ";"))
	return &Style{sequence, nil}
}

// Styles maps style names to styles
var Styles = map[string]*Style{
	"primary":      NewStyle(color.FgCyan, color.Bold),
	"bg-primary":   NewStyle(color.FgHiWhite, color.BgCyan),
	"secondary":    NewStyle(color.FgHiBlack, color.Bold),
	"bg-secondary": NewStyle(color.FgHiWhite, color.BgHiBlack),
	"light":        NewStyle(color.FgHiBlack),
	"info":         NewStyle(color.FgBlue, color.Bold),
	"bg-info":      NewStyle(color.FgWhite, color.BgBlue),
	"success":      NewStyle(color.FgGreen),
	"bg-success":   NewStyle(color.FgHiWhite, color.BgGreen),
	"warning":      NewStyle(color.FgYellow),
	"bg-warning":   NewStyle(color.FgBlack, color.BgYellow),
	"danger":       NewStyle(color.FgHiRed),
	"bg-danger":    NewStyle(color.FgBlack, color.BgHiRed),
}
var reset = NewStyle(color.Reset)

// IsEnabled returns true if the style is enabled, either specifically
// to this style or globally via color.NoColor
func (style *Style) IsEnabled() bool {
	if style == nil {
		return false
	}
	if style.enable != nil {
		return *style.enable
	}
	return !color.NoColor
}

// Enable or disable this style
func (style *Style) Enable(val bool) *Style {
	style.enable = &val
	return style
}

// Colorize some text, printing the output to color.Output
func Colorize(template string, values ...interface{}) {
	fmt.Fprint(color.Output, SColorize(template, values...))
}

// Colorize some text, returning a string
func SColorize(template string, values ...interface{}) string {
	var b strings.Builder
	state := ' '
	var nameStart int
	var style *Style
	var frag strings.Builder
	for i, ch := range template {
		switch state {
		case ' ':
			if ch == '{' {
				nameStart = i
				state = ch
			} else {
				b.WriteRune(ch)
				if ch == '\\' {
					state = ch
				}
			}
		case '\\':
			b.WriteRune(ch)
			state = ' '
		case '{':
			if ch == ':' {
				name := template[nameStart+1 : i]
				style = Styles[name]
				frag.Reset()
				state = '}'
			}
		case '}':
			if ch == '\\' {
				state = ']'
			} else if ch == '}' {
				var enabled = style.IsEnabled()
				if enabled {
					b.WriteString(style.sequence)
				}
				b.WriteString(frag.String())
				if enabled {
					b.WriteString(reset.sequence)
				}
				state = ' '
			} else {
				frag.WriteRune(ch)
			}
		case ']':
			frag.WriteRune(ch)
			state = '}'

		}
	}
	switch state {
	case '{', '}', ']':
		// if we found a '{' w/o completion, append everything from the start
		b.WriteString(template[nameStart:])
	}
	template = b.String()
	return emoji.Sprintf(template, values...)
}
