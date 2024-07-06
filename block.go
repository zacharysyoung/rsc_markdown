// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

// Block is implemented by:
//
//	CodeBlock
//	Document
//	Empty
//	HTMLBlock
//	Heading
//	Item
//	List
//	Paragraph
//	Quote
//	Text
//	ThematicBreak
type Block interface {
	Block()
	Pos() Position
	printHTML(p *printer)
	printMarkdown(p *printer)
}

type Position struct {
	StartLine int
	EndLine   int
}

func (p Position) Pos() Position {
	return p
}