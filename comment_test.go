package obsidian_test

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"

	obsidian "github.com/powerman/goldmark-obsidian"
)

func TestComment(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(obsidian.NewComment()),
	)
	testutil.DoTestCaseFile(markdown, "testdata/comment.txt", t, testutil.ParseCliCaseArg()...)
}
