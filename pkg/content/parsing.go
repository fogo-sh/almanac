package content

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/wikilink"

	"github.com/fogo-sh/almanac/pkg/utils"
)

type Page struct {
	Title         string
	Path          string
	LinksTo       []string
	Backlinks     []string
	ParsedContent []byte
}

func ParsePageFile(path string) (Page, error) {
	f, err := os.Open(path)
	if err != nil {
		return Page{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer utils.DeferredClose(f)

	stat, err := f.Stat()
	if err != nil {
		return Page{}, fmt.Errorf("failed to stat file: %w", err)
	}

	content := make([]byte, stat.Size())
	_, err = f.Read(content)
	if err != nil {
		return Page{}, fmt.Errorf("failed to read file: %w", err)
	}

	var linksTo = make([]string, 0)

	var buf bytes.Buffer
	if err := goldmark.New(goldmark.WithExtensions(
		&frontmatter.Extender{},
		&wikilink.Extender{
			Resolver: WikiLinkResolver{
				recordDestination: func(destination []byte) error {
					linksTo = append(linksTo, string(destination))
					return nil
				},
			},
		},
	)).Convert(content, &buf); err != nil {
		return Page{}, fmt.Errorf("failed to parse markdown: %w", err)
	}

	pageTitle := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return Page{
		Title:         pageTitle,
		LinksTo:       linksTo,
		Path:          path,
		ParsedContent: buf.Bytes(),
	}, nil
}
