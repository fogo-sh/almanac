package content

import (
	"html/template"
)

type PageTemplateData struct {
	AllPageTitles []string
	Page          *Page
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

			{{ if .Page.Meta.Redirect }}
			<p>↳ <a href="/{{ .Page.Meta.Redirect }}">{{ .Page.Meta.Redirect }}</a></p>
			{{ end }}

			{{ if .Page.Meta.Categories }}
			<p>
			{{ range .Page.Meta.Categories }}
				<a href="/$Category:{{ . }}">{{ . }}</a>
			{{ end }}
			</p>
			{{ end }}

			{{ if .Page.Meta.Date }}
			<p>{{ .Page.Meta.Date.Format "Aug 2, 2006" }}</p>
			{{ end }}

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

var linkListingTemplateContent = `<ul>
{{ range .LinkList }}
	<li><a href="/{{ . }}">{{ . }}</a></li>
{{ end }}
</ul>`

var LinkListingTemplate *template.Template

type LinkListingData struct {
	LinkList []string
}

func initTemplate(name string, content string) *template.Template {
	t, err := template.New(name).Parse(content)
	if err != nil {
		panic(err)
	}
	return t
}

func init() {
	PageTemplate = initTemplate("page", pageTemplateContent)
	LinkListingTemplate = initTemplate("linkListing", linkListingTemplateContent)
}
