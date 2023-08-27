package devserver

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/fogo-sh/almanac/pkg/content"
)

type Config struct {
	Addr             string
	ContentDir       string
	UseBundledAssets bool
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

var pageTemplate *template.Template

type PageData struct {
	Title   string
	Content template.HTML
}

func init() {
	t, err := template.New("page").Parse(pageTemplateContent)
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

func serveNotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "page", PageData{
		Title:   "Not found!",
		Content: template.HTML("<p>Looks like this page doesn't exist yet</p>"),
	})
}

func (s *Server) servePage(c echo.Context) error {
	page := c.Param("page")
	contentPath := path.Join(s.config.ContentDir, page+".md")

	_, err := os.Stat(contentPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return serveNotFound(c)
		}
		return fmt.Errorf("error checking file: %w", err)
	}

	file, err := content.ParseFile(contentPath)
	if err != nil {
		return fmt.Errorf("error processing file: %w", err)
	}

	return c.Render(http.StatusOK, "page", PageData{
		Title:   page,
		Content: template.HTML(string(file.ParsedContent)),
	})
}

//go:embed static
var staticFS embed.FS

func NewServer(config Config) *Server {
	echoInst := echo.New()

	var configuredFrontendFS http.FileSystem
	if config.UseBundledAssets {
		configuredFrontendFS = http.FS(staticFS)
	} else {
		configuredFrontendFS = http.Dir("./pkg/devserver/static/")
	}

	echoInst.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: configuredFrontendFS,
	}))

	echoInst.Renderer = &Renderer{}

	echoInst.Use(slogecho.New(slog.Default()))

	server := &Server{
		echoInst: echoInst,
		config:   config,
	}

	echoInst.GET("/:page", server.servePage)

	echoInst.GET("*", serveNotFound)

	return server
}
