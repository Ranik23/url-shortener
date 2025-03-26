package grpc

import (
	"context"

	"github.com/Ranik23/url-shortener/api/proto/gen"
)



type ShortenerServer struct {
	gen.UnimplementedURLShortenerServer
}


//go:generate bash -c "impl -dir . 's *ShortenerServer' gen.URLShortenerServer >> grpc.go"
func NewShortenerServer() gen.URLShortenerServer {
	return &ShortenerServer{}
}

func (s *ShortenerServer) ShortenURL(_ context.Context, _ *gen.ShortenRequest) (*gen.ShortenResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ShortenerServer) GetOriginalURL(_ context.Context, _ *gen.GetRequest) (*gen.GetResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ShortenerServer) GetStats(_ context.Context, _ *gen.StatsRequest) (*gen.StatsResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ShortenerServer) DeleteURL(_ context.Context, _ *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	panic("not implemented") // TODO: Implement
}


