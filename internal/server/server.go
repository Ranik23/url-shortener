package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)


type ShortenerServer struct {
	server *http.Server
	log slog.Logger
}

func NewShortenerServer(addr string, handler http.Handler, logger *slog.Logger) *ShortenerServer {
	return &ShortenerServer{
		server: &http.Server{
			Addr: addr,
			Handler: handler,
		},
		log: *logger,
	}
}

func (ss *ShortenerServer) ListenServe() error {

	ctxGroup, cancel := context.WithCancel(context.Background())
	defer cancel()

	errGroup, ctx := errgroup.WithContext(ctxGroup)
	stopCh := make(chan os.Signal, 1)

	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	errGroup.Go(func() error {
		ss.log.Info("Starting server", slog.String("addr", ss.server.Addr))
		if err := ss.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ss.log.Error("ListenAndServe error", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		ss.log.Info("Selector starting...")
		select {
		case <-stopCh:
			ss.log.Info("Received shutdown signal")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err := errGroup.Wait(); err != nil {
		ss.log.Error("Server encountered an error", slog.String("error", err.Error()))
		return err
	}


	if err := ss.ShutDown(context.Background()); err != nil {
		ss.log.Error("Server shutdown failed", slog.String("error", err.Error()))
	}

	ss.log.Info("Server is shutdown")

	return nil
}


func (ss *ShortenerServer) ShutDown(ctx context.Context) error {

	ss.log.Info("Shutting down the server")

	if err := ss.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

