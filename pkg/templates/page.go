package templates

import (
	"html/template"

	"github.com/fogo-sh/almanac/pkg/content"
)

type PageTemplateData struct {
	AllPageTitles []string
	Page          *content.Page
	Content       template.HTML
}

var pageTemplateContent = `<!DOCTYPE html>
<html>
	<head>
		<title>{{ .Page.Title }}</title>
		<link rel="stylesheet" href="/assets/css/main.css">
		<link rel="icon" type="image/svg+xml" href="/favicon.svg">
	</head>
	<body>
		<nav>
			<ul>
			{{ range .AllPageTitles }}
				<li><a href="/{{ . }}">{{ . }}</a></li>
			{{ end }}
			</ul>
		</nav>
		<main>
			<h1>{{ .Page.Title }}</h1>

			{{ if .Page.Meta.YoutubeId }}
			<iframe
			  width="100%"
			  height="600px"
			  src="https://www.youtube.com/embed/{{ .Page.Meta.YoutubeId }}"
			  frameborder="0"
			  allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
			  allowfullscreen></iframe>
			{{ end }}

			{{ .Content }}

			{{ if .Page.Backlinks }}
			<section>
				<h2>Backlinks</h2>
				<ul>
				{{ range .Page.Backlinks }}
					<li><a href="/{{ . }}">{{ . }}</a></li>
				{{ end }}
				</ul>
			</section>
			{{ end }}
		</main>
	</body>
</html>`

var PageTemplate *template.Template

type PageData struct {
	Title   string
	Content template.HTML
}

func init() {
	t, err := template.New("page").Parse(pageTemplateContent)
	if err != nil {
		panic(err)
	}
	PageTemplate = t
}
