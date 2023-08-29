package content

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
)

func CreateSpecialPages(pages map[string]*Page) error {
	pagesByCategory := PagesByCategory(pages)
	allCategories := AllCategories(pages)

	for _, category := range allCategories {
		keysOfPagesInCategory := make([]string, 0, len(pagesByCategory[category]))
		for _, page := range pagesByCategory[category] {
			keysOfPagesInCategory = append(keysOfPagesInCategory, page.Title)
		}

		var buf bytes.Buffer
		err := LinkListingTemplate.Execute(&buf, LinkListingData{
			LinkList: keysOfPagesInCategory,
		})
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		pages[fmt.Sprintf("$Category:%s", category)] = &Page{
			Title:         fmt.Sprintf("$Category:%s", category),
			LinksTo:       keysOfPagesInCategory,
			ParsedContent: buf.Bytes(),
		}
	}

	return nil
}

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

	err := CreateSpecialPages(pages)
	if err != nil {
		return nil, fmt.Errorf("failed to create special pages: %w", err)
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

func PagesByCategory(pages map[string]*Page) map[string][]*Page {
	pagesByCategory := make(map[string][]*Page)

	for _, page := range pages {
		for _, category := range page.Meta.Categories {
			pagesByCategory[category] = append(pagesByCategory[category], page)
		}
	}

	return pagesByCategory
}

func AllCategories(pages map[string]*Page) []string {
	categories := map[string]struct{}{}
	for _, page := range pages {
		for _, category := range page.Meta.Categories {
			categories[category] = struct{}{}
		}
	}

	keys := make([]string, 0)
	for k := range categories {
		keys = append(keys, k)
	}

	return keys
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
