package obsidian

import (
	"bytes"

	"github.com/powerman/goldmark-obsidian/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var obsidianCommentDelimiter = []byte("%%")

// CommentParser is an Obsidian comment parser.
type CommentParser struct{}

// NewCommentParser returns a new CommentParser.
func NewCommentParser() CommentParser {
	return CommentParser{}
}

func (CommentParser) Trigger() []byte {
	return []byte{'%'}
}

func (CommentParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	savedLineNum, savedSegment := block.Position()
	node := ast.NewComment()
	line, segment := block.PeekLine()

	offset := len(obsidianCommentDelimiter)
	line = line[offset:]
	for {
		index := bytes.Index(line, obsidianCommentDelimiter)
		if index > -1 {
			node.Segments.Append(segment.WithStop(segment.Start + offset + index + len(obsidianCommentDelimiter)))
			block.Advance(offset + index + len(obsidianCommentDelimiter))
			return node
		}
		offset = 0
		node.Segments.Append(segment)
		block.AdvanceLine()
		line, segment = block.PeekLine()
		if line == nil {
			break
		}
	}

	block.SetPosition(savedLineNum, savedSegment)
	return nil
}

// CommentHTMLRenderer is a HTML renderer for an Obsidian comment.
//
// Current implementation does not render comments at all (like Obsidian's "Reading" mode).
type CommentHTMLRenderer struct{}

// NewCommentHTMLRenderer returns a new CommentHTMLRenderer.
func NewCommentHTMLRenderer() *CommentHTMLRenderer {
	return &CommentHTMLRenderer{}
}

// RegisterFuncs implements [renderer.NodeRenderer].
func (r *CommentHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindComment, r.renderComment)
}

func (r *CommentHTMLRenderer) renderComment(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	return gast.WalkSkipChildren, nil
}

// Comment is an extension that helps setup Obsidian [comment] parser and HTML renderer.
//
// [comment]: https://help.obsidian.md/syntax#Comments
type Comment struct{}

// NewComment returns a new Comment extension.
func NewComment() Comment {
	return Comment{}
}

// Extend implements [goldmark.Extender].
func (Comment) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewCommentBlockParser(), prioHighest), // should have higher priority than inline parser
		),
		parser.WithInlineParsers(
			util.Prioritized(NewCommentParser(), prioInlineParserLowest),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewCommentBlockHTMLRenderer(), prioHTMLRenderer),
		),
		renderer.WithNodeRenderers(
			util.Prioritized(NewCommentHTMLRenderer(), prioHTMLRenderer),
		),
	)
}
