package templates

import "html/template"

type PageTemplateData struct {
	AllPageTitles []string
	Title         string
	Content       template.HTML
	Backlinks     []string
}

var pageTemplateContent = `<!DOCTYPE html>
<html>
	<head>
		<title>{{ .Title }}</title>
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
			<h1>{{ .Title }}</h1>
			{{ .Content }}
			{{ if .Backlinks }}
				<section>
					<h2>Backlinks</h2>
					<ul>
					{{ range .Backlinks }}
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
