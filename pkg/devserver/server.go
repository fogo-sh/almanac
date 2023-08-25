package devserver

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"

	"github.com/fogo-sh/almanac/pkg/content"
)

type Config struct {
	Addr       string
	ContentDir string
}

type Server struct {
	echoInst *echo.Echo
	config   Config
}

func (s *Server) Start() error {
	err := s.echoInst.Start(s.config.Addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

type PageTemplate struct {
	Title   string
	Content string
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

var pageTemplate *template.Template = nil

type PageData struct {
	Title   string
	Content template.HTML
}

func init() {
	var t, err = template.New("page").Parse(pageTemplateContent)
	if err != nil {
		panic(err)
	}
	pageTemplate = t
}

type Renderer struct{}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if name != "page" {
		return fmt.Errorf("unknown template: %s", name)
	}

	return pageTemplate.Execute(w, data)
}

func (s *Server) servePage(c echo.Context) error {
	page := c.Param("page")

	file, err := content.ParseFile(path.Join(s.config.ContentDir, page+".md"))
	if err != nil {
		return fmt.Errorf("error processing file: %w", err)
	}

	return c.Render(http.StatusOK, "page", PageData{
		Title:   page,
		Content: template.HTML(string(file.ParsedContent)),
	})
}

func NewServer(config Config) *Server {
	echoInst := echo.New()

	echoInst.Renderer = &Renderer{}

	echoInst.Use(slogecho.New(slog.Default()))

	server := &Server{
		echoInst: echoInst,
		config:   config,
	}

	echoInst.Static("/assets", "assets")

	echoInst.GET("/:page", server.servePage)

	return server
}
