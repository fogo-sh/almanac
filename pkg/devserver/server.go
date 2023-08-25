package devserver

import (
	"errors"
	"fmt"
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

func (s *Server) servePage(c echo.Context) error {
	file, err := content.ParseFile(path.Join(s.config.ContentDir, c.Param("page")+".md"))
	if err != nil {
		return fmt.Errorf("error processing file: %w", err)
	}

	return c.HTMLBlob(http.StatusOK, file.ParsedContent)
}

func NewServer(config Config) *Server {
	echoInst := echo.New()

	echoInst.Use(slogecho.New(slog.Default()))

	server := &Server{
		echoInst: echoInst,
		config:   config,
	}

	echoInst.GET("/:page", server.servePage)

	return server
}
