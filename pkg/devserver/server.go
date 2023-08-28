package devserver

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"

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
		Content: template.HTML("<p>Looks like this page doesn't exist yet</p>"),
		Page: &content.Page{
			Title: "Not Found",
		},
	})
}

func (s *Server) servePage(c echo.Context) error {
	pageKey := c.Param("page")

	pages, err := content.DiscoverPages(s.config.ContentDir)

	if err != nil {
		return fmt.Errorf("error discovering pages: %w", err)
	}

	var page *content.Page

	if pageKey == "" {
		page, err = content.FindRootPage(pages)

		if err != nil {
			return serveNotFound(c)
		}
	} else {
		var ok bool
		page, ok = pages[pageKey]

		if !ok {
			return serveNotFound(c)
		}
	}

	allPageTitles := content.AllPageTitles(pages)

	return c.Render(http.StatusOK, "page", templates.PageTemplateData{
		AllPageTitles: allPageTitles,
		Content:       template.HTML(string(page.ParsedContent)),
		Page:          page,
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
	echoInst.GET("/", server.servePage)

	echoInst.GET("*", serveNotFound)

	return server
}
