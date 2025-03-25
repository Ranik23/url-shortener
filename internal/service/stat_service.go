package service

import "context"

type statService struct {}

func (s *statService) GetStats(ctx context.Context, shortURL string) (any, error) {
	panic("unimplemented")
}

func NewStatService() StatService {
	return &statService{}
}
