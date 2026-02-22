package service

const (
	validUsername = "student"
	validPassword = "student"
	validToken    = "demo-token"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(username, password string) (string, bool) {
	if username == validUsername && password == validPassword {
		return validToken, true
	}
	return "", false
}

func (s *AuthService) VerifyToken(token string) bool {
	return token == validToken
}
