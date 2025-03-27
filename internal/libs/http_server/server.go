package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
	config Config
}

type Config struct {
	Host            string
	Port            int
	StartMsg        string
	ShutdownTimeout time.Duration
}

func New(logger *slog.Logger, config Config, handler http.Handler) *Server {
	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
	}

	s := Server{
		logger: logger,
		server: server,
	}

	return &s
}

func (a *Server) Start(ctx context.Context) error {
	//a.logger.Info(a.config.StartMsg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		err := a.server.Shutdown(ctx) //nolint:contextcheck // sic
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		err := a.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				// ok
			} else {
				return err
			}
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}
