package server

import (
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestShortenerServer_StartAndShutdown(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := NewShortenerServer(":2345", handler, logger) 
	done   := make(chan struct{})

	go func() {
		defer close(done)
		if err := server.ListenServe(); err != nil {
			t.Errorf("Server failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost" + server.server.Addr)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: got %v, want %v", resp.StatusCode, http.StatusOK)
	}

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}
	_ = p.Signal(syscall.SIGTERM)


	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Server did not shut down in time")
	}

	// Проверяем, что сервер больше не отвечает
	_, err = http.Get("http://localhost" + server.server.Addr)
	if err == nil {
		t.Fatal("Server is still responding after shutdown")
	}
}
