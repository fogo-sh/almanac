package content

import (
	"bytes"

	"go.abhg.dev/goldmark/wikilink"
)

type WikiLinkResolver struct {
	recordDestination func([]byte) error
}

func (r WikiLinkResolver) ResolveWikilink(node *wikilink.Node) (destination []byte, err error) {
	destination, err = wikilink.DefaultResolver.ResolveWikilink(node)

	if err != nil {
		return destination, err
	}

	destination = bytes.Replace(destination, []byte(".html"), []byte(""), -1)

	if r.recordDestination != nil {
		err = r.recordDestination(destination)
	}

	return destination, err
}

var _ wikilink.Resolver = WikiLinkResolver{}
