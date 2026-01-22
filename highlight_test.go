package obsidian_test

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"

	obsidian "github.com/powerman/goldmark-obsidian"
)

func TestHighlight(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(obsidian.NewHighlight()),
	)
	testutil.DoTestCaseFile(markdown, "testdata/highlight.txt", t, testutil.ParseCliCaseArg()...)
}
