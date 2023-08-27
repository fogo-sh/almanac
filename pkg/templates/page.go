package templates

import "html/template"

type PageTemplateData struct {
	Title   string
	Content template.HTML
}

var pageTemplateContent = `<!DOCTYPE html>
<html>
	<head>
		<title>{{ .Title }}</title>
		<link rel="stylesheet" href="/assets/css/main.css">
	</head>
	<body>
		<h1>{{ .Title }}</h1>
		{{ .Content }}
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
