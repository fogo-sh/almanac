package devserver

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"golang.org/x/oauth2"

	"github.com/fogo-sh/almanac/pkg/content"
	"github.com/fogo-sh/almanac/pkg/static"
)

type Config struct {
	Addr             string
	ContentDir       string
	UseBundledAssets bool

	UseDiscordOAuth     bool
	DiscordClientId     string
	DiscordClientSecret string
	DiscordCallbackUrl  string
	SessionSecret       string
}

type Server struct {
	echoInst *echo.Echo
	config   Config
	oauth    *oauth2.Config
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

	return content.PageTemplate.Execute(w, data)
}

func serveNotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "page", content.PageTemplateData{
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

	return c.Render(http.StatusOK, "page", content.PageTemplateData{
		AllPageTitles: allPageTitles,
		Content:       template.HTML(string(page.ParsedContent)),
		Page:          page,
	})
}

func NewServer(config Config) *Server {
	slog.Debug(
		"Creating server",
		"config", config,
	)

	echoInst := echo.New()

	var oauthConfig *oauth2.Config

	if config.UseDiscordOAuth {
		if config.DiscordClientId == "" {
			slog.Error("Discord OAuth enabled but missing client_id value")
			os.Exit(1)
		}

		if config.DiscordClientSecret == "" {
			slog.Error("Discord OAuth enabled but missing client_secret value")
			os.Exit(1)
		}

		if config.DiscordCallbackUrl == "" {
			slog.Error("Discord OAuth enabled but missing callback_url value")
			os.Exit(1)
		}

		if config.SessionSecret == "" {
			slog.Error("Discord OAuth enabled but missing session_secret value")
			os.Exit(1)
		}

		oauthConfig = &oauth2.Config{
			ClientID:     config.DiscordClientId,
			ClientSecret: config.DiscordClientSecret,
			Scopes:       []string{"identify"},
			RedirectURL:  config.DiscordCallbackUrl,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discordapp.com/api/oauth2/authorize",
				TokenURL: "https://discordapp.com/api/oauth2/token",
			},
		}
	}

	echoInst.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))

	var configuredFrontendFS http.FileSystem
	if config.UseBundledAssets {
		configuredFrontendFS = http.FS(static.StaticFS)
	} else {
		configuredFrontendFS = http.Dir("./pkg/static/")
	}

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
		oauth:    oauthConfig,
	}

	echoInst.GET("/:page", server.servePage)
	echoInst.GET("/", server.servePage)

	echoInst.GET("*", serveNotFound)

	return server
}
