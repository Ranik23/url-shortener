package grpc

import (
	"context"

	"github.com/Ranik23/url-shortener/api/proto/gen"
	"github.com/Ranik23/url-shortener/internal/service"
)



type ShortenerServer struct {
	gen.UnimplementedURLShortenerServer
	service service.Service
}

//go:generate bash -c "impl -dir . 's *ShortenerServer' gen.URLShortenerServer >> grpc.go"
func NewShortenerServer(service service.Service) gen.URLShortenerServer {
	return &ShortenerServer{
		service: service,
	}
}

func (s *ShortenerServer) ShortenURL(ctx context.Context, req *gen.ShortenRequest) (*gen.ShortenResponse, error) {
	shortURL, err := s.service.CreateShortURL(ctx, req.GetOriginalUrl())
	if err != nil {
		return nil, err
	}
	return &gen.ShortenResponse{ShortenedUrl: shortURL}, nil
}

func (s *ShortenerServer) GetOriginalURL(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	originalURL, err := s.service.ResolveShortURL(ctx, req.GetShortenedUrl())
	if err != nil {
		return nil, service.ErrNotFound
	}
	return &gen.GetResponse{OriginalUrl: originalURL}, nil
}

func (s *ShortenerServer) GetStats(ctx context.Context, req *gen.StatsRequest) (*gen.StatsResponse, error) {
	stats, err := s.service.GetStats(ctx, req.GetShortenedUrl())
	if err != nil {
		return nil, err
	}
	return &gen.StatsResponse{OriginalUrl: stats.(string)}, nil // переделать
}

func (s *ShortenerServer) DeleteURL(ctx context.Context, req *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	err := s.service.DeleteShortURL(ctx, req.GetShortenedUrl())
	if err != nil {
		return nil, err
	}
	return &gen.DeleteResponse{Message: "successfull"}, nil
}


