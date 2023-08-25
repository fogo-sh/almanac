package content

import (
	"bytes"

	"go.abhg.dev/goldmark/wikilink"
)

type WikiLinkResolver struct{}

func (r WikiLinkResolver) ResolveWikilink(node *wikilink.Node) (destination []byte, err error) {
	destination, err = wikilink.DefaultResolver.ResolveWikilink(node)
	destination = bytes.Replace(destination, []byte(".html"), []byte(""), -1)
	return destination, err
}

var _ wikilink.Resolver = WikiLinkResolver{}
