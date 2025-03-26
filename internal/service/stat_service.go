package service

import "context"

type StatService interface {
	GetStats(ctx context.Context, shortURL string) (any, error)
}

type statService struct {}


// здесь будет клинет gRPC который будет отправлять запрос на сервис статистики
func (s *statService) GetStats(ctx context.Context, shortURL string) (any, error) {
	panic("unimplemented")
}

func NewStatService() StatService {
	return &statService{}
}
