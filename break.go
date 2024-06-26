// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import (
	"bytes"
)

type ThematicBreak struct {
	Position
}

func (b *ThematicBreak) PrintHTML(buf *bytes.Buffer) {
	buf.WriteString("<hr />\n")
}

func (b *ThematicBreak) printMarkdown(buf *markOut, s mdState) {
	buf.maybeNL()
	buf.WriteString("***")
}

func newHR(p *parseState, s line) (line, bool) {
	if isHR(s) {
		p.doneBlock(&ThematicBreak{Position{p.lineno, p.lineno}})
		return line{}, true
	}
	return s, false
}

func isHR(s line) bool {
	t := s
	t.trimSpace(0, 3, false)
	switch c := t.peek(); c {
	case '-', '_', '*':
		for i := 0; ; i++ {
			if !t.trim(c) {
				if i >= 3 {
					break
				}
				return false
			}
			t.skipSpace()
		}
		return t.eof()
	}
	return false
}

type HardBreak struct{}

func (*HardBreak) Inline() {}

func (x *HardBreak) PrintHTML(buf *bytes.Buffer) {
	buf.WriteString("<br />\n")
}

func (x *HardBreak) printMarkdown(buf *markOut) {
	buf.WriteString("\\")
	buf.NL()
}

func (x *HardBreak) PrintText(buf *bytes.Buffer) {
	buf.WriteString("\n")
}

type SoftBreak struct{}

func (*SoftBreak) Inline() {}

func (x *SoftBreak) PrintHTML(buf *bytes.Buffer) {
	buf.WriteString("\n")
}

func (x *SoftBreak) printMarkdown(buf *markOut) {
	buf.NL()
}

func (x *SoftBreak) PrintText(buf *bytes.Buffer) {
	buf.WriteString("\n")
}

func parseBreak(_ *parseState, s string, i int) (Inline, int, int, bool) {
	start := i
	for start > 0 && (s[start-1] == ' ' || s[start-1] == '\t') {
		start--
	}
	end := i + 1
	// TODO: Do tabs count? That would be a mess.
	if i >= 2 && s[i-1] == ' ' && s[i-2] == ' ' {
		return &HardBreak{}, start, end, true
	}
	return &SoftBreak{}, start, end, true
}
