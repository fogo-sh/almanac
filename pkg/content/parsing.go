package content

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/wikilink"

	"pkg.fogo.sh/almanac/pkg/content/extensions"
	"pkg.fogo.sh/almanac/pkg/utils"
)

type PageMeta struct {
	Categories []string   `toml:"categories"`
	Date       *time.Time `toml:"date"`
	Redirect   *string    `toml:"redirect"`
	Root       bool       `toml:"root"`
	YoutubeId  string     `toml:"youtube_id"`
}

type Page struct {
	Title         string
	Path          *string
	LinksTo       []string
	Backlinks     []string
	Meta          PageMeta
	ParsedContent []byte
}

type Parser struct {
	DiscordUserResolver *extensions.DiscordUserResolver
}

func (p *Parser) ParsePageFile(path string) (Page, error) {
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

	md := goldmark.New(goldmark.WithExtensions(
		&frontmatter.Extender{},
		&wikilink.Extender{
			Resolver: WikiLinkResolver{
				recordDestination: func(destination []byte) error {
					linksTo = append(linksTo, string(destination))
					return nil
				},
			},
		},
		extensions.NewDiscordMention(p.DiscordUserResolver),
	))

	ctx := parser.NewContext()

	var buf bytes.Buffer
	err = md.Convert(content, &buf, parser.WithContext(ctx))
	if err != nil {
		return Page{}, fmt.Errorf("failed to parse markdown: %w", err)
	}

	var pageMeta PageMeta

	data := frontmatter.Get(ctx)

	if data != nil {
		if err := data.Decode(&pageMeta); err != nil {
			return Page{}, fmt.Errorf("failed to decode frontmatter: %w", err)
		}
	}

	pageTitle := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return Page{
		Title:         pageTitle,
		LinksTo:       linksTo,
		Path:          &path,
		Meta:          pageMeta,
		ParsedContent: buf.Bytes(),
	}, nil
}
