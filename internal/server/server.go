package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	Host         string
	Port         string
	Handler      http.Handler
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ErrorLog     *log.Logger
	logger       *zerolog.Logger
	WG           *sync.WaitGroup
}

type Options func(*Server)

func WithHost(host string) Options {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port string) Options {
	return func(s *Server) {
		s.Port = port
	}
}

func WithHandler(handler http.Handler) Options {
	return func(s *Server) {
		s.Handler = handler
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Options {
	return func(s *Server) {
		s.IdleTimeout = idleTimeout
	}
}

func WithReadTimeout(readTimeout time.Duration) Options {
	return func(s *Server) {
		s.ReadTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Options {
	return func(s *Server) {
		s.WriteTimeout = writeTimeout
	}
}

func WithErrorLog(errorLog *log.Logger) Options {
	return func(s *Server) {
		s.ErrorLog = errorLog
	}
}

func WithWG(wg *sync.WaitGroup) Options {
	return func(s *Server) {
		s.WG = wg
	}
}

func WithLogger(logger *zerolog.Logger) Options {
	return func(s *Server) {
		s.logger = logger
	}
}

func (s *Server) Connect() error {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      s.Handler,
		IdleTimeout:  s.IdleTimeout,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		ErrorLog:     s.ErrorLog,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		se := <-quit

		s.logger.Info().Str("signal", se.String()).Msg("caught shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		s.logger.Info().Str("addr", server.Addr).Msg("completing background task")
		s.WG.Wait()
		shutdownError <- nil
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownError; err != nil {
		return err
	}

	s.logger.Info().Str("addr", addr).Msg("server stopped")

	return nil
}

func NewServer(opts ...Options) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
