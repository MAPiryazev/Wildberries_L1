package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultAddr         = ":80"
	defaultShutdown     = 3 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdown,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		logrus.Infof("starting http server on %s", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			logrus.Errorf("server stopped with error: %v", err)
		} else {
			logrus.Info("http server stopped gracefully")
		}
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	logrus.Infof("shutting down http server (timeout %v)...", s.shutdownTimeout)
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	logrus.Info("server shutdown complete")
	return nil
}
