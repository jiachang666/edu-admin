package service

type Service struct {
	devToken string
}

func New(devToken string) *Service {
	return &Service{devToken: devToken}
}

func (s *Service) Login(username string) map[string]any {
	return map[string]any{
		"accessToken": s.devToken,
		"expiresIn":   7200,
		"user": map[string]any{
			"id":          1,
			"username":    username,
			"displayName": "System Admin",
		},
		"permissions": []string{
			"dashboard:view",
			"students:view",
			"classes:view",
		},
	}
}

func (s *Service) Me() map[string]any {
	return map[string]any{
		"id":          1,
		"username":    "admin",
		"displayName": "System Admin",
		"roles":       []string{"super_admin"},
	}
}

func (s *Service) Permissions() []string {
	return []string{
		"dashboard:view",
		"users:view",
		"roles:view",
		"teachers:view",
		"students:view",
		"courses:view",
		"classes:view",
		"schedules:view",
		"attendance:view",
		"notices:view",
	}
}
