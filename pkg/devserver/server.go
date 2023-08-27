package devserver

import (
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
	"github.com/fogo-sh/almanac/pkg/static"
	"github.com/fogo-sh/almanac/pkg/templates"
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

type Renderer struct{}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if name != "page" {
		return fmt.Errorf("unknown template: %s", name)
	}

	return templates.PageTemplate.Execute(w, data)
}

func serveNotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "page", templates.PageTemplateData{
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

	return c.Render(http.StatusOK, "page", templates.PageTemplateData{
		Title:   page,
		Content: template.HTML(string(file.ParsedContent)),
	})
}

func NewServer(config Config) *Server {
	echoInst := echo.New()

	var configuredFrontendFS http.FileSystem
	if config.UseBundledAssets {
		configuredFrontendFS = http.FS(static.StaticFS)
	} else {
		configuredFrontendFS = http.Dir("./pkg/static/")
	}

	println(config.UseBundledAssets)

	echoInst.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "./static/",
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
