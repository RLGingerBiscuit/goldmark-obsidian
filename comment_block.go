package obsidian

import (
	"bytes"

	"github.com/powerman/goldmark-obsidian/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type CommentBlockParser struct{}

func NewCommentBlockParser() CommentBlockParser {
	return CommentBlockParser{}
}

func (CommentBlockParser) Trigger() []byte {
	return []byte{'%'}
}

func (CommentBlockParser) Open(parent gast.Node, block text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, segment := block.PeekLine()

	if bytes.HasPrefix(line, []byte{' ', ' ', ' ', ' '}) || !bytes.Contains(line, obsidianCommentDelimiter) {
		return nil, parser.NoChildren
	}

	node := ast.NewCommentBlock()
	block.AdvanceToEOL()
	node.Lines().Append(segment)
	return node, parser.NoChildren
}

func (CommentBlockParser) Continue(node gast.Node, block text.Reader, pc parser.Context) parser.State {
	lines := node.Lines()
	line, segment := block.PeekLine()

	if lines.Len() == 1 {
		firstLine := lines.At(0)
		if bytes.Count(firstLine.Value(block.Source()), obsidianCommentDelimiter) >= 2 {
			return parser.Close
		}
	}

	if index := bytes.Index(line, obsidianCommentDelimiter); index > -1 {
		segment = segment.WithStop(segment.Start + index + len(obsidianCommentDelimiter))
		node.Lines().Append(segment)
		block.Advance(index + len(obsidianCommentDelimiter))
		return parser.Close
	}

	node.Lines().Append(segment)
	block.AdvanceToEOL()
	return parser.Continue | parser.NoChildren
}

func (CommentBlockParser) Close(node gast.Node, block text.Reader, pc parser.Context) {
	// nothing to do
}

func (CommentBlockParser) CanInterruptParagraph() bool {
	return true
}

func (CommentBlockParser) CanAcceptIndentedLine() bool {
	return false
}

// CommentBlockHTMLRenderer is a HTML renderer for an Obsidian comment block.
//
// Current implementation does not render comment blocks at all (like Obsidian's "Reading" mode).
type CommentBlockHTMLRenderer struct{}

// NewCommentBlockHTMLRenderer returns a new CommentBlockHTMLRenderer.
func NewCommentBlockHTMLRenderer() *CommentBlockHTMLRenderer {
	return &CommentBlockHTMLRenderer{}
}

// RegisterFuncs implements [renderer.NodeRenderer].
func (r *CommentBlockHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindCommentBlock, r.renderCommentBlock)
}

func (r *CommentBlockHTMLRenderer) renderCommentBlock(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	return gast.WalkSkipChildren, nil
}
