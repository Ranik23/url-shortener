package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/Ranik23/url-shortener/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandlerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	log.Printf("gRPC method called: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInternal):
			return nil, status.Errorf(codes.Internal, "Internal Server Error")
		case errors.Is(err, service.ErrNotFound):
			return nil, status.Errorf(codes.NotFound, "Not Found")
		case errors.Is(err, service.ErrEmptyURL):
			return nil, status.Errorf(codes.InvalidArgument, "Empty URL")
		default:
			return nil, status.Errorf(codes.Unknown, "Unknown Mistake")
		}
	}
	return resp, nil
}
