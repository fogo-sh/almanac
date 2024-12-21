package content

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func (r *Renderer) CreateSpecialPages(pages map[string]*Page) error {
	specialPages := make([]string, 0)

	pagesByCategory := r.PagesByCategory(pages)
	allCategories := r.AllCategories(pages)

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

		pageTitle := fmt.Sprintf("$Category:%s", category)
		page := &Page{
			Title:         pageTitle,
			LinksTo:       keysOfPagesInCategory,
			ParsedContent: buf.Bytes(),
		}

		pages[pageTitle] = page
		specialPages = append(specialPages, pageTitle)
	}

	var buf bytes.Buffer
	err := LinkListingTemplate.Execute(&buf, LinkListingData{
		LinkList: specialPages,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	page := &Page{
		Title:         "$SpecialPages",
		LinksTo:       specialPages,
		ParsedContent: buf.Bytes(),
	}

	pages["$SpecialPages"] = page

	return nil
}

func (r *Renderer) DiscoverPages(path string) (map[string]*Page, error) {
	paths, error := filepath.Glob(filepath.Join(path, "*.md"))

	if error != nil {
		return nil, fmt.Errorf("failed to glob files: %w", error)
	}

	pages := make(map[string]*Page)

	for _, path := range paths {
		page, error := r.ParsePageFile(path)
		if error != nil {
			return nil, fmt.Errorf("failed to parse page: %w", error)
		}

		pages[page.Title] = &page
	}

	err := r.CreateSpecialPages(pages)
	if err != nil {
		return nil, fmt.Errorf("failed to create special pages: %w", err)
	}

	r.PopulateBacklinks(pages)

	return pages, nil
}

func (r *Renderer) PopulateBacklinks(pages map[string]*Page) {
	for _, page := range pages {
		for _, link := range page.LinksTo {
			if _, ok := pages[link]; ok {
				found := false

				for _, backlink := range pages[link].Backlinks {
					if backlink == page.Title {
						found = true
						break
					}
				}

				if !found {
					pages[link].Backlinks = append(pages[link].Backlinks, page.Title)
				}
			}
		}
	}

	for _, page := range pages {
		sort.Slice(page.Backlinks, func(i, j int) bool {
			return strings.ToLower(page.Backlinks[i]) < strings.ToLower(page.Backlinks[j])
		})
	}
}

func (r *Renderer) AllPageTitles(pages map[string]*Page) []string {
	allPageTitles := make([]string, 0, len(pages))
	for key := range pages {
		allPageTitles = append(allPageTitles, key)
	}

	sort.Slice(allPageTitles, func(i, j int) bool {
		return strings.ToLower(allPageTitles[i]) < strings.ToLower(allPageTitles[j])
	})

	return allPageTitles
}

func (r *Renderer) PagesByCategory(pages map[string]*Page) map[string][]*Page {
	pagesByCategory := make(map[string][]*Page)

	for _, page := range pages {
		for _, category := range page.Meta.Categories {
			pagesByCategory[category] = append(pagesByCategory[category], page)
		}
	}

	return pagesByCategory
}

func (r *Renderer) AllCategories(pages map[string]*Page) []string {
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

func (r *Renderer) FindRootPage(pages map[string]*Page) (*Page, error) {
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
