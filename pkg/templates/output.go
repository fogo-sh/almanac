package templates

import (
	"bufio"
	"html/template"
	"os"
	"os/exec"

	"github.com/fogo-sh/almanac/pkg/content"
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

	err = exec.Command("cp", "-r", "pkg/static/static/.", outputDir).Run()
	if err != nil {
		return err
	}

	return nil
}
