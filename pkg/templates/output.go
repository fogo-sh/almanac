package templates

import (
	"bufio"
	"fmt"
	"html/template"
	"os"

	"github.com/fogo-sh/almanac/pkg/content"

	cp "github.com/otiai10/copy"
)

func OutputAllPagesToDisk(pages map[string]*content.Page, outputDir string) error {
	os.RemoveAll(outputDir)

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	allPageTitles := content.AllPageTitles(pages)

	for key, page := range pages {
		outputPath := outputDir + key + ".html"

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}

		defer f.Close()

		w := bufio.NewWriter(f)

		err = PageTemplate.Execute(w, PageTemplateData{
			AllPageTitles: allPageTitles,
			Page:          page,
			Content:       template.HTML(string(page.ParsedContent)),
		})
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		err = w.Flush()
		if err != nil {
			return fmt.Errorf("failed to flush template: %w", err)
		}
	}

	err = cp.Copy("pkg/static/static/.", outputDir)
	if err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	return nil
}
