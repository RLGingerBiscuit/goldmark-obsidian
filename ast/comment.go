package ast

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var KindComment = gast.NewNodeKind("Comment")

type Comment struct {
	gast.BaseInline
	Segments *text.Segments
}

func (n *Comment) Kind() gast.NodeKind {
	return KindComment
}

func (n *Comment) Dump(source []byte, level int) {
	m := map[string]string{}
	t := []string{}
	for i := 0; i < n.Segments.Len(); i++ {
		segment := n.Segments.At(i)
		t = append(t, string(segment.Value(source)))
	}
	m["Segments"] = strings.Join(t, "")
	gast.DumpHelper(n, source, level, m, nil)
}

func NewComment() *Comment {
	return &Comment{
		Segments: text.NewSegments(),
	}
}
