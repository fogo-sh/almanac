package templates

import (
	"bufio"
	"html/template"
	"os"

	"github.com/fogo-sh/almanac/pkg/content"

	cp "github.com/otiai10/copy"
)

func OutputAllPagesToDisk(pages map[string]*content.Page, outputDir string) error {
	os.RemoveAll(outputDir)

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	allPageTitles := content.AllPageTitles(pages)

	for key, page := range pages {
		var outputPath = outputDir + key + ".html"

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		defer f.Close()

		w := bufio.NewWriter(f)

		err = PageTemplate.Execute(w, PageTemplateData{
			AllPageTitles: allPageTitles,
			Title:         page.Title,
			Content:       template.HTML(string(page.ParsedContent)),
			Backlinks:     page.Backlinks,
		})
		if err != nil {
			return err
		}

		err = w.Flush()
		if err != nil {
			return err
		}
	}

	err = cp.Copy("pkg/static/static/.", outputDir)
	if err != nil {
		return err
	}

	return nil
}
