package service

import "edu-admin/internal/modules/demo"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Overview() map[string]any {
	return demo.Overview()
}
