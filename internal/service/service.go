package service

type Service struct {
	AuthService *AuthService
}

func NewService(authService *AuthService) *Service {
	return &Service{
		AuthService: authService,
	}
}
