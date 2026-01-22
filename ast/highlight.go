package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// KindHighlight is a NodeKind of the Highlight node.
var KindHighlight = gast.NewNodeKind("Highlight") // Const.

// A Highlight struct represents a highlight of Obsidian text.
type Highlight struct {
	gast.BaseInline
}

// NewHighlight returns a new Highlight node.
func NewHighlight() *Highlight {
	return &Highlight{}
}

// Dump implements [ast.Node].
func (n *Highlight) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// Kind implements [ast.Node].
func (n *Highlight) Kind() gast.NodeKind {
	return KindHighlight
}
