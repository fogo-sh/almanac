package devserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

type Config struct {
	Addr string
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

func NewServer(config Config) *Server {
	echoInst := echo.New()

	echoInst.Use(slogecho.New(slog.Default()))

	return &Server{
		echoInst: echoInst,
		config:   config,
	}
}
