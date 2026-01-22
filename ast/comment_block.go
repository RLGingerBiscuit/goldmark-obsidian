package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

var KindCommentBlock = gast.NewNodeKind("CommentBlock")

type CommentBlock struct {
	gast.BaseBlock
}

// NewCommentBlock returns a new CommentBlock node.
func NewCommentBlock() *CommentBlock {
	return &CommentBlock{}
}

// Dump implements [ast.Node].
func (n *CommentBlock) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// Kind implements [ast.Node].
func (n *CommentBlock) Kind() gast.NodeKind {
	return KindCommentBlock
}
