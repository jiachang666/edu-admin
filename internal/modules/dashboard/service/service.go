package service

import eduservice "edu-admin/internal/modules/edu/service"

type Service struct {
	eduService *eduservice.Service
}

func New(eduService *eduservice.Service) *Service {
	return &Service{eduService: eduService}
}

func (s *Service) Overview() (map[string]any, error) {
	return s.eduService.Overview()
}
