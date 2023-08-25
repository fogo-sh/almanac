package content

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/wikilink"

	"github.com/fogo-sh/almanac/pkg/utils"
)

type File struct {
	Path          string
	ParsedContent []byte
}

func ParseFile(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer utils.DeferredClose(f)

	stat, err := f.Stat()
	if err != nil {
		return File{}, fmt.Errorf("failed to stat file: %w", err)
	}

	content := make([]byte, stat.Size())
	_, err = f.Read(content)
	if err != nil {
		return File{}, fmt.Errorf("failed to read file: %w", err)
	}

	var buf bytes.Buffer
	if err := goldmark.New(goldmark.WithExtensions(
		&frontmatter.Extender{},
		&wikilink.Extender{
			Resolver: WikiLinkResolver{},
		},
	)).Convert(content, &buf); err != nil {
		return File{}, fmt.Errorf("failed to parse markdown: %w", err)
	}

	return File{
		Path:          path,
		ParsedContent: buf.Bytes(),
	}, nil
}
