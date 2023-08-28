package content

import (
	"fmt"
	"path/filepath"
	"sort"
)

func DiscoverPages(path string) (map[string]*Page, error) {
	paths, error := filepath.Glob(filepath.Join(path, "*.md"))

	if error != nil {
		return nil, fmt.Errorf("failed to glob files: %w", error)
	}

	pages := make(map[string]*Page)

	for _, path := range paths {
		page, error := ParsePageFile(path)
		if error != nil {
			return nil, fmt.Errorf("failed to parse page: %w", error)
		}

		pages[page.Title] = &page
	}

	PopulateBacklinks(pages)

	return pages, nil
}

func PopulateBacklinks(pages map[string]*Page) {
	for _, page := range pages {
		for _, link := range page.LinksTo {
			if _, ok := pages[link]; ok {
				pages[link].Backlinks = append(pages[link].Backlinks, page.Title)
			}
		}
	}
}

func AllPageTitles(pages map[string]*Page) []string {
	allPageTitles := make([]string, 0, len(pages))
	for key := range pages {
		allPageTitles = append(allPageTitles, key)
	}

	sort.Slice(allPageTitles, func(i, j int) bool {
		return allPageTitles[i] < allPageTitles[j]
	})

	return allPageTitles
}

func FindRootPage(pages map[string]*Page) (*Page, error) {
	var rootPage *Page

	for _, page := range pages {
		if page.Meta.Root {
			if rootPage != nil {
				return &Page{}, fmt.Errorf("multiple root pages found")
			}

			rootPage = page
		}
	}

	if rootPage == nil {
		return &Page{}, fmt.Errorf("no root page found")
	}

	return rootPage, nil
}
