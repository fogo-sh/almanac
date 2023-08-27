package content

import (
	"fmt"
	"path/filepath"
	"sort"
)

func DiscoverPages(path string) (map[string]Page, error) {
	paths, error := filepath.Glob(filepath.Join(path, "*.md"))

	if error != nil {
		return nil, fmt.Errorf("failed to glob files: %w", error)
	}

	pages := make(map[string]Page)

	for _, path := range paths {
		page, error := ParsePageFile(path)
		if error != nil {
			return nil, fmt.Errorf("failed to parse page: %w", error)
		}

		pages[page.Title] = page
	}

	return pages, nil
}

func AllPageTitles(pages map[string]Page) []string {
	allPageTitles := make([]string, 0, len(pages))
	for key := range pages {
		allPageTitles = append(allPageTitles, key)
	}

	sort.Slice(allPageTitles, func(i, j int) bool {
		return allPageTitles[i] < allPageTitles[j]
	})

	return allPageTitles
}
