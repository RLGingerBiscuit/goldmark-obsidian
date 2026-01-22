package obsidian

import (
	"github.com/powerman/goldmark-obsidian/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type highlightDelimiterProcessor struct {
}

func newHighlightDelimiterProcessor() *highlightDelimiterProcessor {
	return &highlightDelimiterProcessor{}
}

func (p *highlightDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '='
}

func (p *highlightDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *highlightDelimiterProcessor) OnMatch(consumes int) gast.Node {
	return ast.NewHighlight()
}

// HighlightParser is an Obsidian highlight parser.
type HighlightParser struct {
}

// NewHighlightParser return a new HighlightParser.
func NewHighlightParser() HighlightParser {
	return HighlightParser{}
}

// Trigger implements [parser.InlineParser].
func (HighlightParser) Trigger() []byte {
	return []byte{'='}
}

// Parse implements [parser.InlineParser].
func (HighlightParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := parser.ScanDelimiter(line, before, 1, newHighlightDelimiterProcessor())
	if node == nil || node.OriginalLength > 2 || before == '=' {
		return nil
	}

	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
	return node
}

// HighlightHTMLRenderer is a HTML renderer for highlights.
type HighlightHTMLRenderer struct {
	html.Config
}

// NewHighlightHTMLRenderer returns a new HighlightHTMLRenderer.
func NewHighlightHTMLRenderer(opts ...html.Option) *HighlightHTMLRenderer {
	r := &HighlightHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements [renderer.NodeRenderer].
func (r *HighlightHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHighlight, r.renderHighlight)
}

// HighlightAttributeFilter defines attribute names which mark elements can have.
var HighlightAttributeFilter = html.GlobalAttributeFilter

func (r *HighlightHTMLRenderer) renderHighlight(
	w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<mark")
			html.RenderAttributes(w, n, HighlightAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<mark>")
		}
	} else {
		_, _ = w.WriteString("</mark>")
	}
	return gast.WalkContinue, nil
}

// Highlight is an extension that helps setup Obsidian [highlight] parser and HTML renderer.
//
// [highlight]: https://help.obsidian.md/syntax#Bold,%20italics,%20highlights
type Highlight struct {
}

// NewHighlight returns a new Highlight extension.
func NewHighlight() Highlight {
	return Highlight{}
}

// Extend implements [goldmark.Extender].
func (Highlight) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewHighlightParser(), prioInlineParserLowest),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHighlightHTMLRenderer(), prioHTMLRenderer),
	))
}
