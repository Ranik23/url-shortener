package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

//blocking
func (ss *ShortenerServer) ListenServe() {

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := ss.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ss.log.Error("ListenAndServe error", slog.String("error", err.Error()))
		}
	}()

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ss.ShutDown(ctx); err != nil {
		ss.log.Error("Server shutdown failed", slog.String("error", err.Error()))
	}

	ss.log.Info("Server is shutdown")
}


func (ss *ShortenerServer) ShutDown(ctx context.Context) error {

	ss.log.Info("Shutting down the server")

	if err := ss.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

