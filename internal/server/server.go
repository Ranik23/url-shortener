package server

import (
	"context"
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
		
		errCh := make(chan error, 1)

		go func() {
			errCh <- ss.server.ListenAndServe()
		}()

		select {
		case <-ctx.Done():
			ss.log.Info("Context canceled, shutting down server...")
		case err := <-errCh:
			if err != nil && err != http.ErrServerClosed {
				ss.log.Error("Server error", slog.String("error", err.Error()))
				return err
			}
		}
		return nil
	})

	errGroup.Go(func() error {
		select {
		case <-stopCh:
			ss.log.Info("Received shutdown signal")
			cancel() 
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Ожидаем завершения всех горутин
	if err := errGroup.Wait(); err != nil {
		ss.log.Error("Server encountered an error", slog.String("error", err.Error()))
		return err
	}

	// Завершаем работу сервера
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

